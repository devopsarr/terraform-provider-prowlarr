resource "prowlarr_notification_simplepush" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  key   = "Token"
  event = "ringtone:default"
}