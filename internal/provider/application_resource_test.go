package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationResourceConfig("error", "http://localhost:9696") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationResourceConfig("resourceTest", "http://localhost:9696"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application.test", "prowlarr_url", "http://localhost:9696"),
					resource.TestCheckResourceAttrSet("prowlarr_application.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationResourceConfig("error", "http://localhost:9696") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationResourceConfig("resourceTest", "https://localhost:6969"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application.test", "prowlarr_url", "https://localhost:6969"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_application.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application" "test" {
		name = "%s"
		sync_level = "disabled"
		implementation  = "Lidarr"
		config_contract = "LidarrSettings"

		base_url = "http://localhost:8686"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [3000, 3010, 3030]
	}`, name, prowlarr)
}
