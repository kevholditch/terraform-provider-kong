package kong

import (
	"context"
	"fmt"
	"github.com/kong/go-kong/kong"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConsumerACL(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerACLExists("kong_consumer_acl.consumer_acl"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "group", "group1"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "tags.0", "myTag"),
				),
			},
			{
				Config: testUpdateConsumerACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerACLExists("kong_consumer_acl.consumer_acl"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "group", "group2"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "tags.0", "myTag"),
					resource.TestCheckResourceAttr("kong_consumer_acl.consumer_acl", "tags.1", "otherTag"),
				),
			},
		},
	})
}

func TestAccConsumerACLImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerACLConfig,
			},
			{
				ResourceName:      "kong_consumer_acl.consumer_acl",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckConsumerACLDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.ACLs

	resources := getResourcesByType("kong_consumer_acl", state)

	if len(resources) != 1 {
		return fmt.Errorf("expecting only 1 consumer acl resource found %v", len(resources))
	}

	id, err := splitConsumerID(resources[0].Primary.ID)
	ConsumerACL, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get consumer auth by id: %v", err)
	}

	if ConsumerACL != nil {
		return fmt.Errorf("jwt auth %s still exists, %+v", id.ID, ConsumerACL)
	}

	return nil
}

func testAccCheckConsumerACLExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.ACLs
		id, err := splitConsumerID(rs.Primary.ID)

		ConsumerACL, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

		if err != nil {
			return err
		}

		if ConsumerACL == nil {
			return fmt.Errorf("ConsumerACL with id %v not found", id.ID)
		}

		return nil
	}
}

const testCreateConsumerACLConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	config_json = <<EOT
	{
		"allow": ["group1", "group2"]
	}
EOT
}

resource "kong_consumer_acl" "consumer_acl" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	group          = "group1"
	tags           = ["myTag"]
}
`
const testUpdateConsumerACLConfig = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	config_json = <<EOT
	{
		"allow": ["group1", "group2"]
	}
EOT
}

resource "kong_consumer_acl" "consumer_acl" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	group          = "group2"
	tags           = ["myTag", "otherTag"]
}
`
