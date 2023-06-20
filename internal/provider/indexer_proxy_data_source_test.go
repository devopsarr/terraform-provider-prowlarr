package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerProxyDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexerProxyDataSourceConfig("\"error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccIndexerProxyDataSourceConfig("\"error\""),
				ExpectError: regexp.MustCompile("Unable to find indexer_proxy"),
			},
			// Read testing
			{
				Config: testAccIndexerProxyResourceConfig("dataTest", 50) + testAccIndexerProxyDataSourceConfig("prowlarr_indexer_proxy.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_indexer_proxy.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_indexer_proxy.test", "request_timeout", "50")),
			},
		},
	})
}

func testAccIndexerProxyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "prowlarr_indexer_proxy" "test" {
		name = %s
	}
	`, name)
}
