package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagsDetailsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccTagsDetailsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccTagResourceConfig("test-1", "books") + testAccTagResourceConfig("test-2", "comics"),
			},
			// Read testing
			{
				Config: testAccTagsDetailsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.prowlarr_tags_details.test", "tags.*", map[string]string{"label": "books"}),
				),
			},
		},
	})
}

const testAccTagsDetailsDataSourceConfig = `
data "prowlarr_tags_details" "test" {
}
`
