resource "prowlarr_application_lazy_librarian" "example" {
  name            = "Example"
  sync_level      = "addOnly"
  base_url        = "http://localhost:5299"
  prowlarr_url    = "http://localhost:9696"
  api_key         = "APIKey"
  sync_categories = [7000, 7010, 7030]
}