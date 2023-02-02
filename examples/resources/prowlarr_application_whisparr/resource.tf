resource "prowlarr_application_whisparr" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:6969"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [6000, 6010, 6030]
}