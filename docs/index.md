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

You can use environment variables to set the provider properties instead.  The following table shows all the config options, the corresponding environment variables, and their property defaults if you do not set them.  When using the `kong_api_key` parameter ensure that the key name parameter in the key-auth plugin is set to `apikey`.

| Provider property              | Env variable                  | Default if not set    | Use                                                                             |
|:-------------------------------|:------------------------------|:----------------------|:--------------------------------------------------------------------------------|
| kong_admin_uri                 | KONG_ADMIN_ADDR               | http://localhost:8001 | The url of the kong admin api                                                   |
| kong_admin_username            | KONG_ADMIN_USERNAME           | not set               | Username for the kong admin api                                                 |
| kong_admin_password            | KONG_ADMIN_PASSWORD           | not set               | Password for the kong admin api                                                 |
| tls_skip_verify                | TLS_SKIP_VERIFY               | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| kong_api_key                   | KONG_API_KEY                  | not set               | API key used to secure the kong admin API                                       |
| kong_admin_token               | KONG_ADMIN_TOKEN              | not set               | API key used to secure the kong admin API in the Enterprise Edition             |
| strict_plugins_match           | STRICT_PLUGINS_MATCH          | false                 | Should plugins `config_json` field strictly match plugin configuration          |
