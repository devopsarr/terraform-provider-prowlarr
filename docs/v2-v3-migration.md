# Migrating from provider V2 to V3

## Breaking Changes in Indexer Resource Configuration

### `fields` Handling Update

- **Change**: `prowlarr_indexer` resource now excludes fields with `type = info` from the state file
- **Impact**: `type = info` fields must be removed from Terraform configurations

#### Migration Steps

1. Update `devopsarr/prowlarr` provider
2. Remove `info` fields from `prowlarr_indexer.*.fields` like `info_tpp`, `info_flaresolverr`, etc
3. Run `terraform apply` to clean up the state

#### Example

Before:

```hcl
resource "prowlarr_indexer" "example" {
  name = "Example"
  implementation = "Cardigann"
  config_contract = "CardigannSettings"
  fields = [
    { name = "username", text_value = "example" },
    { name = "password", sensitive_value = "example" },
    { name = "info_tpp", text_value = "For best results, change the <b>Torrents per page:</b> setting to <b>100</b> on your account profile." }
  ]
}
```

After:

```hcl
resource "prowlarr_indexer" "example" {
  name = "Example"
  implementation = "Cardigann"
  config_contract = "CardigannSettings"
  fields = [
    { name = "username", text_value = "example" },
    { name = "password", sensitive_value = "example" }
  ]
}
```
