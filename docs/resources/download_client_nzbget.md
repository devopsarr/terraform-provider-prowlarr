---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_download_client_nzbget Resource - terraform-provider-prowlarr"
subcategory: "Download Clients"
description: |-
  Download Client NZBGet resource.
  For more information refer to Download Client https://wiki.servarr.com/prowlarr/settings#download-clients and NZBGet https://wiki.servarr.com/prowlarr/supported#nzbget.
---

# prowlarr_download_client_nzbget (Resource)

<!-- subcategory:Download Clients -->
Download Client NZBGet resource.
For more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [NZBGet](https://wiki.servarr.com/prowlarr/supported#nzbget).

## Example Usage

```terraform
resource "prowlarr_download_client_nzbget" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "nzbget"
  url_base = "/nzbget/"
  port     = 6789
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `category` (String) Category.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `item_priority` (Number) Recent Movie priority. `-100` VeryLow, `-50` Low, `0` Normal, `50` High, `100` VeryHigh, `900` Force.
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
terraform import prowlarr_download_client_nzbget.example 1
```
