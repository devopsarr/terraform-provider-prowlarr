package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerProxyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerProxyResourceConfig("resourceTest", 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy.test", "request_timeout", "60"),
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy.test", "host", "http://localhost:8191/"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerProxyResourceConfig("resourceTest", 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy.test", "request_timeout", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_indexer_proxy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerProxyResourceConfig(name string, timeout int) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer_proxy" "test" {
		name = "%s"
		implementation = "FlareSolverr"
    	config_contract = "FlareSolverrSettings"
		host = "http://localhost:8191/"
		request_timeout = %d
	}`, name, timeout)
}
