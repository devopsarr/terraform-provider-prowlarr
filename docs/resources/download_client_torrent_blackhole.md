---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_download_client_torrent_blackhole Resource - terraform-provider-prowlarr"
subcategory: "Download Clients"
description: |-
  Download Client Torrent Blackhole resource.
  For more information refer to Download Client https://wiki.servarr.com/prowlarr/settings#download-clients and TorrentBlackhole https://wiki.servarr.com/prowlarr/supported#torrentblackhole.
---

# prowlarr_download_client_torrent_blackhole (Resource)

<!-- subcategory:Download Clients -->
Download Client Torrent Blackhole resource.
For more information refer to [Download Client](https://wiki.servarr.com/prowlarr/settings#download-clients) and [TorrentBlackhole](https://wiki.servarr.com/prowlarr/supported#torrentblackhole).

## Example Usage

```terraform
resource "prowlarr_download_client_torrent_blackhole" "example" {
  enable                = true
  priority              = 1
  name                  = "Example"
  magnet_file_extension = ".magnet"
  torrent_folder        = "/torrent/"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.
- `torrent_folder` (String) Torrent folder.

### Optional

- `enable` (Boolean) Enable flag.
- `magnet_file_extension` (String) Magnet file extension.
- `priority` (Number) Priority.
- `save_magnet_files` (Boolean) Save magnet files flag.
- `tags` (Set of Number) List of associated tags.

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
terraform import prowlarr_download_client_torrent_blackhole.example 1
```
