resource "prowlarr_notification" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  implementation  = "CustomScript"
  config_contract = "CustomScriptSettings"

  path = "/scripts/prowlarr.sh"
}