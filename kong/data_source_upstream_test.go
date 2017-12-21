package kong

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
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
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.0", "4"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.1", "3"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.2", "2"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.3", "1"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.4", "5"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.5", "6"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.6", "7"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.7", "8"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.8", "9"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "order_list.9", "10"),
				),
			},
		},
	})
}

const testUpstreamDataSourceConfig = `
resource "kong_upstream" "upstream_orderlist" {
	name  		= "TestUpstream"
	slots 		= 10
	order_list  = [ 4, 3, 2, 1, 5, 6, 7, 8, 9, 10 ]
}

data "kong_upstream" "upstream_data_source" {
	filter = {
		id   = "${kong_upstream.upstream_orderlist.id}"
		name = "TestUpstream"
	}
}
`
