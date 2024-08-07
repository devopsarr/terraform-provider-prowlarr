---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_download_clients Data Source - terraform-provider-prowlarr"
subcategory: "Download Clients"
description: |-
  List all available Download Clients ../resources/download_client.
---

# prowlarr_download_clients (Data Source)

<!-- subcategory:Download Clients -->
List all available [Download Clients](../resources/download_client).

## Example Usage

```terraform
data "prowlarr_download_clients" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `download_clients` (Attributes Set) Download Client list. (see [below for nested schema](#nestedatt--download_clients))
- `id` (String) The ID of this resource.

<a id="nestedatt--download_clients"></a>
### Nested Schema for `download_clients`

Read-Only:

- `add_paused` (Boolean) Add paused flag.
- `add_stopped` (Boolean) Add stopped flag.
- `additional_tags` (Set of Number) Additional tags, `0` TitleSlug, `1` Quality, `2` Language, `3` ReleaseGroup, `4` Year, `5` Indexer, `6` Network.
- `api_key` (String, Sensitive) API key.
- `api_url` (String) API URL.
- `app_id` (String) App ID.
- `app_token` (String, Sensitive) App Token.
- `categories` (Attributes Set) List of mapped categories. (see [below for nested schema](#nestedatt--download_clients--categories))
- `category` (String) Category.
- `config_contract` (String) DownloadClient configuration template.
- `destination` (String) Destination.
- `destination_directory` (String) Movie directory.
- `directory` (String) Directory.
- `enable` (Boolean) Enable flag.
- `field_tags` (Set of String) Field tags.
- `host` (String) host.
- `id` (Number) Download Client ID.
- `implementation` (String) DownloadClient implementation name.
- `initial_state` (Number) Initial state. `0` Start, `1` ForceStart, `2` Pause.
- `intial_state` (Number) Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.
- `item_priority` (Number) Priority. `0` Last, `1` First.
- `magnet_file_extension` (String) Magnet file extension.
- `name` (String) Download Client name.
- `nzb_folder` (String) NZB folder.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `post_im_tags` (Set of String) Post import tags.
- `priority` (Number) Priority.
- `protocol` (String) Protocol. Valid values are 'usenet' and 'torrent'.
- `read_only` (Boolean) Read only flag.
- `rpc_path` (String) RPC path.
- `save_magnet_files` (Boolean) Save magnet files flag.
- `secret_token` (String, Sensitive) Secret token.
- `start_on_add` (Boolean) Start on add flag.
- `station_directory` (String) Directory.
- `strm_folder` (String) STRM folder.
- `tags` (Set of Number) List of associated tags.
- `torrent_folder` (String) Torrent folder.
- `tv_imported_category` (String) TV imported category.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

<a id="nestedatt--download_clients--categories"></a>
### Nested Schema for `download_clients.categories`

Read-Only:

- `categories` (Set of Number) List of categories.
- `name` (String) Name of client category.
