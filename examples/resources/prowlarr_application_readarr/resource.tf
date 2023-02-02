resource "prowlarr_application_readarr" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:8787"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [7000, 7010, 7030]
}