---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "prowlarr_indexer_proxy Resource - terraform-provider-prowlarr"
subcategory: "Indexer Proxies"
description: |-
  Generic Indexer Proxy resource. When possible use a specific resource instead.
  For more information refer to Indexer Proxy https://wiki.servarr.com/prowlarr/settings#indexer-proxies.
---

# prowlarr_indexer_proxy (Resource)

<!-- subcategory:Indexer Proxies -->Generic Indexer Proxy resource. When possible use a specific resource instead.
For more information refer to [Indexer Proxy](https://wiki.servarr.com/prowlarr/settings#indexer-proxies).

## Example Usage

```terraform
resource "prowlarr_indexer_proxy" "example" {
  name            = "Example"
  implementation  = "FlareSolverr"
  config_contract = "FlareSolverrSettings"
  host            = "http://localhost:8191/"
  request_timeout = 60
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config_contract` (String) IndexerProxy configuration template.
- `implementation` (String) IndexerProxy implementation name.
- `name` (String) Indexer Proxy name.

### Optional

- `host` (String) host.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `request_timeout` (Number) Request timeout.
- `tags` (Set of Number) List of associated tags.
- `username` (String) Username.

### Read-Only

- `id` (Number) Indexer Proxy ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import prowlarr_indexer_proxy.example 1
```