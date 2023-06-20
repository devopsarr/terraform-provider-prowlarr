package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexerDataSourceConfig("error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccIndexerDataSourceConfig("error"),
				ExpectError: regexp.MustCompile("Unable to find indexer"),
			},
			// Create a resource be read
			{
				Config: testAccIndexerResourceConfig("DataSourceTest", "https://0magnet.co/"),
			},
			// Read testing
			{
				Config: testAccIndexerResourceConfig("DataSourceTest", "https://0magnet.co/") + testAccIndexerDataSourceConfig("DataSourceTest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_indexer.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_indexer.test", "name", "DataSourceTest"),
				),
			},
		},
	})
}

func testAccIndexerDataSourceConfig(label string) string {
	return fmt.Sprintf(`
	data "prowlarr_indexer" "test" {
		name = "%s"
	}
	`, label)
}
