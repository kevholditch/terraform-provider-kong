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
