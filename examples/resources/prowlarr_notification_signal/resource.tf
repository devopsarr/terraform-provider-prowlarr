resource "prowlarr_notification_signal" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  auth_username = "User"
  auth_password = "Token"

  host          = "localhost"
  port          = 8080
  use_ssl       = true
  sender_number = "1234"
  receiver_id   = "4321"
}