package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxyDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIndexerProxyDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_indexer_proxy.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_indexer_proxy.test", "request_timeout", "50")),
			},
		},
	})
}

const testAccIndexerProxyDataSourceConfig = `
resource "prowlarr_indexer_proxy" "test" {
	name = "dataTest"
	implementation = "FlareSolverr"
	config_contract = "FlareSolverrSettings"
	host = "http://localhost:8191/"
	request_timeout = 50
}

data "prowlarr_indexer_proxy" "test" {
	name = prowlarr_indexer_proxy.test.name
}
`
