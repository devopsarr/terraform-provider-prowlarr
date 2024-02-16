package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationSonarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationSonarrResourceConfig("resourceSonarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationSonarrResourceConfig("resourceSonarrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_sonarr.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_sonarr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationSonarrResourceConfig("resourceSonarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationSonarrResourceConfig("resourceSonarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_sonarr.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_application_sonarr.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationSonarrResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_sonarr" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:8989"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [5010, 5020]
		anime_sync_categories = [5070]
	}`, name, prowlarr)
}
