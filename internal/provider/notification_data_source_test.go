package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccNotificationDataSourceConfig("\"error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccNotificationDataSourceConfig("\"error\""),
				ExpectError: regexp.MustCompile("Unable to find notification"),
			},
			// Read testing
			{
				Config: testAccNotificationResourceConfig("notificationData", "false") + testAccNotificationDataSourceConfig("prowlarr_notification.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

func testAccNotificationDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "prowlarr_notification" "test" {
		name = %s
	}
	`, name)
}
