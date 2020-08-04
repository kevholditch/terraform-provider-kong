package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKongConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.consumer"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "username", "User1"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "custom_id", "123"),
				),
			},
			{
				Config: testUpdateConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.consumer"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "username", "User2"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "custom_id", "456"),
				),
			},
		},
	})
}

func TestAccKongConsumerImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateConsumerConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_consumer.consumer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongConsumerDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient

	consumers := getResourcesByType("kong_consumer", state)

	if len(consumers) != 1 {
		return fmt.Errorf("expecting only 1 consumer resource found %v", len(consumers))
	}

	response, err := client.Consumers().GetById(consumers[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get consumer by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("consumer %s still exists, %+v", consumers[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongConsumerExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient

		api, err := client.Consumers().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("consumer with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateConsumerConfig = `
resource "kong_consumer" "consumer" {
	username  = "User1"
	custom_id = "123"
}
`
const testUpdateConsumerConfig = `
resource "kong_consumer" "consumer" {
	username  = "User2"
	custom_id = "456"
}
`
