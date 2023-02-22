resource "prowlarr_notification_slack" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  web_hook_url = "http://my.slack.com/test"
  username     = "user"
  channel      = "example-channel"
}