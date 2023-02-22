resource "prowlarr_notification_pushbullet" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key    = "Token"
  device_ids = ["test"]
}