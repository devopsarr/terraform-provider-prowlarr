package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationMylarResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationMylarResourceConfig("resourceMylarTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationMylarResourceConfig("resourceMylarTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_mylar.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_mylar.test", "id"),
				),
			},

			// Unauthorized Read
			{
				Config:      testAccApplicationMylarResourceConfig("resourceMylarTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationMylarResourceConfig("resourceMylarTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_mylar.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_application_mylar.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationMylarResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_mylar" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:8090"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [7030]
	}`, name, prowlarr)
}
