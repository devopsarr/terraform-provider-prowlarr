resource "prowlarr_notification_twitter" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  access_token        = "Token"
  access_token_secret = "TokenSecret"
  consumer_key        = "Key"
  consumer_secret     = "Secret"
  mention             = "someone"
}