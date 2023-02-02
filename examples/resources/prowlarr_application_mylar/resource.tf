resource "prowlarr_application_mylar" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:8090"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [7030]
}