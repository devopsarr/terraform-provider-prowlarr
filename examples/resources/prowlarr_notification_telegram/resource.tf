resource "prowlarr_notification_telegram" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  bot_token = "Token"
  chat_id   = "ChatID01"
}