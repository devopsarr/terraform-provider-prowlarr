package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerProxySocks4Resource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerProxySocks4ResourceConfig("resourceSocks4Test", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerProxySocks4ResourceConfig("resourceSocks4Test", "UserName"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_socks4.test", "username", "UserName"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_socks4.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerProxySocks4ResourceConfig("resourceSocks4Test", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerProxySocks4ResourceConfig("resourceSocks4Test", "User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_socks4.test", "username", "User"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_indexer_proxy_socks4.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerProxySocks4ResourceConfig(name, user string) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer_proxy_socks4" "test" {
		name = "%s"
		host = "localhost"
		port = 0
		username = "%s"
		password = "Pass"
	}`, name, user)
}
