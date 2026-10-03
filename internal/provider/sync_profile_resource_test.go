package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/devopsarr/prowlarr-go/prowlarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSyncProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccSyncProfileResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccSyncProfileResourceConfig("ResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_sync_profile.test", "enable_rss", "true"),
					resource.TestCheckResourceAttrSet("prowlarr_sync_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccSyncProfileResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
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

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccSyncProfileResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccSyncProfileResourceConfig("ResourceTest", "false"),
				Check: testAccCheckResourceDisappears("prowlarr_sync_profile.test", func(client *prowlarr.APIClient, id int32) (*http.Response, error) {
					return client.AppProfileAPI.DeleteAppProfile(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccSyncProfileResourceConfig("ResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("prowlarr_sync_profile.test", "id"),
				),
			},
		},
	})
}
