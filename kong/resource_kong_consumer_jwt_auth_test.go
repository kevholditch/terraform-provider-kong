package kong

import (
	"context"
	"fmt"
	"github.com/kong/go-kong/kong"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccJWTAuth(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJWTAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateJWTAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJWTAuthExists("kong_consumer_jwt_auth.consumer_jwt_config"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "algorithm", "HS256"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "key", "my_key"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "secret", "my_secret"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "rsa_public_key", "foo"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "tags.0", "foo"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "tags.1", "bar"),
				),
			},
			{
				Config: testUpdateJWTAuthConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckJWTAuthExists("kong_consumer_jwt_auth.consumer_jwt_config"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "algorithm", "HS256"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "key", "updated_key"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "secret", "updated_secret"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "rsa_public_key", "bar"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_consumer_jwt_auth.consumer_jwt_config", "tags.0", "foo"),
				),
			},
		},
	})
}

func TestAccJWTAuthImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckJWTAuthDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateJWTAuthConfig,
			},
			{
				ResourceName:      "kong_consumer_jwt_auth.consumer_jwt_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckJWTAuthDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.JWTAuths

	resources := getResourcesByType("kong_consumer_jwt_auth", state)

	if len(resources) != 1 {
		return fmt.Errorf("expecting only 1 jwt auth resource found %v", len(resources))
	}

	id, err := splitConsumerID(resources[0].Primary.ID)
	jwtAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get jwt auth by id: %v", err)
	}

	if jwtAuth != nil {
		return fmt.Errorf("jwt auth %s still exists, %+v", id.ID, jwtAuth)
	}

	return nil
}

func testAccCheckJWTAuthExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.JWTAuths
		id, err := splitConsumerID(rs.Primary.ID)

		jwtAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

		if err != nil {
			return err
		}

		if jwtAuth == nil {
			return fmt.Errorf("jwtAuth with id %v not found", id.ID)
		}

		return nil
	}
}

const testCreateJWTAuthConfig = `
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

resource "kong_consumer_jwt_auth" "consumer_jwt_config" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	algorithm      = "HS256"
	key            = "my_key"
	rsa_public_key = "foo"
	secret         = "my_secret"
    tags           = ["foo", "bar"]
}
`
const testUpdateJWTAuthConfig = `
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

resource "kong_consumer_jwt_auth" "consumer_jwt_config" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	algorithm      = "HS256"
	key            = "updated_key"
	rsa_public_key = "bar"
	secret         = "updated_secret"
	tags           = ["foo"]
}
`
