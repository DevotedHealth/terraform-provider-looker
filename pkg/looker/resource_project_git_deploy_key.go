package looker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

const (
	gitDeployKeyURL = "%s/api/3.1/projects/%s/git/deploy_key"
)

func resourceProjectGitDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGitDeployKeyCreate,
		Read:   resourceProjectGitDeployKeyRead,
		Delete: resourceProjectGitDeployKeyDelete,
		Exists: resourceProjectGitDeployKeyExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectGitDeployKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProjectGitDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)

	req, _ := http.NewRequest("POST", fmt.Sprintf(gitDeployKeyURL, session.Config.BaseUrl, projectID), nil)
	session.Authenticate(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	publicKey, _ := ioutil.ReadAll(res.Body)
	key := strings.Fields(string(publicKey))
	d.SetId(projectID)
	d.Set("project_id", projectID)
	d.Set("public_key", fmt.Sprintf("%s %s", key[0], key[1]))

	return resourceProjectGitDeployKeyRead(d, m)
}

func resourceProjectGitDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(gitDeployKeyURL, session.Config.BaseUrl, d.Id()), nil)
	session.Authenticate(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	publicKey, _ := ioutil.ReadAll(res.Body)
	key := strings.Fields(string(publicKey))

	d.Set("project_id", d.Id())
	d.Set("public_key", fmt.Sprintf("%s %s", key[0], key[1]))

	return nil
}

func resourceProjectGitDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	// TODO There is no way to delete a git deploy key, possibly put this into the project resource (but there is no way to delete project either)
	return nil
}

func resourceProjectGitDeployKeyExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	if err := selectAPISession(client, DEV_WORKSPACE); err != nil {
		return false, err
	}

	_, err := client.GitDeployKey(d.Id(), nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceProjectGitDeployKeyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectGitDeployKeyRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
