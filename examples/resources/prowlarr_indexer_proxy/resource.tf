resource "prowlarr_indexer_proxy" "example" {
  name            = "Example"
  implementation  = "FlareSolverr"
  config_contract = "FlareSolverrSettings"
  host            = "http://localhost:8191/"
  request_timeout = 60
}