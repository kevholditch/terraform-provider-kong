package kong

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
)

func TestAccKongUpstream(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongUpstreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateUpstreamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "name", "MyUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "slots", "10"),
				),
			},
			{
				Config: testUpdateUpstreamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "name", "MyUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "slots", "20"),
				),
			},
		},
	})
}

func TestAccKongUpstreamWithOrderList(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongUpstreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateUpstreamWithOrderList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream_orderlist"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "name", "MyOrderListUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "slots", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.0", "3"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.1", "2"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.2", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.3", "4"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.4", "5"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.5", "6"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.6", "7"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.7", "8"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.8", "9"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.9", "10"),
				),
			},
			{
				Config: testUpdateUpstreamWithOrderList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream_orderlist"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "name", "MyOrderListUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "slots", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.0", "7"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.1", "8"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.2", "9"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.3", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.4", "3"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.5", "2"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.6", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.7", "4"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.8", "5"),
					resource.TestCheckResourceAttr("kong_upstream.upstream_orderlist", "order_list.9", "6"),
				),
			},
		},
	})
}

func TestAccKongUpstreamImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongUpstreamDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateUpstreamConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_upstream.upstream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongUpstreamDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*gokong.KongAdminClient)

	upstreams := getResourcesByType("kong_upstream", state)

	if len(upstreams) != 1 {
		return fmt.Errorf("expecting only 1 upstream resource found %v", len(upstreams))
	}

	response, err := client.Upstreams().GetById(upstreams[0].Primary.ID)

	if err != nil {
		return fmt.Errorf("error calling get upstream by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("upstream %s still exists, %+v", upstreams[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongUpstreamExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		api, err := testAccProvider.Meta().(*gokong.KongAdminClient).Upstreams().GetById(rs.Primary.ID)

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("upstream with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateUpstreamConfig = `
resource "kong_upstream" "upstream" {
	name  		= "MyUpstream"
	slots 		= 10
}
`
const testUpdateUpstreamConfig = `
resource "kong_upstream" "upstream" {
	name  		= "MyUpstream"
	slots 		= 20
}
`

const testCreateUpstreamWithOrderList = `
resource "kong_upstream" "upstream_orderlist" {
	name  		= "MyOrderListUpstream"
	slots 		= 10
	order_list  = [ 3, 2, 1, 4, 5, 6, 7, 8, 9, 10 ]
}
`

const testUpdateUpstreamWithOrderList = `
resource "kong_upstream" "upstream_orderlist" {
	name  		= "MyOrderListUpstream"
	slots 		= 10
	order_list  = [ 7, 8, 9, 10, 3, 2, 1, 4, 5, 6,  ]
}
`
