package looker

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

const (
	USER_ACCESS_EDIT = "edit"
	USER_ACCESS_VIEW = "view"
	USER_ACCESS_NONE = "none"
)

func resourceUserAttribute() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserAttributeCreate,
		Read:   resourceUserAttributeRead,
		Update: resourceUserAttributeUpdate,
		Delete: resourceUserAttributeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceUserAttributeImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"string", "number", "datetime", "relative_url", "advanced_filter_datetime", "advanced_filter_number", "advanced_filter_string"}, true),
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_access": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{USER_ACCESS_EDIT, USER_ACCESS_VIEW, USER_ACCESS_NONE}, true),
			},
		},
	}
}

func resourceUserAttributeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)
	userAttributeName := d.Get("name").(string)
	userAttributeLabel := d.Get("label").(string)
	userAttributeType := d.Get("type").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	if userAttributeDefault, defaultSet := d.GetOk("default"); defaultSet {
		stringDefaultValue := userAttributeDefault.(string)
		writeUserAttribute.DefaultValue = &stringDefaultValue
	}

	if userAccessAttribute, userAccessOk := d.GetOk("user_access"); userAccessOk {
		userCan := true
		userCant := false
		if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_VIEW) {
			writeUserAttribute.UserCanEdit = nil
			writeUserAttribute.UserCanView = &userCan
		} else if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_EDIT) {
			writeUserAttribute.UserCanEdit = &userCan
			writeUserAttribute.UserCanView = nil
		} else if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_NONE) {
			writeUserAttribute.UserCanEdit = &userCant
			writeUserAttribute.UserCanView = &userCant
		}
	}

	userAttribute, err := client.CreateUserAttribute(writeUserAttribute, "", nil)
	if err != nil {
		return err
	}

	userAttributeID := *userAttribute.Id
	d.SetId(strconv.Itoa(int(userAttributeID)))

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	userAttribute, err := client.UserAttribute(userAttributeID, "", nil)
	if err != nil {
		return err
	}

	if err = d.Set("name", userAttribute.Name); err != nil {
		return err
	}
	if err = d.Set("type", userAttribute.Type); err != nil {
		return err
	}
	if err = d.Set("label", userAttribute.Label); err != nil {
		return err
	}

	if _, ok := d.GetOk("default"); ok {
		if err = d.Set("default", userAttribute.DefaultValue); err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("user_access"); ok {
		if *userAttribute.UserCanEdit {
			if err = d.Set("user_access", USER_ACCESS_EDIT); err != nil {
				return err
			}
		} else if *userAttribute.UserCanView {
			if err = d.Set("user_access", USER_ACCESS_VIEW); err != nil {
				return err
			}
		} else {
			if err = d.Set("user_access", USER_ACCESS_NONE); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceUserAttributeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	userAttributeName := d.Get("name").(string)
	userAttributeType := d.Get("type").(string)
	userAttributeLabel := d.Get("label").(string)

	writeUserAttribute := apiclient.WriteUserAttribute{
		Name:  &userAttributeName,
		Label: &userAttributeLabel,
		Type:  &userAttributeType,
	}

	if userAttributeDefault, defaultSet := d.GetOk("default"); defaultSet {
		stringDefaultAttribute := userAttributeDefault.(string)
		writeUserAttribute.DefaultValue = &stringDefaultAttribute
	}

	if userAccessAttribute, userAccessOk := d.GetOk("user_access"); userAccessOk {
		userCan := true
		userCant := false
		if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_VIEW) {
			writeUserAttribute.UserCanEdit = &userCant
			writeUserAttribute.UserCanView = &userCan
		} else if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_EDIT) {
			writeUserAttribute.UserCanEdit = &userCan
			writeUserAttribute.UserCanView = &userCan
		} else if strings.EqualFold(userAccessAttribute.(string), USER_ACCESS_NONE) {
			writeUserAttribute.UserCanEdit = &userCant
			writeUserAttribute.UserCanView = &userCant
		}
	}

	_, err = client.UpdateUserAttribute(userAttributeID, writeUserAttribute, "", nil)
	if err != nil {
		return err
	}

	return resourceUserAttributeRead(d, m)
}

func resourceUserAttributeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	userAttributeID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	_, err = client.DeleteUserAttribute(userAttributeID, nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceUserAttributeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceUserAttributeRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
