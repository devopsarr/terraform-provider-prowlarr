resource "prowlarr_notification_email" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  server = "http://email-server.net"
  port   = 587
  from   = "from_email@example.com"
  to     = ["user1@example.com", "user2@example.com"]
}