package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApplicationWhisparrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApplicationWhisparrResourceConfig("resourceWhisparrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_whisparr.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_whisparr.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccApplicationWhisparrResourceConfig("resourceWhisparrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_whisparr.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_application_whisparr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationWhisparrResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_whisparr" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:6969"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [6010, 6020]
	}`, name, prowlarr)
}
