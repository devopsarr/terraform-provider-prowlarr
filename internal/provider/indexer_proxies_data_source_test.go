package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxiesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccIndexerProxyResourceConfig("datasourceTest", 50),
			},
			// Read testing
			{
				Config: testAccIndexerProxiesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.prowlarr_indexer_proxies.test", "indexer_proxies.*", map[string]string{"request_timeout": "50"}),
				),
			},
		},
	})
}

const testAccIndexerProxiesDataSourceConfig = `
data "prowlarr_indexer_proxies" "test" {
}
`
