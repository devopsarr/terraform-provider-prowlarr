package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSyncProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSyncProfileResourceConfig("ResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_sync_profile.test", "enable_rss", "true"),
					resource.TestCheckResourceAttrSet("prowlarr_sync_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccSyncProfileResourceConfig("ResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_sync_profile.test", "enable_rss", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_sync_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSyncProfileResourceConfig(name, rss string) string {
	return fmt.Sprintf(`
		resource "prowlarr_sync_profile" "test" {
  			name = "%s"
			minimum_seeders = 1
			enable_rss = %s
			enable_automatic_search = true
			enable_interactive_search = true
		}
	`, name, rss)
}
