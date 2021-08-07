package looker

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMembershipCreate,
		Read:   resourceGroupMembershipRead,
		Update: resourceGroupMembershipUpdate,
		Delete: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGroupMembershipImport,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGroupMembershipCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	groupIDString := d.Get("group_id").(string)

	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return err
	}

	userIDString := d.Get("user_id").(string)

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return err
	}

	body := apiclient.GroupIdForGroupUserInclusion{
		UserId: &userID,
	}

	_, err = client.AddGroupUser(groupID, body, nil)
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(&groupIDString, &userIDString))

	return resourceGroupMembershipRead(d, m)
}

func resourceGroupMembershipRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	id := d.Id()
	groupID, userID, err := groupIDAndUserIDFromID(id)
	if err != nil {
		return err
	}

	req := apiclient.RequestAllGroupUsers{
		GroupId: groupID,
	}

	users, err := client.AllGroupUsers(req, nil) // todo: imeplement paging
	if err != nil {
		return err
	}

	if !isContains(users, userID) {
		return fmt.Errorf("failed to find target userID: %d", userID)
	}

	if err = d.Set("group_id", strconv.Itoa(int(groupID))); err != nil {
		return nil
	}

	if err = d.Set("user_id", strconv.Itoa(int(userID))); err != nil {
		return nil
	}

	return nil
}

func resourceGroupMembershipUpdate(d *schema.ResourceData, m interface{}) error {
	// there's no update in this resource
	return resourceGroupMembershipCreate(d, m)
}

func resourceGroupMembershipDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerSDK)

	id := d.Id()
	groupID, userID, err := groupIDAndUserIDFromID(id)
	if err != nil {
		return err
	}

	err = client.DeleteGroupUser(groupID, userID, nil)
	if err != nil {
		return err
	}

	return resourceGroupMembershipRead(d, m)
}

func resourceGroupMembershipImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceGroupMembershipRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func isContains(users []apiclient.User, userID int64) bool {
	for _, user := range users {
		if user.Id == &userID {
			return true
		}
	}
	return false
}

func groupIDAndUserIDFromID(id string) (int64, int64, error) {
	groupIDString, userIDString, err := parseTwoPartID(id)
	if err != nil {
		return 0, 0, err
	}

	groupID, err := strconv.ParseInt(groupIDString, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return groupID, userID, err
}
