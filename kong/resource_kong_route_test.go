package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKongRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.route"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_route.route", "hosts.0", "example.com"),
					resource.TestCheckResourceAttr("kong_route.route", "paths.0", "/"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "true"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "false"),
				),
			},
			{
				Config: testUpdateRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.route"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.1", "https"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.1", "POST"),
					resource.TestCheckResourceAttr("kong_route.route", "hosts.0", "example2.com"),
					resource.TestCheckResourceAttr("kong_route.route", "paths.0", "/test"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "false"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "true"),
				),
			},
		},
	})
}

func TestAccKongRouteImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testImportRouteConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_route.route",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongRouteDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient

	routes := getResourcesByType("kong_route", state)

	if len(routes) != 1 {
		return fmt.Errorf("expecting only 1 route resource found %v", len(routes))
	}

	response, err := client.Routes().GetRoute(routes[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get route by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("route %s still exists, %+v", routes[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongRouteExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		route, err := testAccProvider.Meta().(*config).adminClient.Routes().GetRoute(rs.Primary.ID)

		if err != nil {
			return err
		}

		if route == nil {
			return fmt.Errorf("route with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateRouteConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "http" ]
	methods 		= [ "GET" ]
	hosts 			=	[ "example.com" ]
	paths 			= [ "/" ]
	strip_path 		= true
	preserve_host 	= false
	service_id  	= "${kong_service.service.id}"
}
`
const testUpdateRouteConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "http", "https" ]
	methods 		= [ "GET", "POST" ]
	hosts 			= [ "example2.com" ]
	paths 			= [ "/test" ]
	strip_path 		= false
	preserve_host 	= true
	service_id 		= "${kong_service.service.id}"
}
`
const testImportRouteConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "http" ]
	methods 		= [ "GET" ]
	hosts 			= [ "example.com" ]
	paths 			= [ "/" ]
	strip_path 		= true
	preserve_host 	= false
	service_id		= "${kong_service.service.id}"
}
`
