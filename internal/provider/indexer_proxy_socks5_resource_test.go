package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxySocks5Resource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerProxySocks5ResourceConfig("resourceSocks5Test", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerProxySocks5ResourceConfig("resourceSocks5Test", "UserName"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_socks5.test", "username", "UserName"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_socks5.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerProxySocks5ResourceConfig("resourceSocks5Test", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerProxySocks5ResourceConfig("resourceSocks5Test", "User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_socks5.test", "username", "User"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_indexer_proxy_socks5.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerProxySocks5ResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer_proxy_socks5" "test" {
		name = "%s"
		host = "localhost"
		port = 0
		username = "%s"
		password = "Pass"
	}`, name, user)
}
