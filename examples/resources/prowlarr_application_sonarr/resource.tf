resource "prowlarr_application_sonarr" "example" {
  name                  = "Example"
  sync_level            = "addOnly"
  base_url              = "http://localhost:8989"
  prowlarr_url          = "http://localhost:9696"
  api_key               = "APIKey"
  sync_categories       = [5000, 5010, 5030]
  anime_sync_categories = [5070]
}