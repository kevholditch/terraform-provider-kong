package kong

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "plugin_name", "jwt"),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"algorithm":"HS256"`))),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"key":"my_key"`))),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"secret":"my_secret"`))),
				),
			},
			{
				Config: testUpdateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_jwt_config"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "plugin_name", "jwt"),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"algorithm":"HS256"`))),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"key":"updated_key"`))),
					resource.TestMatchResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "config_json", getRegex(regexp.Compile(`"secret":"updated_secret"`))),
				),
			},
		},
	})
}

func TestAccKongConsumerPluginConfigImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerPluginConfig,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerPluginConfig,
			},
			{
				ResourceName:            "kong_consumer_plugin_config.consumer_jwt_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config_json"},
			},
		},
	})
}

func TestAccCheckKongConsumerPluginCreateAndRefreshFromNonExistentConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerPluginConfig,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerPluginConfigExists("kong_consumer_plugin_config.consumer_jwt_config"),
					resource.TestCheckResourceAttr("kong_consumer_plugin_config.consumer_jwt_config", "plugin_name", "jwt"),
					deleteConsumer("kong_consumer.my_consumer"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckKongConsumerPluginConfig(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient

	consumerPluginConfigs := getResourcesByType("kong_consumer_plugin_config", state)

	if len(consumerPluginConfigs) > 1 {
		return fmt.Errorf("expecting max 1 consumer plugin config resource. found %v", len(consumerPluginConfigs))
	}

	if len(consumerPluginConfigs) == 0 {
		return nil
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

		client := testAccProvider.Meta().(*config).adminClient

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

func deleteConsumer(resourceKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if err := testAccProvider.Meta().(*config).adminClient.Consumers().DeleteById(rs.Primary.ID); err != nil {
			return fmt.Errorf("could not delete kong consumer: %v", err)
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
	config_json = <<EOT
	{
		"claims_to_verify": ["exp"]
	}
EOT
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
	config_json = <<EOT
	{
		"claims_to_verify": ["exp"]
	}
EOT
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
