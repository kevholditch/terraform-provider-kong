Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Usage
---------------------

```
# For example, restrict template version in 0.1.x
provider "template" {
  version = "~> 0.1"
}
```


To configure the provider:
```hcl
provider "kong" {
    kong_admin_uri = "http://myKong:8001"
}
```

By convention the provider will first check the env variable `KONG_ADMIN_ADDR` if that variable is not set then it will default to `http://localhost:8001` if
you do not provide a provider block as above.

To create an api:
```hcl
resource "kong_api" "api" {
    name 	= "TestApi"
    hosts   = [ "example.com" ]
    uris 	= [ "/example" ]
    methods = [ "GET", "POST" ]
    upstream_url = "http://localhost:4140"
    strip_uri = false
    preserve_host = false
    retries = 3
    upstream_connect_timeout = 60000
    upstream_send_timeout = 30000
    upstream_read_timeout = 10000
    https_only = false
    http_if_terminated = false
}
```
The api resource maps directly onto the json for creating an API in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/0.11.x/admin-api/#add-api).