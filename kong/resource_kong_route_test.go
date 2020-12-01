package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttr("kong_route.route", "name", "foo"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_route.route", "hosts.0", "example.com"),
					resource.TestCheckResourceAttr("kong_route.route", "paths.0", "/"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "true"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "false"),
					resource.TestCheckResourceAttr("kong_route.route", "regex_priority", "1"),
				),
			},
			{
				Config: testUpdateRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.route"),
					resource.TestCheckResourceAttr("kong_route.route", "name", "bar"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "http"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.1", "https"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.0", "GET"),
					resource.TestCheckResourceAttr("kong_route.route", "methods.1", "POST"),
					resource.TestCheckResourceAttr("kong_route.route", "hosts.0", "example2.com"),
					resource.TestCheckResourceAttr("kong_route.route", "paths.0", "/test"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "false"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "true"),
					resource.TestCheckResourceAttr("kong_route.route", "regex_priority", "2"),
				),
			},
		},
	})
}

func TestAccKongRouteWithSourcesAndDestinations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateRouteWithSourcesAndDestinationsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.route"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "tls"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "true"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "false"),
					resource.TestCheckResourceAttr("kong_route.route", "source.#", "2"),
					resource.TestCheckResourceAttr("kong_route.route", "destination.#", "1"),
					resource.TestCheckResourceAttr("kong_route.route", "snis.0", "foo.com"),
				),
			},
			{
				Config: testUpdateRouteWithSourcesAndDestinationsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongRouteExists("kong_route.route"),
					resource.TestCheckResourceAttr("kong_route.route", "protocols.0", "tls"),
					resource.TestCheckResourceAttr("kong_route.route", "strip_path", "true"),
					resource.TestCheckResourceAttr("kong_route.route", "preserve_host", "false"),
					resource.TestCheckResourceAttr("kong_route.route", "source.#", "1"),
					resource.TestCheckResourceAttr("kong_route.route", "destination.#", "2"),
					resource.TestCheckResourceAttr("kong_route.route", "snis.0", "bar.com"),
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

	response, err := client.Routes().GetById(routes[0].Primary.ID)

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

		route, err := testAccProvider.Meta().(*config).adminClient.Routes().GetById(rs.Primary.ID)

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
	name            = "foo"
	protocols 		= [ "http" ]
	methods 		= [ "GET" ]
	hosts 			=	[ "example.com" ]
	paths 			= [ "/" ]
	strip_path 		= true
	preserve_host 	= false
	regex_priority  = 1
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
	name            = "bar"
	protocols 		= [ "http", "https" ]
	methods 		= [ "GET", "POST" ]
	hosts 			= [ "example2.com" ]
	paths 			= [ "/test" ]
	strip_path 		= false
	preserve_host 	= true
	regex_priority  = 2
	service_id 		= "${kong_service.service.id}"
}
`

const testCreateRouteWithSourcesAndDestinationsConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "tls" ]
	strip_path 		= true
	preserve_host 	= false
	source {
		ip   = "192.168.1.1"
		port = 80 
	}
	source {
		ip   = "192.168.1.2"
	}
	destination {
		ip 	 = "172.10.1.1"
		port = 81
	}
	snis			= ["foo.com"]
	service_id  	= "${kong_service.service.id}"
}
`

const testUpdateRouteWithSourcesAndDestinationsConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "tls" ]
	strip_path 		= true
	preserve_host 	= false
	source {
		ip   = "192.168.1.1"
		port = 80 
	}
	destination {
		ip 	 = "172.10.1.1"
		port = 81
	}
	destination {
		ip 	 = "172.10.1.2"
		port = 82
	}
	snis			= ["bar.com"]
	service_id  	= "${kong_service.service.id}"
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
