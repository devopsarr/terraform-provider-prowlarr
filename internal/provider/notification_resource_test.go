package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationResourceConfig("resourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification.test", "on_health_issue", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_notification.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationResourceConfig("resourceTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification.test", "on_health_issue", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification" "test" {
		on_health_issue                    = %s
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		implementation  = "CustomScript"
		config_contract = "CustomScriptSettings"
	  
		path = "/scripts/test.sh"
	}`, upgrade, name)
}
