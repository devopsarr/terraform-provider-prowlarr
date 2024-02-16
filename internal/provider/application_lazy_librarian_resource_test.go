package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationLazyLibrarianResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_lazy_librarian.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_lazy_librarian.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationLazyLibrarianResourceConfig("resourceLazyLibrarianTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_lazy_librarian.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_application_lazy_librarian.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationLazyLibrarianResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_lazy_librarian" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:5299"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [7010, 7020]
	}`, name, prowlarr)
}
