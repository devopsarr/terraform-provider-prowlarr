resource "prowlarr_notification_sendgrid" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key    = "APIkey"
  from       = "from_sendgrid@example.com"
  recipients = ["user1@example.com", "user2@example.com"]
}