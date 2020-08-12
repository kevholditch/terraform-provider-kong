package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKongApi(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testApiDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "name", "TestDataSourceApi"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "hosts.0", "example.com"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "uris.0", "/example"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "methods.0", "GET"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "methods.1", "POST"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "upstream_url", "http://localhost:4140"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "strip_uri", "false"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "preserve_host", "false"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "retries", "3"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "upstream_connect_timeout", "60000"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "upstream_send_timeout", "30000"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "upstream_read_timeout", "10000"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "https_only", "false"),
					resource.TestCheckResourceAttr("data.kong_api.api_data_source", "http_if_terminated", "false"),
				),
			},
		},
	})
}

const testApiDataSourceConfig = `
resource "kong_api" "my_test_api" {
	name 	= "TestDataSourceApi"
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

data "kong_api" "api_data_source" {
	filter {
		id = "${kong_api.my_test_api.id}"
		name = "TestDataSourceApi"
		upstream_url = "http://localhost:4140"
	}
}
`
