package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
)

func TestAccKongService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "name", "test"),
					resource.TestCheckResourceAttr("kong_service.service", "protocol", "http"),
					resource.TestCheckResourceAttr("kong_service.service", "host", "test.org"),
				),
			},
			{
				Config: testUpdateServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "name", "test2"),
					resource.TestCheckResourceAttr("kong_service.service", "protocol", "https"),
					resource.TestCheckResourceAttr("kong_service.service", "host", "test2.org"),
				),
			},
		},
	})
}

func TestAccKongServiceImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testImportServiceConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_service.service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongServiceDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	services := getResourcesByType("kong_service", state)

	if len(services) != 1 {
		return fmt.Errorf("expecting only 1 service resource found %v", len(services))
	}

	response, err := client.Services().GetServiceById(services[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get service by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("service %s still exists, %+v", services[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongServiceExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		service, err := testAccProvider.Meta().(*gokong.KongAdminClient).Services().GetServiceById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if service == nil {
			return fmt.Errorf("service with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateServiceConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}
`
const testUpdateServiceConfig = `
resource "kong_service" "service" {
	name     = "test2"
	protocol = "https"
	host     = "test2.org"
}
`
const testImportServiceConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}
`
