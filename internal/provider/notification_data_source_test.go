package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotificationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

const testAccNotificationDataSourceConfig = `
resource "prowlarr_notification" "test" {
	on_health_issue                    = false
	on_application_update              = false
  
	include_health_warnings = false
	name                    = "notificationData"
  
	implementation  = "CustomScript"
	config_contract = "CustomScriptSettings"
  
	path = "/scripts/test.sh"
}

data "prowlarr_notification" "test" {
	name = prowlarr_notification.test.name
}
`
