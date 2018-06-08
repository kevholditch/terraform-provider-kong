package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKongCertificate(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCertificateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_certificate.certificate_data_source", "certificate", "public key --- 777 ----"),
					resource.TestCheckResourceAttr("data.kong_certificate.certificate_data_source", "private_key", "private key --- 888 ----"),
				),
			},
		},
	})
}

const testCertificateDataSourceConfig = `
resource "kong_certificate" "test_certificate" {
	certificate  = "public key --- 777 ----"
	private_key  = "private key --- 888 ----"
}

data "kong_certificate" "certificate_data_source" {
	filter = {
		id = "${kong_certificate.test_certificate.id}"
	}
}
`
