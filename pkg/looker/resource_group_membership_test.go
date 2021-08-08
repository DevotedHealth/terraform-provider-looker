package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_GroupMembership(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: groupMembershipConfig(name1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_group_membership.test", "group_id", "1"),
					resource.TestCheckResourceAttr("looker_group_membership.test", "user_id", "1"),
				),
			},
		},
	})
}

func groupMembershipConfig(name string) string {
	return fmt.Sprintf(`
	resource "looker_group" "test" {
		name = "%s"
	}
	resource "looker_user" "test" {
		name = "%s"
	}
	resource "looker_group_membership" "test" {
		group_id = looker_group.test.id
		user_id = looker_user.test.id
	}
	`, name, name)
}
