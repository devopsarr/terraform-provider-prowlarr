resource "prowlarr_notification_gotify" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  server    = "http://gotify-server.net"
  app_token = "Token"
  priority  = 5
}