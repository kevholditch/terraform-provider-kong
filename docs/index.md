Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong (using Docker)!

Terraform provider tested to work against Kong 2.X.

Usage
-----

To configure the provider:
```hcl
provider "kong" {
    kong_admin_uri = "http://localhost:8001"
}
```

Optionally you can configure Username and Password for BasicAuth:
```hcl
provider "kong" {
    kong_admin_uri  = "http://localhost:8001"
    kong_admin_username = "youruser"
    kong_admin_password = "yourpass"
}
```

## Argument Reference

In addition to generic provider arguments (e.g. alias and version), the following arguments are supported in the Kong provider block:

* `kong_admin_uri` - (Required) The URI of the Kong admin API, can be sourced from the `KONG_ADMIN_ADDR` environment variable
* `kong_admin_username` - (Optional) The username for the Kong admin API if set, can be sourced from the `KONG_ADMIN_USERNAME` environment variable
* `kong_admin_password` - (Optional) The password for the Kong admin API if set, can be sourced from the `KONG_ADMIN_PASSWORD` environment variable
* `tls_skip_verify` - (Optional) Whether to skip TLS certificate verification for the kong api when using https, can be sourced from the `TLS_SKIP_VERIFY` environment variable
* `kong_api_key` - (Optional) API key used to secure the kong admin API, can be sourced from the `KONG_API_KEY` environment variable
* `kong_admin_token` - (Optional) API key used to secure the kong admin API in the Enterprise Edition, can be sourced from the `KONG_ADMIN_TOKEN` environment variable
* `kong_workspace` - (Optional) Workspace context (Enterprise Edition)
* `strict_plugins_match` - (Optional) Should plugins `config_json` field strictly match plugin configuration                               
              
