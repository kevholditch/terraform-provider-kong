# kong_plugin

The plugin resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://getkong.org/docs/1.0.x/admin-api/#plugin-object).
The `config_json` is passed through to the plugin to configure it as is.  Note that the old `config` property has been removed due to incompatibility issues with kong v1.0.0.
Having the `config_json` property gives you ultimate flexibility to configure the plugin.

## Example Usage

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

## Argument reference

`plugin_name` - (Required) the name of the plugin you want to configure
`consumer_id` - (Optional) is the consumer id you want to configure the plugin for
`config_json` - (Optional) this is the configuration json for how you want to configure the plugin.  The json is passed straight through to kong as is.  You can get the json config from the Kong documentation
page of the plugin you are configuring

## Import

To import a plugin:

```shell
terraform import kong_plugin.<plugin_identifier> <plugin_id>
```
