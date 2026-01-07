package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientNzbvortexResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_download_client_nzbvortex.test", "host", "nzbvortex"),
					resource.TestCheckResourceAttr("prowlarr_download_client_nzbvortex.test", "url_base", "/nzbvortex/"),
					resource.TestCheckResourceAttrSet("prowlarr_download_client_nzbvortex.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_download_client_nzbvortex.test", "host", "nzbvortex-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_download_client_nzbvortex.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientNzbvortexResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "prowlarr_download_client_nzbvortex" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/nzbvortex/"
		port = 4321
		api_key = "testAPIkey"
		categories = [
			{
				name = "test"
				categories = [1000]
			}
		]
	}`, name, host)
}
