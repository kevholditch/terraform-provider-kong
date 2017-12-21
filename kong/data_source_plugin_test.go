package kong

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccDataSourceKongPlugin(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testPluginDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_plugin.plugin_data_source", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("data.kong_plugin.plugin_data_source", "enabled", "true"),
				),
			},
		},
	})
}

const testPluginDataSourceConfig = `
resource "kong_api" "my_api" {
	name 	                 = "TestApi"
  	hosts                    = [ "example.com" ]
	uris 	                 = [ "/example" ]
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

resource "kong_consumer" "my_consumer" {
	username  = "PluginUser"
	custom_id = "111"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	api_id 		= "${kong_api.my_api.id}"
	consumer_id = "${kong_consumer.my_consumer.id}"
	config 		= {
		limits.sms.minute = 11
	}
}

data "kong_plugin" "plugin_data_source" {
	filter = {
		id          = "${kong_plugin.rate_limit.id}"
		name        = "response-ratelimiting"
		api_id      = "${kong_api.my_api.id}"
		consumer_id = "${kong_consumer.my_consumer.id}"
	}
}
`
