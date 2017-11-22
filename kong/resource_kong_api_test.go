package kong

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
)

func TestAccKongApi_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateApiConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongApiExists("kong_api.api"),
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

		if err == nil {
			return fmt.Errorf("record %s still exists, %+v", rs.Primary.ID, response)
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
			return fmt.Errorf("no Record ID is set")
		}

		client := testAccProvider.Meta().(*gokong.KongAdminClient)

		_, err := client.Apis().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}


		return nil
	}
}

const testCreateApiConfig = `
resource "kong_api" "api" {
	name 	= "TestApi"
  	hosts   = ["example.com"]
	uris 	= ["/example"]
	methods = ["GET", "POST"]
	upstream_url = "http://localhost:4140"
	strip_url = false
	preserve_host = false
	retries = 3
	upstream_connect_timeout = 60000
	upstream_send_timeout = 30000
	upstream_read_timeout = 10000
	https_only = false
	http_if_terminated = false
}
`

