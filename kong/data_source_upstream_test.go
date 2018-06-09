package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKongUpstream(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUpstreamDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "name", "TestUpstream"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "slots", "10"),
				),
			},
		},
	})
}

const testUpstreamDataSourceConfig = `
resource "kong_upstream" "upstream" {
	name  		= "TestUpstream"
	slots 		= 10
}

data "kong_upstream" "upstream_data_source" {
	filter = {
		id   = "${kong_upstream.upstream.id}"
		name = "TestUpstream"
	}
}
`
