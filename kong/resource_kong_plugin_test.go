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
				),
			},
			{
				Config: testUpdatePluginForAllApisAndConsumersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.response_rate_limiting"),
					resource.TestCheckResourceAttr("kong_plugin.response_rate_limiting", "name", "response-ratelimiting"),
				),
			},
		},
	})
}

func testAccCheckKongPluginDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kong_api" {
			continue
		}

		response, err := client.Plugins().GetById(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error calling get plugin by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("plugin %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
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
