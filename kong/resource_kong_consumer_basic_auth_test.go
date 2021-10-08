package kong

import (
	"context"
	"fmt"
	"github.com/kong/go-kong/kong"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConsumerBasicAuth(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerBasicAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerBasicAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerBasicAuthExists("kong_consumer_basic_auth.consumer_basic_auth"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "username", "foo"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "tags.0", "myTag"),
				),
			},
			{
				Config: testUpdateConsumerBasicAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerBasicAuthExists("kong_consumer_basic_auth.consumer_basic_auth"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "username", "foo_updated"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "tags.0", "myTag"),
					resource.TestCheckResourceAttr("kong_consumer_basic_auth.consumer_basic_auth", "tags.1", "anotherTag"),
				),
			},
		},
	})
}

func testAccCheckConsumerBasicAuthDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.BasicAuths

	resources := getResourcesByType("kong_consumer_basic_auth", state)

	if len(resources) != 1 {
		return fmt.Errorf("expecting only 1 consumer basic auth resource found %v", len(resources))
	}

	id, err := splitConsumerID(resources[0].Primary.ID)
	ConsumerBasicAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get consumer auth by id: %v", err)
	}

	if ConsumerBasicAuth != nil {
		return fmt.Errorf("basic auth %s still exists, %+v", id.ID, ConsumerBasicAuth)
	}

	return nil
}

func testAccCheckConsumerBasicAuthExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.BasicAuths
		id, err := splitConsumerID(rs.Primary.ID)

		ConsumerBasicAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

		if err != nil {
			return err
		}

		if ConsumerBasicAuth == nil {
			return fmt.Errorf("ConsumerBasicAuth with id %v not found", id.ID)
		}

		return nil
	}
}

const testCreateConsumerBasicAuthConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "basic_auth_plugin" {
	name        = "basic-auth"
}

resource "kong_consumer_basic_auth" "consumer_basic_auth" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	username       = "foo"
	password       = "bar"
	tags           = ["myTag"]
}
`
const testUpdateConsumerBasicAuthConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "basic_auth_plugin" {
	name        = "basic-auth"
}

resource "kong_consumer_basic_auth" "consumer_basic_auth" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	username       = "foo_updated"
	password       = "bar_updated"
	tags           = ["myTag", "anotherTag"]
}
`
