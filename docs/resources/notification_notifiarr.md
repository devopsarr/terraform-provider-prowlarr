---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_notification_notifiarr Resource - terraform-provider-prowlarr"
subcategory: "Notifications"
description: |-
  Notification Notifiarr resource.
  For more information refer to Notification https://wiki.servarr.com/prowlarr/settings#connect and Notifiarr https://wiki.servarr.com/prowlarr/supported#notifiarr.
---

# prowlarr_notification_notifiarr (Resource)

<!-- subcategory:Notifications -->
Notification Notifiarr resource.
For more information refer to [Notification](https://wiki.servarr.com/prowlarr/settings#connect) and [Notifiarr](https://wiki.servarr.com/prowlarr/supported#notifiarr).

## Example Usage

```terraform
resource "prowlarr_notification_notifiarr" "example" {
  on_health_issue       = false
  on_application_update = false

  include_health_warnings = false
  name                    = "Example"

  api_key = "Token"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) API key.
- `name` (String) NotificationNotifiarr name.

### Optional

- `include_health_warnings` (Boolean) Include health warnings.
- `include_manual_grabs` (Boolean) Include manual grab flag.
- `on_application_update` (Boolean) On application update flag.
- `on_grab` (Boolean) On release grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_health_restored` (Boolean) On health restored flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_notification_notifiarr.example 1
```
