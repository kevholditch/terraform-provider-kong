package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
	"testing"
)

func TestAccKongCertificate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongCertificateExists("kong_certificate.certificate"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "certificate", "public key --- 123 ----"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "private_key", "private key --- 456 ----"),
				),
			},
			{
				Config: testUpdateCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongCertificateExists("kong_certificate.certificate"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "certificate", "public key --- 789 ----"),
					resource.TestCheckResourceAttr("kong_certificate.certificate", "private_key", "private key --- 321 ----"),
				),
			},
		},
	})
}

func testAccCheckKongCertificateDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	certificates := getResourcesByType("kong_certificate", state)

	if len(certificates) != 1 {
		return fmt.Errorf("expecting only 1 certificate resource found %v", len(certificates))
	}

	response, err := client.Certificates().GetById(certificates[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get certificate by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("certificate %s still exists, %+v", certificates[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongCertificateExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*gokong.KongAdminClient).Certificates().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("certificate with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateCertificateConfig = `
resource "kong_certificate" "certificate" {
	certificate  = "public key --- 123 ----"
	private_key = "private key --- 456 ----"
}
`
const testUpdateCertificateConfig = `
resource "kong_certificate" "certificate" {
	certificate  = "public key --- 789 ----"
	private_key = "private key --- 321 ----"
}
`
