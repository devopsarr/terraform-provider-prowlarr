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

func TestAccApplicationRadarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccApplicationRadarrResourceConfig("resourceRadarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccApplicationRadarrResourceConfig("resourceRadarrTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_radarr.test", "prowlarr_url", "false"),
					resource.TestCheckResourceAttrSet("prowlarr_application_radarr.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccApplicationRadarrResourceConfig("resourceRadarrTest", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccApplicationRadarrResourceConfig("resourceRadarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_application_radarr.test", "prowlarr_url", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "prowlarr_application_radarr.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationRadarrResourceConfig(name, prowlarr string) string {
	return fmt.Sprintf(`
	resource "prowlarr_application_radarr" "test" {
		name = "%s"
		sync_level = "disabled"

		base_url = "http://localhost:7878"
		prowlarr_url = "%s"
		api_key = "APIKey"
		sync_categories = [2010, 2020]
	}`, name, prowlarr)
}

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccApplicationRadarrResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccApplicationRadarrResourceConfig("resourceRadarrTest", "true"),
				Check: testAccCheckResourceDisappears("prowlarr_application_radarr.test", func(client *prowlarr.APIClient, id int32) (*http.Response, error) {
					return client.ApplicationAPI.DeleteApplications(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccApplicationRadarrResourceConfig("resourceRadarrTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("prowlarr_application_radarr.test", "id"),
				),
			},
		},
	})
}
