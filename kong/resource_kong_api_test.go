package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
	"testing"
)

func TestAccKongApi(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongApiExists("kong_api.api"),
					resource.TestCheckResourceAttr("kong_api.api", "name", "TestApi"),
					resource.TestCheckResourceAttr("kong_api.api", "hosts.0", "example.com"),
					resource.TestCheckResourceAttr("kong_api.api", "uris.0", "/example"),
					resource.TestCheckResourceAttr("kong_api.api", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_api.api", "methods.1", "POST"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_url", "http://localhost:4140"),
					resource.TestCheckResourceAttr("kong_api.api", "strip_uri", "false"),
					resource.TestCheckResourceAttr("kong_api.api", "preserve_host", "false"),
					resource.TestCheckResourceAttr("kong_api.api", "retries", "3"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_connect_timeout", "60000"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_send_timeout", "30000"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_read_timeout", "10000"),
					resource.TestCheckResourceAttr("kong_api.api", "https_only", "false"),
					resource.TestCheckResourceAttr("kong_api.api", "http_if_terminated", "false"),
				),
			},
			{
				Config: testUpdateApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongApiExists("kong_api.api"),
					resource.TestCheckResourceAttr("kong_api.api", "name", "MyApi"),
					resource.TestCheckResourceAttr("kong_api.api", "hosts.0", "different.com"),
					resource.TestCheckResourceAttr("kong_api.api", "uris.0", "/somedomain"),
					resource.TestCheckResourceAttr("kong_api.api", "methods.0", "PUT"),
					resource.TestCheckResourceAttr("kong_api.api", "methods.1", "PATCH"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_url", "http://localhost:4242"),
					resource.TestCheckResourceAttr("kong_api.api", "strip_uri", "true"),
					resource.TestCheckResourceAttr("kong_api.api", "preserve_host", "true"),
					resource.TestCheckResourceAttr("kong_api.api", "retries", "10"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_connect_timeout", "50000"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_send_timeout", "22000"),
					resource.TestCheckResourceAttr("kong_api.api", "upstream_read_timeout", "11000"),
					resource.TestCheckResourceAttr("kong_api.api", "https_only", "true"),
					resource.TestCheckResourceAttr("kong_api.api", "http_if_terminated", "true"),
				),
			},
		},
	})
}

func testAccCheckKongApiDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "kong_api" {
			continue
		}

		response, err := client.Apis().GetById(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error calling get api by id: %v", err)
		}

		if response != nil {
			return fmt.Errorf("api %s still exists, %+v", rs.Primary.ID, response)
		}
	}

	return nil
}

func testAccCheckKongApiExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*gokong.KongAdminClient).Apis().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("api with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateApiConfig = `
resource "kong_api" "api" {
	name 	= "TestApi"
  	hosts   = [ "example.com" ]
	uris 	= [ "/example" ]
	methods = [ "GET", "POST" ]
	upstream_url = "http://localhost:4140"
	strip_uri = false
	preserve_host = false
	retries = 3
	upstream_connect_timeout = 60000
	upstream_send_timeout = 30000
	upstream_read_timeout = 10000
	https_only = false
	http_if_terminated = false
}
`
const testUpdateApiConfig = `
resource "kong_api" "api" {
	name 	= "MyApi"
  	hosts   = [ "different.com" ]
	uris 	= [ "/somedomain" ]
	methods = [ "PUT", "PATCH" ]
	upstream_url = "http://localhost:4242"
	strip_uri = true
	preserve_host = true
	retries = 10
	upstream_connect_timeout = 50000
	upstream_send_timeout = 22000
	upstream_read_timeout = 11000
	https_only = true
	http_if_terminated = true
}
`
