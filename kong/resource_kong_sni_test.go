package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccKongSni(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongSniDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateSniConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongSniExists("kong_sni.sni"),
					resource.TestCheckResourceAttr("kong_sni.sni", "name", "www.example.com"),
				),
			},
			{
				Config: testUpdateSniConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongSniExists("kong_sni.sni"),
					resource.TestCheckResourceAttr("kong_sni.sni", "name", "www.example.com"),
				),
			},
		},
	})
}

func TestAccKongSniImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongSniDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateSniConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_sni.sni",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongSniDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient

	snis := getResourcesByType("kong_sni", state)

	if len(snis) != 1 {
		return fmt.Errorf("expecting only 1 sni resource found %v", len(snis))
	}

	response, err := client.Snis().GetByName(snis[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get sni by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("sni %s still exists, %+v", snis[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongSniExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*config).adminClient.Snis().GetByName(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("sni with name %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateSniConfig = `
resource "kong_certificate" "certificate1" {
	certificate  = "public key --- 123 ----"
	private_key = "private key --- 456 ----"
}

resource "kong_certificate" "certificate2" {
	certificate  = "public key --- 789 ----"
	private_key = "private key --- 321 ----"
}

resource "kong_sni" "sni" {
	name  		   = "www.example.com"
	certificate_id = "${kong_certificate.certificate1.id}"
}

`
const testUpdateSniConfig = `
resource "kong_certificate" "certificate1" {
	certificate  = "public key --- 123 ----"
	private_key = "private key --- 456 ----"
}

resource "kong_certificate" "certificate2" {
	certificate  = "public key --- 789 ----"
	private_key = "private key --- 321 ----"
}

resource "kong_sni" "sni" {
	name  		   = "www.example.com"
	certificate_id = "${kong_certificate.certificate2.id}"
}
`
