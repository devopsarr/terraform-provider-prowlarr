package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"prowlarr": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("PROWLARR_URL"); v == "" {
		t.Skip("PROWLARR_URL must be set for acceptance tests")
	}

	if v := os.Getenv("PROWLARR_API_KEY"); v == "" {
		t.Skip("PROWLARR_API_KEY must be set for acceptance tests")
	}
}

const testUnauthorizedProvider = `
provider "prowlarr" {
	url = "http://localhost:9696"
	api_key = "ErrorAPIKey"
  }
`
