resource "prowlarr_notification_ntfy" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  priority   = 1
  server_url = "https://ntfy.sh"
  username   = "User"
  password   = "Pass"
  topics     = ["Topic1234", "Topic4321"]
  field_tags = ["warning", "skull"]
}