resource "prowlarr_notification_join" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key  = "Key"
  priority = 2
}