package looker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceProjectGitRepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGitRepoCreate,
		Read:   resourceProjectGitRepoRead,
		Delete: resourceProjectGitRepoDelete,
		Update: resourceProjectGitRepoUpdate,
		Exists: resourceProjectGitRepoExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectGitRepoImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"github"}, false),
			},
			"git_remote_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_production_branch_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "main",
			},
			"deploy_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"git_release_mgmt_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pull_request_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: validation.StringInSlice([]string{"off", "links", "recommended", "required"}, true),
			},
		},
	}
}

func setProjectGitDetails(d *schema.ResourceData, m interface{}, create bool) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	stringRemoteURL := d.Get("git_remote_url").(string)
	branchName := d.Get("git_production_branch_name").(string)
	releaseManagement := d.Get("git_release_mgmt_enabled").(bool)
	pullRequestMode := apiclient.PullRequestMode(d.Get("pull_request_mode").(string))

	createPayload := apiclient.WriteProject{
		GitRemoteUrl: &stringRemoteURL,
	}

	github := "github"
	params := apiclient.WriteProject{
		GitServiceName:          &github,
		GitRemoteUrl:            &stringRemoteURL,
		GitProductionBranchName: &branchName,
		GitReleaseMgmtEnabled:   &releaseManagement,
		PullRequestMode:         &pullRequestMode,
	}

	if deploySecret, deploySecretOk := d.GetOk("deploy_secret"); deploySecretOk && deploySecret.(string) != "" {
		deploySecretString := deploySecret.(string)
		params.DeploySecret = &deploySecretString
	} else if deploySecretOk && deploySecret.(string) == "" {
		unset := true
		params.UnsetDeploySecret = &unset
	}

	if create {
		if _, err := client.UpdateProject(d.Id(), createPayload, "", nil); err != nil {
			return err
		}
	}

	if _, err := client.UpdateProject(d.Id(), params, "", nil); err != nil {
		return err
	}

	// Deploy won't work if we are requiring PRs in GitHub
	if pullRequestMode != "required" {
		_, err := client.DeployToProduction(d.Id(), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceProjectGitRepoCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId(d.Get("project_id").(string))

	err := setProjectGitDetails(d, m, true)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	return resourceProjectGitRepoRead(d, m)
}

func resourceProjectGitRepoRead(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	result, err := client.Project(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("project_id", *result.Id)
	d.Set("git_remote_url", *result.GitRemoteUrl)
	d.Set("git_production_branch_name", *result.GitProductionBranchName)
	d.Set("git_release_mgmt_enabled", *result.GitReleaseMgmtEnabled)
	d.Set("pull_request_mode", string(*result.PullRequestMode))
	d.Set("deploy_secret", d.Get("deploy_secret").(string))

	return nil
}

func resourceProjectGitRepoUpdate(d *schema.ResourceData, m interface{}) error {
	err := setProjectGitDetails(d, m, false)
	if err != nil {
		return err
	}

	return resourceProjectGitRepoRead(d, m)
}

func resourceProjectGitRepoDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Deleting this resource should set the git fields back to blank values. not implementing this yet since leaving the values does not have any negative effect
	return nil
}

func resourceProjectGitRepoExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return false, err
	}

	_, err := client.Project(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceProjectGitRepoImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectGitRepoRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
