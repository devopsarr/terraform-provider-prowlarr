package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccHostResourceConfig("Prowlarr", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccHostResourceConfig("Prowlarr", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_host.test", "port", "9696"),
					resource.TestCheckResourceAttrSet("prowlarr_host.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccHostResourceConfig("Prowlarr", "test") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("ProwlarrTest", "test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_host.test", "port", "9696"),
				),
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("ProwlarrTest", "test1234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_host.test", "port", "9696"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "prowlarr_host.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test1234",
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHostResourceConfig(name, pass string) string {
	return fmt.Sprintf(`
	resource "prowlarr_host" "test" {
		launch_browser = true
		port = 9696
		url_base = ""
		bind_address = "*"
		application_url =  ""
		instance_name = "%s"
		proxy = {
			enabled = false
		}
		ssl = {
			enabled = false
			certificate_validation = "enabled"
		}
		logging = {
			log_level      = "info"
			log_size_limit = 1
		}
		backup = {
			folder = "/backup"
			interval = 5
			retention = 10
		}
		authentication = {
			method = "forms"
			username = "test"
			password = "%s"
			required = "disabledForLocalAddresses"
		}
		update = {
			mechanism = "docker"
			branch = "develop"
		}
	}`, name, pass)
}
