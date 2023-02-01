resource "prowlarr_indexer_proxy_http" "example" {
  name     = "Example"
  host     = "localhost"
  port     = 8080
  username = "User"
  password = "Pass"
}