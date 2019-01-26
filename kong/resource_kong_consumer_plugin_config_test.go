package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
)

func TestAccKongConsumerPluginConfig(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerPluginConfig,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_jwt_config"),
					//resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "plugin_name", "jwt"),
					//resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", `{"algorithm":"HS256","key":"my_key","secret":"my_secret"}`),
				),
			},
			{
				Config: testUpdateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_jwt_config"),
					//resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "plugin_name", "jwt"),
					//resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", `{"algorithm":"HS256","key":"updated_key","secret":"updated_secret"}`),
				),
			},
		},
	})
}

//
//func TestAccKongConsumerPluginConfigKV(t *testing.T) {
//
//	resource.Test(t, resource.TestCase{
//		Providers:    testAccProviders,
//		CheckDestroy: testAccCheckKongConsumerPluginConfig,
//		Steps: []resource.TestStep{
//			{
//				Config: testCreateConsumerPluginConfigKV,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_acl_config"),
//					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_acl_config", "plugin_name", "acls"),
//					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_acl_config", "config.group", "nginx"),
//				),
//			},
//			{
//				Config: testUpdateConsumerPluginConfigKV,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_acl_config"),
//					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_acl_config", "plugin_name", "acls"),
//					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_acl_config", "config.group", "apache"),
//				),
//			},
//		},
//	})
//}

func TestAccKongConsumerPluginConfigImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerPluginConfig,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testImportConsumerPluginConfigKV,
			},

			resource.TestStep{
				ResourceName:      "kong_consumer_plugin_config.consumer_acl_config",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckKongConsumerPluginConfig(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	consumerPluginConfigs := getResourcesByType("kong_consumer_plugin_config", state)

	if len(consumerPluginConfigs) != 1 {
		return fmt.Errorf("expecting only 1 consumer plugin config resource found %v", len(consumerPluginConfigs))
	}

	idFields, err := splitIdIntoFields(consumerPluginConfigs[0].Primary.ID)

	if err != nil {
		return err
	}

	response, err := client.Consumers().GetPluginConfig(idFields.consumerId, idFields.pluginName, idFields.id)

	if err != nil {
		return fmt.Errorf("error calling get consumer plugin config by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("consumer plugin config %s still exists, %+v", consumerPluginConfigs[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongConsumerPluginConfigExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*gokong.KongAdminClient)

		idFields, err := splitIdIntoFields(rs.Primary.ID)

		if err != nil {
			return err
		}

		consumerPluginConfig, err := client.Consumers().GetPluginConfig(idFields.consumerId, idFields.pluginName, idFields.id)

		if err != nil {
			return err
		}

		if consumerPluginConfig == nil {
			return fmt.Errorf("consumer plugin config with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateConsumerPluginConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "jwt_plugin" {
	name        = "jwt"
	config 		= {
		claims_to_verify = ["exp"]
	}
}

resource "kong_consumer_plugin_config" "consumer_jwt_config" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	plugin_name = "jwt"
	config_json = <<EOT
		{
			"algorithm": "HS256",
			"key": "my_key",
			"secret": "my_secret"
		}
EOT
}
`

const testUpdateConsumerPluginConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "jwt_plugin" {
	name        = "jwt"
	config 		= {
		claims_to_verify = ["exp"]
	}
}

resource "kong_consumer_plugin_config" "consumer_jwt_config" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	plugin_name = "jwt"
	config_json = <<EOT
		{
			"algorithm": "HS256",
			"key": "updated_key",
			"secret": "updated_secret"
		}
EOT
}
`

const testCreateConsumerPluginConfigKV = `
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

resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	api_id      = "${kong_api.api.id}"
	config      = {
		whitelist = "nginx"
	}
}

resource "kong_consumer_plugin_config" "consumer_acl_config" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	plugin_name = "acls"
	config = {
		group = "nginx"
	}
}
`

const testUpdateConsumerPluginConfigKV = `
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

resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	api_id      = "${kong_api.api.id}"
	config = {
		whitelist = "apache"
	}
}

resource "kong_consumer_plugin_config" "consumer_acl_config" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	plugin_name = "acls"
	config      = {
		group = "apache"
	}
}
`

const testImportConsumerPluginConfigKV = `
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

resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	api_id      = "${kong_api.api.id}"
	config = {
		whitelist = "apache"
	}
}

resource "kong_consumer_plugin_config" "consumer_acl_config" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	plugin_name = "acls"
	config      = {
		group = "apache"
	}
}
`
