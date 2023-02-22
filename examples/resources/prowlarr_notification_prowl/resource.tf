resource "prowlarr_notification_prowl" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "APIKey"
  priority = -2
}