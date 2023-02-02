resource "prowlarr_application_radarr" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:7878"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [2000, 2010, 2030]
}