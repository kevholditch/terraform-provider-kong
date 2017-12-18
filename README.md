[![Build Status](https://travis-ci.org/kevholditch/terraform-provider-kong.svg?branch=master)](https://travis-ci.org/kevholditch/terraform-provider-kong)

Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Usage
-----

To configure the provider:
```hcl
provider "kong" {
    kong_admin_uri = "http://myKong:8001"
}
```

By convention the provider will first check the env variable `KONG_ADMIN_ADDR` if that variable is not set then it will default to `http://localhost:8001` if
you do not provide a provider block as above.

# Resources

## Apis
```hcl
resource "kong_api" "api" {
    name 	             = "TestApi"
    hosts                    = [ "example.com" ]
    uris 	             = [ "/example" ]
    methods                  = [ "GET", "POST" ]
    upstream_url             = "http://localhost:4140"
    strip_uri                = false
    preserve_host            = false
    retries                  = 3
    upstream_connect_timeout = 60000
    upstream_send_timeout    = 30000
    upstream_read_timeout    = 10000
    https_only               = false
    http_if_terminated       = false
}
```
The api resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/0.11.x/admin-api/#api-object).

## Plugins
```hcl
resource "kong_plugin" "response_rate_limiting" {
	name   = "response-ratelimiting"
	config = {
		limits.sms.minute = 10
	}
}
```

The plugin resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/0.11.x/admin-api/#plugin-object).

Here is a more complex example for creating a plugin for a consumer and an API:

```hcl
resource "kong_api" "api" {
    name 	             = "TestApi"
    hosts                    = [ "example.com" ]
    uris 	             = [ "/example" ]
    methods                  = [ "GET", "POST" ]
    upstream_url             = "http://localhost:4140"
    strip_uri                = false
    preserve_host            = false
    retries                  = 3
    upstream_connect_timeout = 60000
    upstream_send_timeout    = 30000
    upstream_read_timeout    = 10000
    https_only               = false
    http_if_terminated       = false
}

resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "111"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	api_id 	    = "${kong_api.api.id}"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config 	    = {
		limits.sms.minute = 77
	}
}
```


## Consumers
```hcl
resource "kong_consumer" "consumer" {
	username  = "User1"
	custom_id = "123"
}
```

The consumer resource maps directly onto the json for creating an Consumer in Kong.  For more information on the parameters [see the Kong Consumer create documentation](https://getkong.org/docs/0.11.x/admin-api/#consumer-object).

## Certificates
```hcl
resource "kong_certificate" "certificate" {
	certificate  = "public key --- 123 ----"
	private_key = "private key --- 456 ----"
}
```

`certificate` should be the public key of your certificate it is mapped to the `Cert` parameter on the Kong API.
`private_key` should be the private key of your certificate it is mapped to the `Key` parameter on the Kong API.

For more information on creating certificates in Kong [see their documentation](https://getkong.org/docs/0.11.x/admin-api/#certificate-object)

## SNIs
```hcl
resource "kong_certificate" "certificate" {
	certificate  = "public key --- 123 ----"
	private_key = "private key --- 456 ----"
}

resource "kong_sni" "sni" {
	name  		   = "www.example.com"
	certificate_id = "${kong_certificate.certificate.id}"
}
```
`name` is your domain you want to assign to the certificate
`certificate_id` is the id of a certificate

For more information on creating SNIs in Kong [see their documentaton](https://getkong.org/docs/0.11.x/admin-api/#sni-objects)

## Upstreams
```hcl
resource "kong_upstream" "upstream" {
	name  		= "sample_upstream"
	slots 		= 10
	order_list  = [ 3, 2, 1, 4, 5, 6, 7, 8, 9, 10 ]
}
```
`order_list` is optional if not supplied then one will be generated at random by kong and it will be set in the resource state.  For more
information on creating Upstreams in Kong [see their documentaton](https://getkong.org/docs/0.11.x/admin-api/#upstream-objects)

# Data Sources
## Api
To look up an existing api
```hcl
data "kong_api" "api_data_source" {
	filter = {
		id = "de539d26-97d2-4d5b-aaf9-628e51087d9c"
		name = "TestDataSourceApi"
		upstream_url = "http://localhost:4140"
	}
}
```
Each of the filter parameters are optional and they are combined for an AND search against all APIs.   The following output parameters are
returned:

`id` - the id of the API
`name` - the name of the API
`hosts` - a list of the hosts configured on the API
`uris` - a list of the uri prefixes for the API
`methods` - a list of the allowed methods on the API
`upstream_url` - the upstream url for the API
`strip_uri` - whether the API strips the matching prefix from the uri
`preserve_host` - whether the API forwards the host header onto the upstream service
`retries` - number of retries the API executes upon failure to the upstream service
`upstream_connect_timeout` - the timeout in milliseconds for establishing a connection to your upstream service
`upstream_send_timeout` - the timeout in milliseconds between two successive write operations for transmitting a request to your upstream service
`upstream_read_timeout` - the timeout in milliseconds between two successive read operations for transmitting a request to your upstream service
`https_only` - whether the API is served through HTTPS
`http_if_terminated` - whether the API considers the  X-Forwarded-Proto header when enforcing HTTPS only traffic
