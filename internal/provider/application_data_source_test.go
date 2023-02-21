package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApplicationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccApplicationDataSourceConfig("\"error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccApplicationDataSourceConfig("\"error\""),
				ExpectError: regexp.MustCompile("Unable to find application"),
			},
			// Read testing
			{
				Config: testAccApplicationResourceConfig("applicationData", "https://localhost:6969") + testAccApplicationDataSourceConfig("prowlarr_application.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_application.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_application.test", "base_url", "http://localhost:8686")),
			},
		},
	})
}

func testAccApplicationDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "prowlarr_application" "test" {
		name = %s
	}
	`, name)
}
