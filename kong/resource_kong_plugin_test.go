package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
	"testing"
)

func TestAccKongPluginForAllConsumersAndApis(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForAllApisAndConsumersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.response_rate_limiting"),
					resource.TestCheckResourceAttr("kong_plugin.response_rate_limiting", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.response_rate_limiting", "config.limits.sms.minute", "10"),
				),
			},
			{
				Config: testUpdatePluginForAllApisAndConsumersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.response_rate_limiting"),
					resource.TestCheckResourceAttr("kong_plugin.response_rate_limiting", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.response_rate_limiting", "config.limits.sms.minute", "40"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificApi(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.basic_auth"),
					testAccCheckKongApiExists("kong_api.api"),
					testAccCheckForChildIdCorrect("kong_api.api", "kong_plugin.basic_auth", "api_id"),
					resource.TestCheckResourceAttr("kong_plugin.basic_auth", "name", "basic-auth"),
					resource.TestCheckResourceAttr("kong_plugin.basic_auth", "config.hide_credentials", "true"),
				),
			},
			{
				Config: testUpdatePluginForASpecificApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.basic_auth"),
					testAccCheckKongApiExists("kong_api.api"),
					testAccCheckForChildIdCorrect("kong_api.api", "kong_plugin.basic_auth", "api_id"),
					resource.TestCheckResourceAttr("kong_plugin.basic_auth", "name", "basic-auth"),
					resource.TestCheckResourceAttr("kong_plugin.basic_auth", "config.hide_credentials", "false"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckForChildIdCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "config.limits.sms.minute", "20"),
				),
			},
			{
				Config: testUpdatePluginForASpecificConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckForChildIdCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "config.limits.sms.minute", "11"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificApiAndConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificApiAndConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckKongApiExists("kong_api.api"),
					testAccCheckForChildIdCorrect("kong_api.api", "kong_plugin.rate_limit", "api_id"),
					testAccCheckForChildIdCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "config.limits.sms.minute", "77"),
				),
			},
			{
				Config: testUpdatePluginForASpecificApiAndConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckKongApiExists("kong_api.api"),
					testAccCheckForChildIdCorrect("kong_api.api", "kong_plugin.rate_limit", "api_id"),
					testAccCheckForChildIdCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "response-ratelimiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "config.limits.sms.minute", "23"),
				),
			},
		},
	})
}

func testAccCheckKongPluginDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	plugins := getResourcesByType("kong_plugin", state)

	if len(plugins) != 1 {
		return fmt.Errorf("expecting only 1 plugin resource found %v", len(plugins))
	}

	response, err := client.Plugins().GetById(plugins[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get plugin by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("plugin %s still exists, %+v", plugins[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckForChildIdCorrect(parentResource string, childResource string, childIdField string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[parentResource]

		if !ok {
			return fmt.Errorf("not found: %s", parentResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		parentId := rs.Primary.ID

		rs, ok = s.RootModule().Resources[childResource]

		if !ok {
			return fmt.Errorf("not found: %s", parentResource)
		}

		childId, ok := rs.Primary.Attributes[childIdField]

		if !ok {
			return fmt.Errorf("child id field %s not set on %s", childIdField, childResource)
		}

		if parentId != childId {
			return fmt.Errorf("expected %s id of %s to equal %s id field %s of %s", parentResource, parentId, childResource, childIdField, childId)
		}

		return nil
	}
}

func testAccCheckKongPluginExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*gokong.KongAdminClient).Plugins().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("plugin with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreatePluginForAllApisAndConsumersConfig = `
resource "kong_plugin" "response_rate_limiting" {
	name  = "response-ratelimiting"
	config = {
		limits.sms.minute = 10
	}
}
`
const testUpdatePluginForAllApisAndConsumersConfig = `
resource "kong_plugin" "response_rate_limiting" {
	name  = "response-ratelimiting"
	config = {
		limits.sms.minute = 40
	}
}
`
const testCreatePluginForASpecificApiConfig = `
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

resource "kong_plugin" "basic_auth" {
	name   = "basic-auth"
	api_id = "${kong_api.api.id}"
	config = {
		hide_credentials = "true"
	}
}
`

const testUpdatePluginForASpecificApiConfig = `
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

resource "kong_plugin" "basic_auth" {
	name   = "basic-auth"
	api_id = "${kong_api.api.id}"
	config = {
		hide_credentials = "false"
	}
}
`

const testCreatePluginForASpecificConsumerConfig = `
resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "567"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config 		= {
		limits.sms.minute = 20
	}
}
`

const testUpdatePluginForASpecificConsumerConfig = `
resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "567"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config 		= {
		limits.sms.minute = 11
	}
}
`

const testCreatePluginForASpecificApiAndConsumerConfig = `
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

resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "111"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	api_id 		= "${kong_api.api.id}"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config 		= {
		limits.sms.minute = 77
	}
}
`

const testUpdatePluginForASpecificApiAndConsumerConfig = `
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

resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "111"
}

resource "kong_plugin" "rate_limit" {
	name        = "response-ratelimiting"
	api_id 		= "${kong_api.api.id}"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config 		= {
		limits.sms.minute = 23
	}
}
`
