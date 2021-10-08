package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccConsumerKeyAuth(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerKeyAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerKeyAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerKeyAuthExists("kong_consumer_key_auth.consumer_key_auth"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "key", "foo"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.0", "myTag"),
				),
			},
			{
				Config: testUpdateConsumerKeyAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerKeyAuthExists("kong_consumer_key_auth.consumer_key_auth"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "key", "foo_updated"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.0", "myTag"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.1", "anotherTag"),
				),
			},
		},
	})
}

func TestAccConsumerKeyAuthComputed(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerKeyAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerKeyAuthConfigKeyComputed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerKeyAuthExists("kong_consumer_key_auth.consumer_key_auth"),
					resource.TestCheckResourceAttrSet("kong_consumer_key_auth.consumer_key_auth", "key"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.0", "myTag"),
				),
			},
			{
				Config: testUpdateConsumerKeyAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerKeyAuthExists("kong_consumer_key_auth.consumer_key_auth"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "key", "foo_updated"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.0", "myTag"),
					resource.TestCheckResourceAttr("kong_consumer_key_auth.consumer_key_auth", "tags.1", "anotherTag"),
				),
			},
		},
	})
}

func testAccCheckConsumerKeyAuthDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.KeyAuths

	resources := getResourcesByType("kong_consumer_key_auth", state)

	if len(resources) != 1 {
		return fmt.Errorf("expecting only 1 consumer key auth resource found %v", len(resources))
	}

	id, err := splitConsumerID(resources[0].Primary.ID)
	ConsumerKeyAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get consumer auth by id: %v", err)
	}

	if ConsumerKeyAuth != nil {
		return fmt.Errorf("key auth %s still exists, %+v", id.ID, ConsumerKeyAuth)
	}

	return nil
}

func testAccCheckConsumerKeyAuthExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.KeyAuths
		id, err := splitConsumerID(rs.Primary.ID)

		ConsumerKeyAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

		if err != nil {
			return err
		}

		if ConsumerKeyAuth == nil {
			return fmt.Errorf("ConsumerKeyAuth with id %v not found", id.ID)
		}

		return nil
	}
}

const testCreateConsumerKeyAuthConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "key_auth_plugin" {
	name = "key-auth"
}

resource "kong_consumer_key_auth" "consumer_key_auth" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	key         = "foo"
	tags        = ["myTag"]
}
`
const testUpdateConsumerKeyAuthConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "key_auth_plugin" {
	name = "key-auth"
}

resource "kong_consumer_key_auth" "consumer_key_auth" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	key         = "foo_updated"
	tags        = ["myTag", "anotherTag"]
}
`
const testCreateConsumerKeyAuthConfigKeyComputed = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "key_auth_plugin" {
	name = "key-auth"
}

resource "kong_consumer_key_auth" "consumer_key_auth" {
	consumer_id = "${kong_consumer.my_consumer.id}"
	tags        = ["myTag"]
}
`
