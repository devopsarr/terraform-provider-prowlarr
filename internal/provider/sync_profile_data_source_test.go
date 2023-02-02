package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSyncProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSyncProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_sync_profile.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_sync_profile.test", "minimum_seeders", "10")),
			},
		},
	})
}

const testAccSyncProfileDataSourceConfig = `
resource "prowlarr_sync_profile" "test" {
	name = "dataTest"
  	minimum_seeders = 10
  	enable_rss = true
  	enable_automatic_search = true
  	enable_interactive_search = true
}

data "prowlarr_sync_profile" "test" {
	name = prowlarr_sync_profile.test.name
}
`
