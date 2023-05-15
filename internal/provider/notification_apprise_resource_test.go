package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationAppriseResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationAppriseResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_apprise.test", "auth_password", "key1"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_apprise.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationAppriseResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_apprise.test", "auth_password", "key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_apprise.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationAppriseResourceConfig(name, key string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_apprise" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"

		server_url = "http://localhost:8000"
		configuration_key = "ConfigKey"
		auth_username = "User"
		auth_password = "%s"
		field_tags = ["warning","skull"]
	}`, name, key)
}
