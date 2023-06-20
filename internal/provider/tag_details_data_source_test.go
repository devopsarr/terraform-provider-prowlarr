package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTagDetailsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccTagDetailsDataSourceConfig("error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccTagDetailsDataSourceConfig("error"),
				ExpectError: regexp.MustCompile("Unable to find tag"),
			},
			// Create a resource be read
			{
				Config: testAccTagResourceConfig("test", "tag_details_datasource"),
			},
			// Read testing
			{
				Config: testAccTagResourceConfig("test", "tag_details_datasource") + testAccTagDetailsDataSourceConfig("tag_details_datasource"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_tag_details.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_tag_details.test", "label", "tag_details_datasource"),
				),
			},
		},
	})
}

func testAccTagDetailsDataSourceConfig(label string) string {
	return fmt.Sprintf(`
	data "prowlarr_tag_details" "test" {
		label = "%s"
	}
	`, label)
}
