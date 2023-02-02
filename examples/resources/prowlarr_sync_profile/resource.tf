resource "prowlarr_sync_profile" "example" {
  name                      = "Example"
  minimum_seeders           = 1
  enable_rss                = true
  enable_automatic_search   = true
  enable_interactive_search = true
}