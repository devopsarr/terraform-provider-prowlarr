package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxyFlaresolverrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerProxyFlaresolverrResourceConfig("resourceFlaresolverrTest", 10) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerProxyFlaresolverrResourceConfig("resourceFlaresolverrTest", 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_flaresolverr.test", "request_timeout", "10"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_flaresolverr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerProxyFlaresolverrResourceConfig("resourceFlaresolverrTest", 10) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerProxyFlaresolverrResourceConfig("resourceFlaresolverrTest", 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_flaresolverr.test", "request_timeout", "20"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_indexer_proxy_flaresolverr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerProxyFlaresolverrResourceConfig(name string, timeout int) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer_proxy_flaresolverr" "test" {
		name = "%s"
		host = "http://localhost:8191/"
		request_timeout = %d
	}`, name, timeout)
}
