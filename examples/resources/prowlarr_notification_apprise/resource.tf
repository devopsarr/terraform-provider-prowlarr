resource "prowlarr_notification_apprise" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  base_url          = "http://localhost:8000"
  configuration_key = "ConfigKey"
  auth_username     = "User"
  auth_password     = "Pass"
  field_tags        = ["test", "test1"]
}