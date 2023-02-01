package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxyHTTPResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_http.test", "username", "UserName"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_http.test", "id"),
				),
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
