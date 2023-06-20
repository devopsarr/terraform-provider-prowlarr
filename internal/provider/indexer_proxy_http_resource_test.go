package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerProxyHTTPResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_http.test", "username", "UserName"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_http.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_http.test", "username", "User"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_indexer_proxy_http.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerProxyHTTPResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer_proxy_http" "test" {
		name = "%s"
		host = "localhost"
		port = 0
		username = "%s"
		password = "Pass"
	}`, name, user)
}
