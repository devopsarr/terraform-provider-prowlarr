package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientDelugeResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge-host") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_download_client_deluge.test", "host", "deluge"),
					resource.TestCheckResourceAttr("prowlarr_download_client_deluge.test", "url_base", "/deluge/"),
					resource.TestCheckResourceAttrSet("prowlarr_download_client_deluge.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge-host") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientDelugeResourceConfig("resourceDelugeTest", "deluge-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_download_client_deluge.test", "host", "deluge-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_download_client_deluge.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientDelugeResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "prowlarr_download_client_deluge" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/deluge/"
		port = 9091
		categories = [
			{
				name = "test"
				categories = [1000]
			}
		]
	}`, name, host)
}
