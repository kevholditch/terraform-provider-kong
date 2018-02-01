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
		CheckDestroy: testAccCheckKongConsumerPluginConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.config"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "plugin_name", "jwt"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "consumer_id", "123"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "config_json", `<<EOF
						{
							"key": "e71829c351aa4242c2719cbfbe671c09",
							"secret": "a36c3049b36249a3c9f8891cb127243c"
						}
						EOF`),
				),
			},
			{
				Config: testUpdateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.config"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "plugin_name", "jwt"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "consumer_id", "123"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.config", "config_json", `<<EOF
						{
							"key": "a36c3049b36249a3c9f8891cb127243c",
							"secret": "e71829c351aa4242c2719cbfbe671c09"
						}
						EOF`),
				),
			},
		},
	})
}

func TestAccKongConsumerPluginConfigImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerPluginConfigDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateConsumerPluginConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_consumer_plugin_config.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongConsumerPluginConfigDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	configs := getResourcesByType("kong_consumer_plugin_config", state)

	if len(configs) != 1 {
		return fmt.Errorf("expecting only 1 consumer config resource found %v", len(configs))
	}

	response, err := client.Consumers().GetPluginConfig("123", "jwt", configs[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get consumer by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("consumer %s still exists, %+v", configs[0].Primary.ID, response)
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

		plugin, err := client.Consumers().GetPluginConfig("123", "jwt", rs.Primary.ID)

		if err != nil {
			return err
		}

		if plugin == nil {
			return fmt.Errorf("consumer with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateConsumerPluginConfig = `
resource "kong_consumer_plugin_config" "config" {
	plugin_name  = "jwt"
	consumer_id = "123"
	{
		"key": "e71829c351aa4242c2719cbfbe671c09",
		"secret": "a36c3049b36249a3c9f8891cb127243c"
	}
}
`
const testUpdateConsumerPluginConfig = `
resource "kong_consumer_plugin_config" "config" {
	plugin_name  = "jwt"
	consumer_id = "123"
	config_json = <<EOF
	{
		"key": "a36c3049b36249a3c9f8891cb127243c",
		"secret": "e71829c351aa4242c2719cbfbe671c09"
	}
	EOF
}
`
