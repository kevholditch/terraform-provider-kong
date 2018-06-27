[![Build Status](https://travis-ci.org/kevholditch/terraform-provider-kong.svg?branch=master)](https://travis-ci.org/kevholditch/terraform-provider-kong)

Terraform Provider Kong
=======================
The Kong Terraform Provider tested against real Kong!


Notice
------
I have recently updated the provider to use `v1.0.0` of [gokong](http://github.com/kevholditch/gokong) this pulls in the changes to use pointers to all api fields.  If you update to the latest provider
be aware of this change.  Terraform may want to update some api resources as this fixes a bug where if you set a string from a value to `""` it will now be treated as empty string and not ignored.  If you
have set any of your api fields to empty string this will now be picked up.

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

Optionally you can configure Username and Password for BasicAuth:
```hcl
provider "kong" {
    kong_admin_uri  = "http://myKong:8001"
    kong_admin_username = "youruser"
    kong_admin_password = "yourpass"
}
```


You can use environment variables to set the provider properties instead.  The following table shows all of the config options, the corresponding environment variables and their property defaults if you do not set them.

| Provider property     | Env variable         | Default if not set    | Use                                                                             |
|:----------------------|:---------------------|:----------------------|:--------------------------------------------------------------------------------|
| kong_admin_uri        | KONG_ADMIN_ADDR      | http://localhost:8001 | The url of the kong admin api                                                   |
| kong_admin_username   | KONG_ADMIN_USERNAME  | not set               | Username for the kong admin api                                                 |
| kong_admin_password   | KONG_ADMIN_PASSWORD  | not set               | Password for the kong admin api                                                 |
| tls_skip_verify       | TLS_SKIP_VERIFY      | false                 | Whether to skip tls certificate verification for the kong api when using https  |
| kong_api_key          | KONG_API_KEY         | not set               | API key used to secure the kong admin API                                       |




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
    api_id 	= "${kong_api.api.id}"
    consumer_id = "${kong_consumer.plugin_consumer.id}"
    config      = {
        limits.sms.minute = 77
    }
}
```

### Configure plugins for a consumer
Some plugins allow you to configure them for a specific consumer for example the [jwt](https://getkong.org/plugins/jwt/#create-a-jwt-credential) and [key-auth](https://getkong.org/plugins/key-authentication/#create-an-api-key) plugins.
To configure a plugin for a consumer this terraform provider provides a generic way to do this for all plugins the "kong_consumer_plugin_config" resource.

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
    private_key  = "private key --- 456 ----"
}

resource "kong_sni" "sni" {
    name  	   = "www.example.com"
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
}
```


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
  * `name` - the name of the found upstream
  * `slots` - the number of slots on the found upstream
  * `order_list` - a list containing the slot order on the found upstream
  
  
# Contributing
I would love to get contributions to the project so please feel free to submit a PR.  To setup your dev station you need go and docker installed.

Once you have cloned the repository the `env TF_ACC=1 make` command will build the code and run all of the tests.  If they all pass then you are good to go!

If when you run the make command you get the following error:
```
gofmt needs running on the following files:
```
Then all you need to do is run `make goimports` this will reformat all of the code (I know awesome)!!

Please write tests for your new feature/bug fix, PRs will only be accepted with covering tests and where all tests pass.  If you want to start work on a feature feel free to open a PR early so we can discuss it or if you need help.
