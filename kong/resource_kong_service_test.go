package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
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
					resource.TestCheckResourceAttr("kong_service.service", "port", "80"),
					resource.TestCheckResourceAttr("kong_service.service", "path", "/mypath"),
					resource.TestCheckResourceAttr("kong_service.service", "retries", "5"),
					resource.TestCheckResourceAttr("kong_service.service", "connect_timeout", "1000"),
					resource.TestCheckResourceAttr("kong_service.service", "write_timeout", "2000"),
					resource.TestCheckResourceAttr("kong_service.service", "read_timeout", "3000"),
					resource.TestCheckResourceAttr("kong_service.service", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_service.service", "tags.0", "foo"),
					resource.TestCheckResourceAttr("kong_service.service", "tags.1", "bar"),
				),
			},
			{
				Config: testUpdateServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "name", "test2"),
					resource.TestCheckResourceAttr("kong_service.service", "protocol", "https"),
					resource.TestCheckResourceAttr("kong_service.service", "host", "test2.org"),
					resource.TestCheckResourceAttr("kong_service.service", "path", "/"),
					resource.TestCheckResourceAttr("kong_service.service", "retries", "5"),
					resource.TestCheckResourceAttr("kong_service.service", "connect_timeout", "6000"),
					resource.TestCheckResourceAttr("kong_service.service", "write_timeout", "5000"),
					resource.TestCheckResourceAttr("kong_service.service", "read_timeout", "4000"),
					resource.TestCheckResourceAttr("kong_service.service", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_service.service", "tags.0", "foo"),
					resource.TestCheckResourceAttr("kong_service.service", "tls_verify", "true"),
					resource.TestCheckResourceAttr("kong_service.service", "tls_verify_depth", "2"),
				),
			},
		},
	})
}

func TestAccKongDefaultService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateServiceConfigZero,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "retries", "0"),
				),
			},
			{
				Config: testUpdateServiceConfigZero,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "retries", "0"),
				),
			},
		},
	})
}

func TestAccKongServiceWithClientCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testServiceWithClientCertificateConfig, testCert1, testKey1, testCert2, testKey2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongServiceExists("kong_service.service"),
					resource.TestCheckResourceAttr("kong_service.service", "name", "test"),
					resource.TestCheckResourceAttr("kong_service.service", "protocol", "https"),
					resource.TestCheckResourceAttr("kong_service.service", "host", "test.org"),
					func(s *terraform.State) error {
						module := s.RootModule()
						cert, ok := module.Resources["kong_certificate.certificate"]
						if !ok {
							return fmt.Errorf("could not find certificate resource")
						}

						service, ok := module.Resources["kong_service.service"]
						if !ok {
							return fmt.Errorf("could not find service resource")
						}

						v, ok := service.Primary.Attributes["client_certificate_id"]
						if !ok {
							return fmt.Errorf("could not find client_certificate_id property")
						}

						if v != cert.Primary.ID {
							return fmt.Errorf("client_certificate_id does not match certificate id")
						}
						return nil
					},
					func(s *terraform.State) error {
						module := s.RootModule()
						cert, ok := module.Resources["kong_certificate.ca"]
						if !ok {
							return fmt.Errorf("could not find ca certificate resource")
						}

						service, ok := module.Resources["kong_service.service"]
						if !ok {
							return fmt.Errorf("could not find service resource")
						}

						v, ok := service.Primary.Attributes["ca_certificate_ids.0"]
						if !ok {
							return fmt.Errorf("could not find ca_certificate_ids property")
						}

						if v != cert.Primary.ID {
							return fmt.Errorf("ca_certificate_ids does not match ca certificate id")
						}
						return nil
					},
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

	client := testAccProvider.Meta().(*config).adminClient.Services

	services := getResourcesByType("kong_service", state)

	if len(services) != 1 {
		return fmt.Errorf("expecting only 1 service resource found %v", len(services))
	}

	response, err := client.Get(context.Background(), kong.String(services[0].Primary.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
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

		client := testAccProvider.Meta().(*config).adminClient.Services
		service, err := client.Get(context.Background(), kong.String(rs.Primary.ID))

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
	name     		 = "test"
	protocol 		 = "http"
	host     		 = "test.org"
	path     		 = "/mypath"
	retries  		 = 5
	connect_timeout  = 1000
	write_timeout 	 = 2000
	read_timeout  	 = 3000
	tags             = ["foo", "bar"]
}
`
const testUpdateServiceConfig = `
resource "kong_service" "service" {
	name     		 = "test2"
	protocol 		 = "https"
	host     		 = "test2.org"
	port     		 = 8081
	path     		 = "/"
	connect_timeout  = 6000
	write_timeout 	 = 5000
	read_timeout  	 = 4000
	tags             = ["foo"]
	tls_verify       = true
	tls_verify_depth = 2
}
`
const testCreateServiceConfigZero = `
resource "kong_service" "service" {
	name     		= "test"
	protocol 		= "http"
	host     		= "test.org"
	path     		= "/mypath"
	retries  		= 0
	connect_timeout = 1000
	write_timeout 	= 2000
	read_timeout  	= 3000
}
`
const testUpdateServiceConfigZero = `
resource "kong_service" "service" {
	name     		= "test2"
	protocol 		= "https"
	host     		= "test2.org"
	port     		= 8081
	path     		= "/"
	connect_timeout = 6000
	write_timeout 	= 5000
	read_timeout  	= 4000
	retries         = 0
}
`

const testServiceWithClientCertificateConfig = `
resource "kong_certificate" "certificate" {
	certificate  = <<EOF
%s
EOF
	private_key =  <<EOF
%s
EOF
   snis			= ["foo.com"]
}

resource "kong_certificate" "ca" {
	certificate  = <<EOF
%s
EOF
	private_key =  <<EOF
%s
EOF
   snis			= ["ca.com"]
}

resource "kong_service" "service" {
	name                  = "test"
	protocol              = "https"
	host                  = "test.org"
	client_certificate_id = kong_certificate.certificate.id
	ca_certificate_ids    = [kong_certificate.ca.id]
}`

const testImportServiceConfig = `
resource "kong_service" "service" {
	name     		= "test"
	protocol 		= "http"
	host     		= "test.org"
	port     		= 8080
	path     		= "/mypath"
	retries  		= 5
	connect_timeout = 8000
	write_timeout 	= 9000
	read_timeout  	= 10000
}
`
