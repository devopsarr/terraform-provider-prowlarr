package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSignalResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSignalResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_signal.test", "auth_password", "token123"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_signal.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSignalResourceConfig("error", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSignalResourceConfig("resourceSignalTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_signal.test", "auth_password", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_signal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSignalResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_signal" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"

		auth_username = "User"
		auth_password = "%s"

		host = "localhost"
		port = 8080
		use_ssl = true
		sender_number = "1234"
		receiver_id = "4321"
	}`, name, token)
}
