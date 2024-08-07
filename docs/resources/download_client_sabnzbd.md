---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_download_client_sabnzbd Resource - terraform-provider-prowlarr"
subcategory: "Download Clients"
description: |-
  Download Client Sabnzbd resource.
  For more information refer to Download Client https://wiki.servarr.com/prowlarr/settings#download-clients and Sabnzbd https://wiki.servarr.com/prowlarr/supported#sabnzbd.
---

# prowlarr_download_client_sabnzbd (Resource)

<!-- subcategory:Download Clients -->
Download Client Sabnzbd resource.
For more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [Sabnzbd](https://wiki.servarr.com/prowlarr/supported#sabnzbd).

## Example Usage

```terraform
resource "prowlarr_download_client_sabnzbd" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "sabnzbd"
  url_base = "/sabnzbd/"
  port     = 9091
  api_key  = "test"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `api_key` (String, Sensitive) API key.
- `category` (String) Category.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `item_priority` (Number) Recent Movie priority. `-100` Default, `-2` Paused, `-1` Low, `0` Normal, `1` High, `2` Force.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `tags` (Set of Number) List of associated tags.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `categories` (Attributes Set) List of mapped categories. (see [below for nested schema](#nestedatt--categories))
- `id` (Number) Download Client ID.

<a id="nestedatt--categories"></a>
### Nested Schema for `categories`

Optional:

- `categories` (Set of Number) List of categories.
- `name` (String) Name of client category.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_download_client_sabnzbd.example 1
```
