package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationDiscordResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationDiscordResourceConfig("error", "dog-picture") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationDiscordResourceConfig("resourceDiscordTest", "dog-picture"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_discord.test", "avatar", "dog-picture"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_discord.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationDiscordResourceConfig("error", "dog-picture") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationDiscordResourceConfig("resourceDiscordTest", "cat-picture"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_discord.test", "avatar", "cat-picture"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_discord.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationDiscordResourceConfig(name, avatar string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_discord" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		web_hook_url = "http://discord-web-hook.com"
		username = "User"
		avatar = "%s"
		grab_fields = [0,1,2,3,4,5,6,7,8,9]
	}`, name, avatar)
}
