package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSyncProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccSyncProfileResourceConfig("datasourceTest", "false"),
			},
			// Read testing
			{
				Config: testAccSyncProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.prowlarr_sync_profiles.test", "sync_profiles.*", map[string]string{"enable_rss": "true"}),
				),
			},
		},
	})
}

const testAccSyncProfilesDataSourceConfig = `
data "prowlarr_sync_profiles" "test" {
}
`
