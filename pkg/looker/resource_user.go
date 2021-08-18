package looker

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserImport,
		},
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)

	writeUser := apiclient.WriteUser{
		FirstName: &firstName,
		LastName:  &lastName,
	}

	user, err := client.CreateUser(writeUser, "", nil)
	if err != nil {
		return err
	}

	userID := *user.Id
	d.SetId(strconv.Itoa(int(userID)))

	writeCredentialsEmail := apiclient.WriteCredentialsEmail{
		Email: &email,
	}
	_, err = client.CreateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
	if err != nil {
		if _, err = client.DeleteUser(userID, nil); err != nil {
			return err
		}
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	user, err := client.User(userID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("email", user.Email); err != nil {
		return err
	}
	if err = d.Set("first_name", user.FirstName); err != nil {
		return err
	}
	if err = d.Set("last_name", user.LastName); err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	if d.HasChanges("first_name", "last_name") {
		firstName := d.Get("first_name").(string)
		lastName := d.Get("last_name").(string)
		writeUser := apiclient.WriteUser{
			FirstName: &firstName,
			LastName:  &lastName,
		}
		_, err = client.UpdateUser(userID, writeUser, "", nil)
		if err != nil {
			return err
		}
	}

	if d.HasChange("email") {
		email := d.Get("email").(string)
		writeCredentialsEmail := apiclient.WriteCredentialsEmail{
			Email: &email,
		}
		_, err = client.UpdateUserCredentialsEmail(userID, writeCredentialsEmail, "", nil)
		if err != nil {
			return err
		}
	}

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	client := apiclient.NewLookerSDK(m.(*rtl.AuthSession))

	userID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	_, err = client.DeleteUser(userID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
