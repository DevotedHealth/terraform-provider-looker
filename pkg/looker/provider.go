package looker

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
)

const (
	defaultAPIVersion = "4.0"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"ini_file_path": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LOOKER_INI_FILE_PATH", nil),
				Description:  "Path to the looker.ini file with connection information.",
				AtLeastOneOf: []string{"ini_file_path", "client_id", "client_secret"},
			},
			"ini_section": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_INI_SECTION", nil),
				Description: "Section of the ini to use with this connection. Default: Looker",
			},
			"client_id": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LOOKER_API_CLIENT_ID", nil),
				Description:  "Client ID to authenticate with Looker",
				AtLeastOneOf: []string{"ini_file_path", "client_id", "client_secret"},
			},
			"client_secret": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				DefaultFunc:  schema.EnvDefaultFunc("LOOKER_API_CLIENT_SECRET", nil),
				Description:  "Client Secret to authenticate with Looker",
				AtLeastOneOf: []string{"ini_file_path", "client_id", "client_secret"},
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_BASE_URL", nil),
				Description: "Looker API Base URL",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_API_VERSION", defaultAPIVersion),
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_VERIFY_SSL", true),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LOOKER_TIMEOUT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"looker_user":                   resourceUser(),
			"looker_user_roles":             resourceUserRoles(),
			"looker_permission_set":         resourcePermissionSet(),
			"looker_model_set":              resourceModelSet(),
			"looker_group":                  resourceGroup(),
			"looker_role":                   resourceRole(),
			"looker_role_groups":            resourceRoleGroups(),
			"looker_user_attribute":         resourceUserAttribute(),
			"looker_project":                resourceProject(),
			"looker_project_git_deploy_key": resourceProjectGitDeployKey(),
			"looker_project_git_repo":       resourceProjectGitRepo(),
			"looker_connection":             resourceConnection(),
			"looker_lookml_model":           resourceLookMLModel(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	path := d.Get("ini_file_path").(string)
	section := d.Get("ini_section").(string)

	apiSettings, err := rtl.NewSettingsFromFile(path, &section)
	if err != nil {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "No Looker ini file found",
			Detail:   fmt.Sprintf("%s not found", path),
		})
	}

	baseUrl, baseUrlOk := d.GetOk("base_url")
	clientID, clientIDOk := d.GetOk("client_id")
	clientSecret, clientSecretOk := d.GetOk("client_secret")
	apiVersion, apiVersionOk := d.GetOk("api_version")
	timeout, timeoutOk := d.GetOk("timeout")

	if baseUrlOk {
		apiSettings.BaseUrl = baseUrl.(string)
	}

	if clientIDOk {
		apiSettings.ClientId = clientID.(string)
	}

	if clientSecretOk {
		apiSettings.ClientSecret = clientSecret.(string)
	}

	if apiVersionOk {
		apiSettings.ApiVersion = apiVersion.(string)
	}

	if timeoutOk {
		apiSettings.Timeout = timeout.(int32)
	}

	if apiSettings.ClientId == "" || apiSettings.ClientSecret == "" || apiSettings.BaseUrl == "" {
		diagnostics = append(diagnostics, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No credentials found",
			Detail:   fmt.Sprintf("ClientID/ClientSecret/BaseURL were not found after parsing provider and %s", path),
		})

		return nil, diagnostics
	}

	authSession := rtl.NewAuthSession(apiSettings)
	// client := apiclient.NewLookerSDK(authSession)

	return authSession, diagnostics
}
