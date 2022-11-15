resource "prowlarr_notification_custom_script" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  path = "/scripts/prowlarr.sh"
}