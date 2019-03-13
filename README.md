[![Build Status](https://travis-ci.org/kevholditch/terraform-provider-kong.svg?branch=master)](https://travis-ci.org/kevholditch/terraform-provider-kong)

Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!


Notice
------
**CURRENTLY NOT COMPATIBLE WITH KONG 1.0.0  - Im working on support for it**

I have recently updated the provider to use `v1.0.0` of [gokong](http://github.com/kevholditch/gokong) this pulls in the changes to use pointers to all api fields.  If you update to the latest provider
be aware of this change.  Terraform may want to update some api resources as this fixes a bug where if you set a string from a value to `""` it will now be treated as empty string and not ignored.  If you
have set any of your api fields to empty string this will now be picked up.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

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

| Provider property     | Env variable         | Default if not set    | Use                                                                             |
|:----------------------|:---------------------|:----------------------|:--------------------------------------------------------------------------------|
| kong_admin_uri        | KONG_ADMIN_ADDR      | http://localhost:8001 | The url of the kong admin api                                                   |
| kong_admin_username   | KONG_ADMIN_USERNAME  | not set               | Username for the kong admin api                                                 |
| kong_admin_password   | KONG_ADMIN_PASSWORD  | not set               | Password for the kong admin api                                                 |
| tls_skip_verify       | TLS_SKIP_VERIFY      | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| kong_api_key          | KONG_API_KEY         | not set               | API key used to secure the kong admin API                                       |
| kong_admin_token      | KONG_ADMIN_TOKEN     | not set               | API key used to secure the kong admin API in the Enterprise Edition             |



# Resources

## Services
```hcl
resource "kong_service" "service" {
	name     	= "test"
	protocol 	= "http"
	host     	= "test.org"
	port     	= 8080
	path     	= "/mypath"
	retries  	= 5
	connect_timeout = 1000
	write_timeout 	= 2000
	read_timeout  	= 3000

}
```
The service resource maps directly onto the json for the service endpoint in Kong.  For more information on the parameters [see the Kong Service create documentation](https://getkong.org/docs/0.13.x/admin-api/#service-object).

To import a service:
```
terraform import kong_service.<service_identifier> <service_id>
```

## Routes
```hcl
resource "kong_route" "route" {
	protocols 	= [ "http", "https" ]
	methods 	= [ "GET", "POST" ]
	hosts 		= [ "example2.com" ]
	paths 		= [ "/test" ]
	strip_path 	= false
	preserve_host 	= true
	service_id 	= "${kong_service.service.id}"
}

```
The route resource maps directly onto the json for the route endpoint in Kong.  For more information on the parameters [see the Kong Route create documentation](https://getkong.org/docs/0.13.x/admin-api/#route-object).

To import a route:
```
terraform import kong_route.<route_identifier> <route_id>
```

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
The api resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/0.13.x/admin-api/#api-object).

To import an API:
```
terraform import kong_api.<api_identifier> <api_id>
```

## Plugins
```hcl
resource "kong_plugin" "response_rate_limiting" {
    name   = "response-ratelimiting"
    config = {
        limits.sms.minute = 10
    }
}
```

The plugin resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/0.13.x/admin-api/#plugin-object).

To import a plugin:
```
terraform import kong_plugin.<plugin_identifier> <plugin_id>
```

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
    api_id 	= "${kong_api.api.id}"
    consumer_id = "${kong_consumer.plugin_consumer.id}"
    config      = {
        limits.sms.minute = 77
    }
}
```

### Configure plugins for a consumer
Some plugins allow you to configure them for a specific consumer for example the [jwt](https://getkong.org/plugins/jwt/#create-a-jwt-credential) and [key-auth](https://getkong.org/plugins/key-authentication/#create-an-api-key) plugins.
To configure a plugin for a consumer this terraform provider provides a generic way to do this for all plugins the `kong_consumer_plugin_config` resource.

```hcl
resource "kong_consumer_plugin_config" "consumer_jwt_config" {
    consumer_id = "876bf719-8f18-4ce5-cc9f-5b5af6c36007"
    plugin_name = "jwt"
    config_json = <<EOT
        {
	    "key": "my_key",
	    "secret": "my_secret"
	}
EOT
}
```

The example above shows configuring the jwt plugin for a consumer.

`consumer_id` is the consumer id you want to configure the plugin for
`plugin_name` the name of the plugin you want to configure
`config_json` this is the configuration json for how you want to configure the plugin.  The json is passed straight through to kong as is.  You can get the json config from the Kong documentation
page of the plugin you are configuring

Other plugins must be configured using key/value pairs, for example the [acl](https://getkong.org/plugins/acl/) plugin.  To update a plugin using key value pairs configure the "kong_consumer_plugin_config" resource.

```hcl
resource "kong_consumer_plugin_config" "consumer_acl_config" {
    consumer_id = "876bf719-8f18-4ce5-cc9f-5b5af6c36007"
    plugin_name = "acls"
    config      = {
        group = "your_acl_group"
    }
}
```

All parameters are the same as above except the `config` parameter.
`config` is a map of key/value pairs you wish to pass as the configuration.

#### NOTE:  You can only have either config or config_json configured, not both.


## Consumers
```hcl
resource "kong_consumer" "consumer" {
    username  = "User1"
    custom_id = "123"
}
```

The consumer resource maps directly onto the json for creating an Consumer in Kong.  For more information on the parameters [see the Kong Consumer create documentation](https://getkong.org/docs/0.13.x/admin-api/#consumer-object).

To import a consumer:
```
terraform import kong_consumer.<consumer_identifier> <consumer_id>
```

## Certificates
```hcl
resource "kong_certificate" "certificate" {
    certificate  = "public key --- 123 ----"
    private_key = "private key --- 456 ----"
}
```

`certificate` should be the public key of your certificate it is mapped to the `Cert` parameter on the Kong API.
`private_key` should be the private key of your certificate it is mapped to the `Key` parameter on the Kong API.

For more information on creating certificates in Kong [see their documentation](https://getkong.org/docs/0.13.x/admin-api/#certificate-object)

To import a certificate:
```
terraform import kong_certificate.<certifcate_identifier> <certificate_id>
```

## SNIs
```hcl
resource "kong_certificate" "certificate" {
    certificate  = "public key --- 123 ----"
    private_key  = "private key --- 456 ----"
}

resource "kong_sni" "sni" {
    name  	   = "www.example.com"
    certificate_id = "${kong_certificate.certificate.id}"
}
```
`name` is your domain you want to assign to the certificate
`certificate_id` is the id of a certificate

For more information on creating SNIs in Kong [see their documentaton](https://getkong.org/docs/0.13.x/admin-api/#sni-objects)

To import a SNI:
```
terraform import kong_sni.<sni_identifier> <sni_id>
```

## Upstreams
```hcl
resource "kong_upstream" "upstream" {
    name                 = "sample_upstream"
    slots                = 10
    hash_on              = "header"
    hash_fallback        = "consumer"
    hash_on_header       = "HeaderName"
    hash_fallback_header = "FallbackHeaderName"
    healthchecks         = {
        active = {
            http_path                = "/status"
            timeout                  = 10
            concurrency              = 20
            healthy = {
                successes = 1
                interval  = 5
                http_statuses = [200, 201]
            }
            unhealthy = {
                timeouts      = 7
                interval      = 3
                tcp_failures  = 1
                http_failures = 2
                http_statuses = [500, 501]
            }
        }
        passive = {
            healthy = {
                successes = 1
                http_statuses = [200, 201, 202]
            }
            unhealthy = {
                timeouts      = 3
                tcp_failures  = 5
                http_failures = 6
                http_statuses = [500, 501, 502]
            }
        }
    }
}
```

  * `name` is a hostname, which must be equal to the host of a Service.
  * `slots` is the number of slots in the load balancer algorithm (10-65536, defaults to 10000).
  * `hash_on` is a hashing input type: `none `(resulting in a weighted-round-robin scheme with no hashing), `consumer`, `ip`, `header`, or `cookie`. Defaults to `none`.
  * `hash_fallback` is a hashing input type if the primary `hash_on` does not return a hash (eg. header is missing, or no consumer identified). One of: `none`, `consumer`, `ip`, `header`, or `cookie`. Not available if `hash_on` is set to `cookie`. Defaults to `none`.
  * `hash_on_header` is a header name to take the value from as hash input. Only required when `hash_on` is set to `header`. Default `nil`.
  * `hash_fallback_header` is a header name to take the value from as hash input. Only required when `hash_fallback` is set to `header`. Default `nil`.
  * `healthchecks.active.timeout` is a socket timeout for active health checks (in seconds). Defaults to `1`.
  * `healthchecks.active.concurrency` is a number of targets to check concurrently in active health checks. Defaults to `10`.
  * `healthchecks.active.http_path` is a path to use in GET HTTP request to run as a probe on active health checks. Defaults to `/`.
  * `healthchecks.active.healthy.interval` is an interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed. Defaults to `0`.
  * `healthchecks.active.healthy.successes` is a number of successes in active probes (as defined by `healthchecks.active.healthy.http_statuses`) to consider a target healthy. Defaults to `0`.
  * `healthchecks.active.healthy.http_statuses` is an array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks. Defaults to `[200, 302]`.
  * `healthchecks.active.unhealthy.interval` is an interval between active health checks for unhealthy targets (in seconds). A value of zero indicates that active probes for unhealthy targets should not be performed. Defaults to `0`.
  * `healthchecks.active.unhealthy.tcp_failures` is a number of TCP failures in active probes to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.http_failures` is a number of HTTP failures in active probes (as defined by `healthchecks.active.unhealthy.http_statuses`) to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.timeouts` is a number of timeouts in active probes to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.http_statuses` is an array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks. Defaults to `[429, 404, 500, 501, 502, 503, 504, 505]`.
  * `healthchecks.passive.healthy.successes` is a Number of successes in proxied traffic (as defined by `healthchecks.passive.healthy.http_statuses`) to consider a target healthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.healthy.http_statuses` is an array of HTTP statuses which represent healthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308]`.
  * `healthchecks.passive.unhealthy.tcp_failures` is a number of TCP failures in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.http_failures` is a number of HTTP failures in proxied traffic (as defined by `healthchecks.passive.unhealthy.http_statuses`) to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.timeouts` is a number of timeouts in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.http_statuses` is an array of HTTP statuses which represent unhealthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[429, 500, 503]`.


# Data Sources
## APIs
To look up an existing api you can do so by using a filter:
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

  * `id` - the id of the API
  * `name` - the name of the API
  * `hosts` - a list of the hosts configured on the API
  * `uris` - a list of the uri prefixes for the API
  * `methods` - a list of the allowed methods on the API
  * `upstream_url` - the upstream url for the API
  * `strip_uri` - whether the API strips the matching prefix from the uri
  * `preserve_host` - whether the API forwards the host header onto the upstream service
  * `retries` - number of retries the API executes upon failure to the upstream service
  * `upstream_connect_timeout` - the timeout in milliseconds for establishing a connection to your upstream service
  * `upstream_send_timeout` - the timeout in milliseconds between two successive write operations for transmitting a request to your upstream service
  * `upstream_read_timeout` - the timeout in milliseconds between two successive read operations for transmitting a request to your upstream service
  * `https_only` - whether the API is served through HTTPS
  * `http_if_terminated` - whether the API considers the  X-Forwarded-Proto header when enforcing HTTPS only traffic

## Certificates
To look up an existing certificate:
```hcl
data "kong_certificate" "certificate_data_source" {
    filter = {
        id = "471c625a-4eba-4b78-985f-86cf54a2dc12"
    }
}
```
You can only find existing certificates by their id in Kong.  The following output parameters are returned:

  * `id` - the Kong id for the certificate
  * `certificate` - the public key of the certificate
  * `private_key` - the private key of the certificate

## Consumers
To look up an existing consumer:
```hcl
data "kong_consumer" "consumer_data_source" {
    filter = {
        id 	  = "8086a91b-cb5a-4e60-90b0-ca6650e82464"
	username  = "User777"
	custom_id = "123456"
    }
}
```
Each of the filter parameters are optional and they are combined for an AND search against all consumers.   The following output parameters are
returned:

  * `id` - the Kong id of the found consumer
  * `username` - the username of the found consumer
  * `custom_id` - the custom id of the found consumer

## Plugins
To look up an existing plugin:
```hcl
data "kong_plugin" "plugin_data_source" {
    filter = {
        id          = "f0e656af-ad53-4622-ac73-ffd46ae05289"
	name        = "response-ratelimiting"
	api_id      = "51694bcd-3c72-43b3-b414-a09bbf4e3c30"
	consumer_id = "88154fd2-7a0e-41b1-97ba-4a59ebe2cc39"
    }
}
```
Each of the filter parameters are optional and they are combined for an AND search against all plugins.  The following output parameters are returned:

  * `id` - the Kong id of the found plugin
  * `name` - the name of the found plugin
  * `api_id` - the API id the found plugin is associated with (might be empty if not associated with an API)
  * `consumer_id` - the consumer id the found plugin is associated with (might be empty if not associated with a consumer)
  * `enabled` - whether the plugin is enabled

## Upstreams
To lookup an existing upstream:
```hcl
data "kong_upstream" "upstream_data_source" {
    filter = {
        id   = "893a49a8-090f-421e-afce-ba70b02ce958"
	name = "TestUpstream"
    }
}
```
Each of the filter parameters are optional and they are combined for an AND search against all upstreams.  The following output parameters are returned:

  * `id` - the Kong id of the found upstream
  * `name` is a hostname like name that can be referenced in the upstream_url field of a service.
  * `slots` is the number of slots in the load balancer algorithm (10-65536, defaults to 10000).
  * `hash_on` is a hashing input type: `none `(resulting in a weighted-round-robin scheme with no hashing), `consumer`, `ip`, `header`, or `cookie`. Defaults to `none`.
  * `hash_fallback` is a hashing input type if the primary `hash_on` does not return a hash (eg. header is missing, or no consumer identified). One of: `none`, `consumer`, `ip`, `header`, or `cookie`. Not available if `hash_on` is set to `cookie`. Defaults to `none`.
  * `hash_on_header` is a header name to take the value from as hash input. Only required when `hash_on` is set to `header`. Default `nil`.
  * `hash_fallback_header` is a header name to take the value from as hash input. Only required when `hash_fallback` is set to `header`. Default `nil`.
  * `healthchecks.active.timeout` is a socket timeout for active health checks (in seconds). Defaults to `1`.
  * `healthchecks.active.concurrency` is a number of targets to check concurrently in active health checks. Defaults to `10`.
  * `healthchecks.active.http_path` is a path to use in GET HTTP request to run as a probe on active health checks. Defaults to `/`.
  * `healthchecks.active.healthy.interval` is an interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed. Defaults to `0`.
  * `healthchecks.active.healthy.successes` is a number of successes in active probes (as defined by `healthchecks.active.healthy.http_statuses`) to consider a target healthy. Defaults to `0`.
  * `healthchecks.active.healthy.http_statuses` is an array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks. Defaults to `[200, 302]`.
  * `healthchecks.active.unhealthy.interval` is an interval between active health checks for unhealthy targets (in seconds). A value of zero indicates that active probes for unhealthy targets should not be performed. Defaults to `0`.
  * `healthchecks.active.unhealthy.tcp_failures` is a number of TCP failures in active probes to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.http_failures` is a number of HTTP failures in active probes (as defined by `healthchecks.active.unhealthy.http_statuses`) to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.timeouts` is a number of timeouts in active probes to consider a target unhealthy. Defaults to `0`.
  * `healthchecks.active.unhealthy.http_statuses` is an array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks. Defaults to `[429, 404, 500, 501, 502, 503, 504, 505]`.
  * `healthchecks.passive.healthy.successes` is a Number of successes in proxied traffic (as defined by `healthchecks.passive.healthy.http_statuses`) to consider a target healthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.healthy.http_statuses` is an array of HTTP statuses which represent healthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308]`.
  * `healthchecks.passive.unhealthy.tcp_failures` is a number of TCP failures in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.http_failures` is a number of HTTP failures in proxied traffic (as defined by `healthchecks.passive.unhealthy.http_statuses`) to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.timeouts` is a number of timeouts in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
  * `healthchecks.passive.unhealthy.http_statuses` is an array of HTTP statuses which represent unhealthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[429, 500, 503]`.
  * `order_list` - a list containing the slot order on the found upstream

To import an upstream:
```
terraform import kong_upstream.<upstream_identifier> <upstream_id>
```

## Targets
```hcl
resource "kong_target" "target" {
    target  		= "sample_target:80"
    weight 	  	= 10
    upstream_id = "${kong_upstream.upstream.id}"
}
```
`target` is the target address (IP or hostname) and port. If omitted the port defaults to 8000.
`weight` is the weight this target gets within the upstream load balancer (0-1000, defaults to 100).
`upstream_id` is the id of the upstream to apply this target to.


To import a target use a combination of the upstream id and the target id as follows:
```
terraform import kong_target.<target_identifier> <upstream_id>/<target_id>
```

# Contributing
I would love to get contributions to the project so please feel free to submit a PR.  To setup your dev station you need go and docker installed.

Once you have cloned the repository the `env TF_ACC=1 make` command will build the code and run all of the tests.  If they all pass then you are good to go!

If when you run the make command you get the following error:
```
goimports needs running on the following files:
```
Then all you need to do is run `make goimports` this will reformat all of the code (I know awesome)!!

Please write tests for your new feature/bug fix, PRs will only be accepted with covering tests and where all tests pass.  If you want to start work on a feature feel free to open a PR early so we can discuss it or if you need help.
