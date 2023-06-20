package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexersDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccIndexerResourceConfig("DataTest", "https://0magnet.co/"),
			},
			// Read testing
			{
				Config: testAccIndexersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.prowlarr_indexers.test", "indexers.*", map[string]string{"name": "DataTest"}),
				),
			},
		},
	})
}

const testAccIndexersDataSourceConfig = `
data "prowlarr_indexers" "test" {
}
`
