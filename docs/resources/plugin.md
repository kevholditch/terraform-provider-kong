# kong_plugin

The plugin resource maps directly onto the json for the API endpoint in Kong.  For more information on the parameters [see the Kong Api create documentation](https://docs.konghq.com/gateway-oss/2.5.x/admin-api/#plugin-object).
The `config_json` is passed through to the plugin to configure it as is.  

## Example Usage

```hcl
resource "kong_plugin" "rate_limit" {
  name = "rate-limiting"
  config_json = <<EOT
	{
		"second": 5,
		"hour" : 1000
	}
EOT
}
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

## Argument reference

* `plugin_name` - (Required) the name of the plugin you want to configure
* `consumer_id` - (Optional) the consumer id you want to configure the plugin for
* `service_id`  - (Optional) the service id that you want to configure the plugin for
* `route_id` - (Optional) the route id that you want to configure the plugin for
* `enabled` - (Optional) whether the plugin is enabled or not, use if you want to keep the plugin installed but disable it
* `config_json` - (Optional) this is the configuration json for how you want to configure the plugin.  The json is passed straight through to kong as is.  You can get the json config from the Kong documentation
page of the plugin you are configuring
* `tags` - (Optional) A list of strings associated with the Plugin for grouping and filtering

## Import

To import a plugin:

```shell
terraform import kong_plugin.<plugin_identifier> <plugin_id>
```
