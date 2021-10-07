package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
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
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.0", "a"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.1", "b"),
				),
			},
			{
				Config: testUpdateConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.consumer"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "username", "User2"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "custom_id", "456"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.0", "a"),
				),
			},
		},
	})
}

func TestAccKongConsumerNilIDs(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongConsumerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerConfigNoCustomID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongConsumerExists("kong_consumer.consumer"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "username", "User3"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "custom_id", ""),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer.consumer", "tags.0", "c"),
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

	client := testAccProvider.Meta().(*config).adminClient.Consumers

	consumers := getResourcesByType("kong_consumer", state)

	if len(consumers) != 1 {
		return fmt.Errorf("expecting only 1 consumer resource found %v", len(consumers))
	}

	consumer, err := client.Get(context.Background(), kong.String(consumers[0].Primary.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get consumer by id: %v", err)
	}

	if consumer != nil {
		return fmt.Errorf("consumer %s still exists, %+v", consumers[0].Primary.ID, consumer)
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

		client := testAccProvider.Meta().(*config).adminClient.Consumers

		api, err := client.Get(context.Background(), kong.String(rs.Primary.ID))

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
    tags      = ["a", "b"]
}
`
const testUpdateConsumerConfig = `
resource "kong_consumer" "consumer" {
	username  = "User2"
	custom_id = "456"
    tags      = ["a"] 
}
`
const testCreateConsumerConfigNoCustomID = `
resource "kong_consumer" "consumer" {
	username = "User3"
	tags     = ["c"]
}
`
