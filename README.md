
Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!

**`v5.0.0` of the provider supports Terraform `0.12`**

IMPORTANT
------
The provider has been updated to support Kong `v2.X`, there were some breaking changes made between Kong `v1` and `v2`.  To use Kong `v1` use provider version `v6.X.X`.  That version will no longer be maintained.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.x
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

Usage
-----

First, install the desired [plugin release](https://github.com/kevholditch/terraform-provider-kong/releases) following Terraform's [Third-party plugin docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

To configure the provider:
```hcl
provider "kong" {
    kong_admin_uri = "http://myKong:8001"
}
```

Optionally you can configure Username and Password for BasicAuth:
```hcl
provider "kong" {
    kong_admin_uri  = "http://myKong:8001"
    kong_admin_username = "youruser"
    kong_admin_password = "yourpass"
}
```


You can use environment variables to set the provider properties instead.  The following table shows all of the config options, the corresponding environment variables and their property defaults if you do not set them.  When using the `kong_api_key` parameter ensure that the key name parameter in the key-auth plugin is set to `apikey`.

| Provider property              | Env variable                  | Default if not set    | Use                                                                             |
|:-------------------------------|:------------------------------|:----------------------|:--------------------------------------------------------------------------------|
| kong_admin_uri                 | KONG_ADMIN_ADDR               | http://localhost:8001 | The url of the kong admin api                                                   |
| kong_admin_username            | KONG_ADMIN_USERNAME           | not set               | Username for the kong admin api                                                 |
| kong_admin_password            | KONG_ADMIN_PASSWORD           | not set               | Password for the kong admin api                                                 |
| tls_skip_verify                | TLS_SKIP_VERIFY               | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| kong_api_key                   | KONG_API_KEY                  | not set               | API key used to secure the kong admin API                                       |
| kong_admin_token               | KONG_ADMIN_TOKEN              | not set               | API key used to secure the kong admin API in the Enterprise Edition             |
| strict_plugins_match           | STRICT_PLUGINS_MATCH          | false                 | Should plugins `config_json` field strictly match plugin configuration          |

# Documentation
For documentation on how to use the provider see the documentation on the [Hashicorp Terraform Registry for this provider](https://registry.terraform.io/providers/kevholditch/kong/latest/docs)

# Contributing
I would love to get contributions to the project so please feel free to submit a PR.  To setup your dev station you need go and docker installed.

Once you have cloned the repository the `env TF_ACC=1 make` command will build the code and run all of the tests.  If they all pass then you are good to go!

If when you run the make command you get the following error:
```
goimports needs running on the following files:
```
Then all you need to do is run `make goimports` this will reformat all of the code (I know awesome)!!

Please write tests for your new feature/bug fix, PRs will only be accepted with covering tests and where all tests pass.  If you want to start work on a feature feel free to open a PR early so we can discuss it or if you need help.
