package kong

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccKongTarget(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongTargetExists("kong_target.target"),
					resource.TestCheckResourceAttr("kong_target.target", "target", "mytarget:4000"),
					resource.TestCheckResourceAttr("kong_target.target", "weight", "100"),
					resource.TestCheckResourceAttr("kong_target.target", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_target.target", "tags.0", "a"),
					resource.TestCheckResourceAttr("kong_target.target", "tags.1", "b"),
				),
			},
			{
				Config: testUpdateTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongTargetExists("kong_target.target"),
					resource.TestCheckResourceAttr("kong_target.target", "target", "mytarget:4000"),
					resource.TestCheckResourceAttr("kong_target.target", "weight", "200"),
					resource.TestCheckResourceAttr("kong_target.target", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_target.target", "tags.0", "a"),
				),
			},
		},
	})
}

func TestAccKongTargetDelete(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongTargetExists("kong_target.target"),
					resource.TestCheckResourceAttr("kong_target.target", "target", "mytarget:4000"),
					resource.TestCheckResourceAttr("kong_target.target", "weight", "100"),
				),
			},
			{
				Config: testDeleteTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongTargetDoesNotExist("kong_target.target", "kong_upstream.upstream"),
				),
			},
		},
	})
}

func TestAccKongTargetCreateAndRefreshFromNonExistentUpstream(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateTargetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongTargetExists("kong_target.target"),
					resource.TestCheckResourceAttr("kong_target.target", "target", "mytarget:4000"),
					resource.TestCheckResourceAttr("kong_target.target", "weight", "100"),
					deleteUpstream("kong_upstream.upstream"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKongTargetImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateTargetConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_target.target",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongTargetDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.Targets

	targets := getResourcesByType("kong_target", state)

	if len(targets) > 1 {
		return fmt.Errorf("expecting max 1 target resource found %v", len(targets))
	}

	if len(targets) == 0 {
		return nil
	}

	response, _, _ := client.List(context.Background(), kong.String(targets[0].Primary.Attributes["upstream_id"]), nil)

	if response != nil {
		for _, element := range response {
			if *element.ID == targets[0].Primary.ID {
				return fmt.Errorf("target %s still exists, %+v", targets[0].Primary.ID, response)
			}
		}
	}

	return nil
}

func testAccCheckKongTargetExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		var ids = strings.Split(rs.Primary.ID, "/")
		client := testAccProvider.Meta().(*config).adminClient.Targets
		api, _, err := client.List(context.Background(), kong.String(ids[0]), nil)

		if !kong.IsNotFoundErr(err) && err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("target with id %v not found", rs.Primary.ID)
		}

		for _, element := range api {
			if *element.ID == ids[1] {
				break
			}

			return fmt.Errorf("target with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckKongTargetDoesNotExist(targetResourceKey string, upstreamResourceKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[targetResourceKey]

		if ok {
			return fmt.Errorf("Found target: %s", targetResourceKey)
		}

		rs, ok = s.RootModule().Resources[upstreamResourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", upstreamResourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no upstream ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.Targets
		targets, _, err := client.List(context.Background(), kong.String(rs.Primary.ID), nil)

		if len(targets) > 0 {
			return fmt.Errorf("expecting zero target resources found %v", len(targets))
		}

		if err != nil {
			return fmt.Errorf("error thrown when trying to read target: %v", err)
		}

		return nil
	}
}

func deleteUpstream(upstreamResourceKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[upstreamResourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", upstreamResourceKey)
		}

		client := testAccProvider.Meta().(*config).adminClient.Upstreams
		if err := client.Delete(context.Background(), kong.String(rs.Primary.ID)); err != nil {
			return fmt.Errorf("could not delete kong upstream: %v", err)
		}

		return nil
	}
}

const testCreateTargetConfig = `
resource "kong_upstream" "upstream" {
	name				= "MyUpstream"
	slots				= 10
}

resource "kong_target" "target" {
	target			= "mytarget:4000"
	weight			= 100
	upstream_id	    = "${kong_upstream.upstream.id}"
    tags            = ["a", "b"]
}
`
const testUpdateTargetConfig = `
resource "kong_upstream" "upstream" {
	name				= "MyUpstream"
	slots 			= 10
}

resource "kong_target" "target" {
	target			= "mytarget:4000"
	weight			= 200
	upstream_id  	= "${kong_upstream.upstream.id}"
	tags            = ["a"]
}
`
const testDeleteTargetConfig = `
resource "kong_upstream" "upstream" {
	name				= "MyUpstream"
	slots				= 10
}
`
