package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationReadarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationReadarrResourceConfig("resourceReadarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationReadarrResourceConfig("resourceReadarrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_readarr.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_readarr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationReadarrResourceConfig("resourceReadarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationReadarrResourceConfig("resourceReadarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_readarr.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_application_readarr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationReadarrResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_readarr" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:8787"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [2010, 2020]
	}`, name, prowlarr)
}
