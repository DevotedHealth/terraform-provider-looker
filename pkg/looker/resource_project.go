package looker

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Exists: resourceProjectExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	body := apiclient.WriteProject{}

	projectName := d.Get("name").(string)
	body.Name = &projectName

	err := selectAPISession(client, DEV_WORKSPACE)
	if err != nil {
		return err
	}

	project, err := client.CreateProject(body, nil)
	if err != nil {
		return err
	}

	d.SetId(*project.Id)

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	project, err := client.Project(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	projectName := project.Name
	d.Set("name", &projectName)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	body := apiclient.WriteProject{}
	projectName := d.Get("name").(string)
	body.Name = &projectName

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	_, err := client.UpdateProject(d.Id(), body, "", nil)
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: Looker doesn't appear to support deleting projects from the API
	return nil
}

func resourceProjectExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return false, err
	}

	_, err := client.Project(d.Id(), "", nil)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func resourceProjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
