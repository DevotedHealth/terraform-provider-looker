package looker

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceLookMLModel() *schema.Resource {
	return &schema.Resource{
		Create: resourceLookMLModelCreate,
		Read:   resourceLookMLModelRead,
		Update: resourceLookMLModelUpdate,
		Delete: resourceLookMLModelDelete,
		Exists: resourceLookMLModelExists,
		Importer: &schema.ResourceImporter{
			State: resourceLookMLModelImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_db_connection_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceLookMLModelCreate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	params := apiclient.WriteLookmlModel{}

	modelName := d.Get("name").(string)
	params.Name = &modelName

	projectName := d.Get("project_name").(string)
	params.ProjectName = &projectName

	var connectionNames []string
	for _, modelName := range d.Get("allowed_db_connection_names").(*schema.Set).List() {
		connectionNames = append(connectionNames, modelName.(string))
	}
	params.AllowedDbConnectionNames = &connectionNames

	model, err := client.CreateLookmlModel(params, nil)
	if err != nil {
		log.Printf("[WARN] Error creating a model., %s", err.Error())
		return err
	}

	d.SetId(*model.Name)

	return resourceLookMLModelRead(d, m)
}

func resourceLookMLModelRead(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	model, err := client.LookmlModel(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(*model.Name)
	d.Set("name", *model.Name)
	d.Set("project_name", *model.ProjectName)
	d.Set("allowed_db_connection_names", *model.AllowedDbConnectionNames)

	return nil
}

func resourceLookMLModelUpdate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	params := apiclient.WriteLookmlModel{}

	modelName := d.Get("name").(string)
	params.Name = &modelName

	projectName := d.Get("project_name").(string)
	params.ProjectName = &projectName

	var connectionNames []string
	for _, modelName := range d.Get("allowed_db_connection_names").(*schema.Set).List() {
		connectionNames = append(connectionNames, modelName.(string))
	}
	params.AllowedDbConnectionNames = &connectionNames

	_, err := client.UpdateLookmlModel(d.Id(), params, nil)
	if err != nil {
		return err
	}

	return resourceLookMLModelRead(d, m)
}

func resourceLookMLModelDelete(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	_, err := client.DeleteLookmlModel(d.Id(), nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceLookMLModelExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	_, err := client.LookmlModel(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceLookMLModelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceLookMLModelRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
