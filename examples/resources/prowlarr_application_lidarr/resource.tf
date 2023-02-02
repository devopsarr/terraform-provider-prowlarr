resource "prowlarr_application_lidarr" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:8686"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [3000, 3010, 3030]
}