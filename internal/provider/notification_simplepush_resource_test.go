package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationSimplepushResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationSimplepushResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "key1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_simplepush.test", "key", "key1"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_simplepush.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationSimplepushResourceConfig("error", "key1") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationSimplepushResourceConfig("resourceSimplepushTest", "key2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_simplepush.test", "key", "key2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_simplepush.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationSimplepushResourceConfig(name, key string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_simplepush" "test" {
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		key = "%s"
		event = "ringtone:default"
	}`, name, key)
}
