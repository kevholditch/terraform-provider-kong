[![Build Status](https://travis-ci.org/kevholditch/terraform-provider-kong.svg?branch=master)](https://travis-ci.org/kevholditch/terraform-provider-kong)

Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!

**`v5.0.0` of the provider supports Terraform `0.12`**


IMPORTANT
------
This provider now supports kong `v1.0.0` and onwards **ONLY** (from `v2.0.0` onwards of provider).  Since the release of Kong `v1.0.0` has many breaking changes (e.g. removing APIs) this provider is
no longer compatible with version of kong pre `v1.0.0`.  If you want to use the provider with versions of kong pre `v1.0.0` then please checkout branch `kong-pre-1.0.0` or
use a version of the provider `v1.9.2` or less.

Due to compatibility issues I have had to remove some of the properties on the resources.  Most notability for a plugin you can only configure it using the `config_json` property
the `config` property has been removed.  This is due to some internal changes that have been made to kong in `v1.0.0`.

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

| Provider property              | Env variable                  | Default if not set    | Use                                                                             |
|:-------------------------------|:------------------------------|:----------------------|:--------------------------------------------------------------------------------|
| kong_admin_uri                 | KONG_ADMIN_ADDR               | http://localhost:8001 | The url of the kong admin api                                                   |
| kong_admin_username            | KONG_ADMIN_USERNAME           | not set               | Username for the kong admin api                                                 |
| kong_admin_password            | KONG_ADMIN_PASSWORD           | not set               | Password for the kong admin api                                                 |
| tls_skip_verify                | TLS_SKIP_VERIFY               | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| kong_api_key                   | KONG_API_KEY                  | not set               | API key used to secure the kong admin API                                       |
| kong_admin_token               | KONG_ADMIN_TOKEN              | not set               | API key used to secure the kong admin API in the Enterprise Edition             |
| strict_plugins_match           | STRICT_PLUGINS_MATCH          | false                 | Should plugins `config_json` field strictly match plugin configuration          |



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
The service resource maps directly onto the json for the service endpoint in Kong.  For more information on the parameters [see the Kong Service create documentation](https://getkong.org/docs/1.0.x/admin-api/#service-object).

To import a service:
```
terraform import kong_service.<service_identifier> <service_id>
```

## Routes
```hcl
resource "kong_route" "route" {
    name            = "MyRoute"
	protocols 	    = [ "http", "https" ]
	methods 	    = [ "GET", "POST" ]
	hosts 		    = [ "example2.com" ]
	paths 		    = [ "/test" ]
	strip_path 	    = false
	preserve_host 	= true
	regex_priority 	= 1
	service_id 	    = "${kong_service.service.id}"
}

```
The route resource maps directly onto the json for the route endpoint in Kong.  For more information on the parameters [see the Kong Route create documentation](https://getkong.org/docs/1.0.x/admin-api/#route-object).

To create a tcp/tls route you set `sources` and `destinations` by repeating the corresponding element (`source` or `destination`) for each
 source or destination you want, for example:

```hcl

resource "kong_route" "route" {
	protocols 		= [ "tcp" ]
	strip_path 		= true
	preserve_host 	= false
	source {
		ip   = "192.168.1.1"
		port = 80
	}
	source {
		ip   = "192.168.1.2"
	}
	destination {
		ip 	 = "172.10.1.1"
		port = 81
	}
	snis			= ["foo.com"]
	service_id  	= "${kong_service.service.id}"
}
```

To import a route:
```
terraform import kong_route.<route_identifier> <route_id>
```


## Plugins
```hcl
resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"	
	config_json = <<EOT
	{
		"second": 5,
		"hour" : 1000
	}
EOT
```

The `config_json` is passed through to the plugin to configure it as is.  Note that the old `config` property has been removed due to incompatibility issues with kong v1.0.0.
Having the `config_json` property gives you ultimate flexibility to configure the plugin.

To apply a plugin to a consumer use the `consumer_id` property, for example:
```hcl
resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "567"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config_json = <<EOT
	{
		"second": 5,
		"hour" : 1000
	}
EOT
}
```

To apply a plugin to a service use the `service_id` property, for example:

```hcl
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	service_id = "${kong_service.service.id}"
	config_json = <<EOT
	{
		"second": 10,
		"hour" : 2000
	}
EOT
}
```

To apply a plugin to a route use the `route_id` property, for example:

```hcl
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	enabled     = true
	service_id = "${kong_service.service.id}"
	config_json = <<EOT
	{
		"second": 11,
		"hour" : 4000
	}
EOT
}
```


The plugin resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/1.0.x/admin-api/#plugin-object).

To import a plugin:
```
terraform import kong_plugin.<plugin_identifier> <plugin_id>
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

Here is another example using the [acl](https://getkong.org/plugins/acl/) plugin:  

```hcl
resource "kong_consumer_plugin_config" "consumer_acl_config" {
consumer_id = "876bf719-8f18-4ce5-cc9f-5b5af6c36007"
	plugin_name = "acls"
	config_json = <<EOT
	{
		"group": "your_acl_group"
	}
EOT
}
```

## Consumers
```hcl
resource "kong_consumer" "consumer" {
    username  = "User1"
    custom_id = "123"
}
```

The consumer resource maps directly onto the json for creating an Consumer in Kong.  For more information on the parameters [see the Kong Consumer create documentation](https://getkong.org/docs/1.0.x/admin-api/#consumer-object).

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

For more information on creating certificates in Kong [see their documentation](https://getkong.org/docs/1.0.x/admin-api/#certificate-object)

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

For more information on creating SNIs in Kong [see their documentaton](https://getkong.org/docs/1.0.x/admin-api/#sni-objects)

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
    hash_fallback        = "cookie"
    hash_on_header       = "HeaderName"
    hash_fallback_header = "FallbackHeaderName"
    hash_on_cookie       = "CookieName"
    hash_on_cookie_path  = "/path"
    healthchecks {
        active {
            type                     = "https"
            http_path                = "/status"
            timeout                  = 10
            concurrency              = 20
            https_verify_certificate = false
            https_sni                = "some.domain.com"
            healthy {
                successes = 1
                interval  = 5
                http_statuses = [200, 201]
            }
            unhealthy {
                timeouts      = 7
                interval      = 3
                tcp_failures  = 1
                http_failures = 2
                http_statuses = [500, 501]
            }
        }
        passive {
            type    = "https"
            healthy {
                successes = 1
                http_statuses = [200, 201, 202]
            }
            unhealthy {
                timeouts      = 3
                tcp_failures  = 5
                http_failures = 6
                http_statuses = [500, 501, 502]
            }
        }
    }
}
```

- `name` is a hostname, which must be equal to the host of a Service.
- `slots` is the number of slots in the load balancer algorithm (10-65536, defaults to 10000).
- `hash_on` is a hashing input type: `none `(resulting in a weighted-round-robin scheme with no hashing), `consumer`, `ip`, `header`, or `cookie`. Defaults to `none`.
- `hash_fallback` is a hashing input type if the primary `hash_on` does not return a hash (eg. header is missing, or no consumer identified). One of: `none`, `consumer`, `ip`, `header`, or `cookie`. Not available if `hash_on` is set to `cookie`. Defaults to `none`.
- `hash_on_header` is a header name to take the value from as hash input. Only required when `hash_on` is set to `header`. Default `nil`.
- `hash_fallback_header` is a header name to take the value from as hash input. Only required when `hash_fallback` is set to `header`. Default `nil`.
- `hash_on_cookie` is a cookie name to take the value from as hash input. Only required when `hash_on` or `hash_fallback` is set to `cookie`. If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response. Default `nil`.
- `hash_on_cookie_path` is a cookie path to set in the response headers. Only required when `hash_on` or `hash_fallback` is set to `cookie`. Defaults to `/`.
- `healthchecks.active.type` is a active health check type. HTTP or HTTPS, or just attempt a TCP connection. Possible values are `tcp`, `http` or `https`. Defaults to `http`.
- `healthchecks.active.timeout` is a socket timeout for active health checks (in seconds). Defaults to `1`.
- `healthchecks.active.concurrency` is a number of targets to check concurrently in active health checks. Defaults to `10`.
- `healthchecks.active.http_path` is a path to use in GET HTTP request to run as a probe on active health checks. Defaults to `/`.
- `healthchecks.active.https_verify_certificate` check the validity of the SSL certificate of the remote host when performing active health checks using HTTPS. Defaults to `true`.
- `healthchecks.active.https_sni` is the hostname to use as an SNI (Server Name Identification) when performing active health checks using HTTPS. This is particularly useful when Targets are configured using IPs, so that the target host’s certificate can be verified with the proper SNI. Default `nil`.
- `healthchecks.active.healthy.interval` is an interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed. Defaults to `0`.
- `healthchecks.active.healthy.successes` is a number of successes in active probes (as defined by `healthchecks.active.healthy.http_statuses`) to consider a target healthy. Defaults to `0`.
- `healthchecks.active.healthy.http_statuses` is an array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks. Defaults to `[200, 302]`.
- `healthchecks.active.unhealthy.interval` is an interval between active health checks for unhealthy targets (in seconds). A value of zero indicates that active probes for unhealthy targets should not be performed. Defaults to `0`.
- `healthchecks.active.unhealthy.tcp_failures` is a number of TCP failures in active probes to consider a target unhealthy. Defaults to `0`.
- `healthchecks.active.unhealthy.http_failures` is a number of HTTP failures in active probes (as defined by `healthchecks.active.unhealthy.http_statuses`) to consider a target unhealthy. Defaults to `0`.
- `healthchecks.active.unhealthy.timeouts` is a number of timeouts in active probes to consider a target unhealthy. Defaults to `0`.
- `healthchecks.active.unhealthy.http_statuses` is an array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks. Defaults to `[429, 404, 500, 501, 502, 503, 504, 505]`.
- `healthchecks.passive.type` is a passive health check type. Interpreting HTTP/HTTPS statuses, or just check for TCP connection success. Possible values are `tcp`, `http` or `https` (in passive checks, `http` and `https` options are equivalent.). Defaults to `http`.
- `healthchecks.passive.healthy.successes` is a Number of successes in proxied traffic (as defined by `healthchecks.passive.healthy.http_statuses`) to consider a target healthy, as observed by passive health checks. Defaults to `0`.
- `healthchecks.passive.healthy.http_statuses` is an array of HTTP statuses which represent healthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308]`.
- `healthchecks.passive.unhealthy.tcp_failures` is a number of TCP failures in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
- `healthchecks.passive.unhealthy.http_failures` is a number of HTTP failures in proxied traffic (as defined by `healthchecks.passive.unhealthy.http_statuses`) to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
- `healthchecks.passive.unhealthy.timeouts` is a number of timeouts in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
- `healthchecks.passive.unhealthy.http_statuses` is an array of HTTP statuses which represent unhealthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[429, 500, 503]`.

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
