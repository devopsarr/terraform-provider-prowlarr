resource "prowlarr_notification_webhook" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  url      = "https://example.webhook.com/example"
  method   = 1
  username = "exampleUser"
  password = "examplePass"
}