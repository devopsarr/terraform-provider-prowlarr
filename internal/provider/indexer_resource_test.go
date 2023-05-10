package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerResourceConfig("resourceTest", "https://0magnet.co/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "https://0magnet.co/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// resource.TestCheckResourceAttr("prowlarr_indexer.test", "enable_automatic_search", "false"),
					resource.TestCheckTypeSetElemNestedAttrs("prowlarr_indexer.test", "fields.*", map[string]string{"name": "baseSettings.queryLimit"}),
					resource.TestCheckResourceAttrSet("prowlarr_indexer.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerResourceConfig("resourceTest", "https://0magnet.co/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "https://13mag.net/"),
				Check:  resource.ComposeAggregateTestCheckFunc(
				// resource.TestCheckResourceAttr("prowlarr_indexer.test", "enable_automatic_search", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerResourceConfig(name, url string) string {
	return fmt.Sprintf(`
	resource "prowlarr_indexer" "test" {
		enable = false
		name = "%s"
		implementation = "Cardigann"
    	config_contract = "CardigannSettings"
		protocol = "torrent"
		tags = []

		fields = [
			{
				name = "definitionFile"
				text_value = "0magnet"
			},
			{
				name = "baseUrl"
				text_value = "%s"
			},
			{
				name = "baseSettings.queryLimit"
				number_value = 2
			},
			{
				name = "torrentBaseSettings.seedRatio"
				number_value = 0.5
			}
		]
	}

	resource "prowlarr_indexer" "test2" {
		enable = false
		name = "HDits"
		implementation = "HDBits"
    	config_contract = "HDBitsSettings"
		protocol = "torrent"
		tags = []

		fields = [
			{
				name = "username"
				text_value = "test"
			},
			{
				name = "apiKey"
				text_value = "test"
			},
			{
				name = "codecs"
				set_value = [1,5]
			},
			{
				name = "mediums"
				set_value = [1,3]
			},
			{
				name = "torrentBaseSettings.seedRatio"
				number_value = 0.5
			},
			{
				name = "torrentBaseSettings.seedTime"
				number_value = 5
			},
		]
	}
	`, name, url)
}
