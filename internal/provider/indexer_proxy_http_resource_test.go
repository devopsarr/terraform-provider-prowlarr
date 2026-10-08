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

func TestAccIndexerProxyHTTPResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("prowlarr_indexer_proxy_http.test", "username", "UserName"),
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_http.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "UserName") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
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
				ResourceName:            "prowlarr_indexer_proxy_http.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
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

//nolint:paralleltest // deletes an object outside Terraform, a parallel test could otherwise take over its freed ID
func TestAccIndexerProxyHTTPResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "User"),
				Check: testAccCheckResourceDisappears("prowlarr_indexer_proxy_http.test", func(client *prowlarr.APIClient, id int32) (*http.Response, error) {
					return client.IndexerProxyAPI.DeleteIndexerProxy(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccIndexerProxyHTTPResourceConfig("resourceHTTPTest", "User"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("prowlarr_indexer_proxy_http.test", "id"),
				),
			},
		},
	})
}
