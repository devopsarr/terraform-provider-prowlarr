package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerSchemaDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccIndexerSchemaDataSourceConfig("error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccIndexerSchemaDataSourceConfig("error"),
				ExpectError: regexp.MustCompile("Unable to find indexer_schema"),
			},
			// Read testing
			{
				Config: testAccIndexerSchemaDataSourceConfig("AlphaRatio"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_indexer_schema.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_indexer_schema.test", "name", "AlphaRatio"),
				),
			},
		},
	})
}

func testAccIndexerSchemaDataSourceConfig(label string) string {
	return fmt.Sprintf(`
	data "prowlarr_indexer_schema" "test" {
		name = "%s"
	}
	`, label)
}
