package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func TestAcc_UserAttributeUserValue(t *testing.T) {
	user := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	attributeValue1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	attributeValue2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: userAttributeUserValueConfig(user, attributeValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeUserValueExists("looker_user_attribute_user_value.test_user_attr"),
					resource.TestCheckResourceAttr("looker_user_attribute_user_value.test_user_attr", "value", attributeValue1),
				),
			},
			// Test: Update
			{
				Config: userAttributeUserValueConfig(user, attributeValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributeUserValueExists("looker_user_attribute_user_value.test_user_attr"),
					resource.TestCheckResourceAttr("looker_user_attribute_user_value.test_user_attr", "value", attributeValue2),
				),
			},
			// Test: Import
			{
				ResourceName:            "looker_user_attribute_user_value.test_user_attr",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_attribute_id", "user_id", "value"},
			},
		},
		CheckDestroy: testAccCheckUserAttributeUserValueDestroy,
	})
}

func testAccCheckUserAttributeUserValueExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("user attribute user value setting not found from resources: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no user attribute user value setting ID is set")
		}

		client := testAccProvider.Meta().(*apiclient.LookerSDK)
		userID, userAttributeID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		userAttributeIDs := rtl.DelimString{userAttributeID}
		request := apiclient.RequestUserAttributeUserValues{
			UserId:           userID,
			UserAttributeIds: &userAttributeIDs,
		}

		_, err = client.UserAttributeUserValues(request, nil)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckUserAttributeUserValueDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_user_attribute_user_value" {
			continue
		}

		userAttributeID, userID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		userAttributeIDs := rtl.DelimString{userAttributeID}
		request := apiclient.RequestUserAttributeUserValues{
			UserId:           userID,
			UserAttributeIds: &userAttributeIDs,
		}

		userAttributeUserValues, err := client.UserAttributeUserValues(request, nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil // successfully destroyed
			}
		}

		if len(userAttributeUserValues) != 1 {
			return err
		}

		if userAttributeUserValues[0].Value != nil {
			return err
		}
	}

	return nil
}

func userAttributeUserValueConfig(user, attributeValue string) string {
	return fmt.Sprintf(`
	resource "looker_user" "test_user_attr" {
        first_name = "%s"
        last_name  = "%s"
        email      = "%s@jason.com"
	}
	resource "looker_user_attribute" "test_user_attr" {
        name  = "test_x"
        type  = "string"
        label = "test label x"
	}
	resource "looker_user_attribute_user_value" "test_user_attr" {
		user_id           = looker_user.test_user_attr.id
		user_attribute_id = looker_user_attribute.test_user_attr.id
        value             = "%s"
	}
	`, user, user, user, attributeValue)
}
