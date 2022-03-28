package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccConsumerOAuth2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsumerOAuth2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateConsumerOAuth2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerOAuth2Exists("kong_consumer_oauth2.consumer_oauth2"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "name", "test_application"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "client_id", "client_id"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "client_secret", "client_secret"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.0", "https://asdf.com/callback"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.1", "https://test.com/callback"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "tags.0", "myTag"),
				),
			},
			{
				Config: testUpdateConsumerOAuth2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsumerOAuth2Exists("kong_consumer_oauth2.consumer_oauth2"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "name", "test_application_updated"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "client_id", "client_id_updated"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "client_secret", "client_secret_updated"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.0", "https://asdf.com/callback"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "redirect_uris.1", "https://test.cl/callback"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "tags.0", "myTag"),
					resource.TestCheckResourceAttr("kong_consumer_oauth2.consumer_oauth2", "tags.1", "anotherTag"),
				),
			},
		},
	})
}

func testAccCheckConsumerOAuth2Destroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.Oauth2Credentials

	resources := getResourcesByType("kong_consumer_oauth2", state)

	if len(resources) != 1 {
		return fmt.Errorf("expecting only 1 consumer oauth2 resource found %v", len(resources))
	}

	id, _ := splitConsumerID(resources[0].Primary.ID)
	ConsumerOAuth2, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get consumer oauth2 by id: %v", err)
	}

	if ConsumerOAuth2 != nil {
		return fmt.Errorf("oauth2 %s still exists, %+v", id.ID, ConsumerOAuth2)
	}

	return nil
}

func testAccCheckConsumerOAuth2Exists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.Oauth2Credentials
		id, _ := splitConsumerID(rs.Primary.ID)

		ConsumerOAuth2, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

		if err != nil {
			return err
		}

		if ConsumerOAuth2 == nil {
			return fmt.Errorf("ConsumerOAuth2 with id %v not found", id.ID)
		}

		return nil
	}
}

const testCreateConsumerOAuth2Config = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "oauth2_plugin" {
	name = "oauth2"
	config_json = <<EOT
	{
		"global_credentials": true,
		"enable_password_grant": true,
		"token_expiration": 180,
		"refresh_token_ttl": 180,
		"provision_key": "testprovisionkey"
	}
EOT
}

resource "kong_consumer_oauth2" "consumer_oauth2" {
	name          = "test_application"
	consumer_id   = "${kong_consumer.my_consumer.id}"
	client_id     = "client_id"
	client_secret = "client_secret"
	redirect_uris = ["https://asdf.com/callback", "https://test.com/callback"]
	tags          = ["myTag"]
}
`
const testUpdateConsumerOAuth2Config = `
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "oauth2_plugin" {
	name = "oauth2"
	config_json = <<EOT
	{
		"global_credentials": true,
		"enable_password_grant": true,
		"token_expiration": 180,
		"refresh_token_ttl": 180,
		"provision_key": "testprovisionkey"
	}
EOT
}

resource "kong_consumer_oauth2" "consumer_oauth2" {
	name          = "test_application_updated"
	consumer_id   = "${kong_consumer.my_consumer.id}"
	client_id     = "client_id_updated"
	client_secret = "client_secret_updated"
	redirect_uris = ["https://asdf.com/callback", "https://test.cl/callback"]
	tags          = ["myTag", "anotherTag"]
}
`
