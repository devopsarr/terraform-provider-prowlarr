package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApplicationsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccApplicationsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to have a value to check
			{
				Config: testAccApplicationResourceConfig("datasourceTest", "http://localhost:9696"),
			},
			// Read testing
			{
				Config: testAccApplicationsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.prowlarr_applications.test", "applications.*", map[string]string{"prowlarr_url": "http://localhost:9696"}),
				),
			},
		},
	})
}

const testAccApplicationsDataSourceConfig = `
data "prowlarr_applications" "test" {
}
`
