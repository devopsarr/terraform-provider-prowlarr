package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationCustomScriptResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_custom_script.test", "on_health_issue", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_custom_script.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationCustomScriptResourceConfig("resourceScriptTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_custom_script.test", "on_health_issue", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_custom_script.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationCustomScriptResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_custom_script" "test" {
		on_health_issue                    = %s
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		path = "/scripts/test.sh"
	}`, upgrade, name)
}
