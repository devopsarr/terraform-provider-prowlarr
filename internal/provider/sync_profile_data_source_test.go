package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSyncProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccSyncProfileDataSourceConfig("\"error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccSyncProfileDataSourceConfig("\"error\""),
				ExpectError: regexp.MustCompile("Unable to find sync_profile"),
			},
			// Read testing
			{
				Config: testAccSyncProfileResourceConfig("dataTest", "false") + testAccSyncProfileDataSourceConfig("prowlarr_sync_profile.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_sync_profile.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_sync_profile.test", "minimum_seeders", "1")),
			},
		},
	})
}

func testAccSyncProfileDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "prowlarr_sync_profile" "test" {
		name = %s
	}
	`, name)
}
