package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApplicationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccApplicationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.prowlarr_application.test", "id"),
					resource.TestCheckResourceAttr("data.prowlarr_application.test", "base_url", "http://localhost:8686")),
			},
		},
	})
}

const testAccApplicationDataSourceConfig = `
resource "prowlarr_application" "test" {
	name                    = "applicationData"
	sync_level = "disabled"
	implementation  = "Lidarr"
	config_contract = "LidarrSettings"

	base_url = "http://localhost:8686"
	prowlarr_url = "http://localhost:9696"
	api_key = "APIKey"
	sync_categories = [3000, 3010, 3030]
}

data "prowlarr_application" "test" {
	name = prowlarr_application.test.name
}
`
