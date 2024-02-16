package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationLidarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationLidarrResourceConfig("resourceLidarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationLidarrResourceConfig("resourceLidarrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_lidarr.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_lidarr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationLidarrResourceConfig("resourceLidarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationLidarrResourceConfig("resourceLidarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_lidarr.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_application_lidarr.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationLidarrResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_lidarr" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:9696"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [3010, 3020]
	}`, name, prowlarr)
}
