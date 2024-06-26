---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Prowlarr Provider"
subcategory: ""
description: |-
  The Prowlarr provider is used to interact with any Prowlarr https://prowlarr.com/ installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.
---

# Prowlarr Provider

The Prowlarr provider is used to interact with any [Prowlarr](https://prowlarr.com/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.

## Example Usage

```terraform
provider "prowlarr" {
  url     = "http://example.prowlarr.tv:8989"
  api_key = "APIkey-example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) API key for Prowlarr authentication. Can be specified via the `PROWLARR_API_KEY` environment variable.
- `extra_headers` (Attributes Set) Extra headers to be sent along with all Prowlarr requests. If this attribute is unset, it can be specified via environment variables following this pattern `PROWLARR_EXTRA_HEADER_${Header-Name}=${Header-Value}`. (see [below for nested schema](#nestedatt--extra_headers))
- `url` (String) Full Prowlarr URL with protocol and port (e.g. `https://test.prowlarr.audio:8686`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `PROWLARR_URL` environment variable.

<a id="nestedatt--extra_headers"></a>
### Nested Schema for `extra_headers`

Required:

- `name` (String) Header name.
- `value` (String) Header value.
