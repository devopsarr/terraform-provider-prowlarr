package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationWebhookResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_webhook.test", "on_health_issue", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_notification_webhook.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccNotificationWebhookResourceConfig("resourceWebhookTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_notification_webhook.test", "on_health_issue", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_notification_webhook.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationWebhookResourceConfig(name, upgrade string) string {
	return fmt.Sprintf(`
	resource "prowlarr_notification_webhook" "test" {
		on_health_issue                    = %s
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		url = "http://transmission:9091"
		method = 1
	}`, upgrade, name)
}
