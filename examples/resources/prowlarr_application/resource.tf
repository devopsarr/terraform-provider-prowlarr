resource "prowlarr_application" "example" {
  name            = "Example"
  sync_level      = "disabled"
  implementation  = "Lidarr"
  config_contract = "LidarrSettings"
  base_url        = "http://localhost:8686"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [3000, 3010, 3030]
}