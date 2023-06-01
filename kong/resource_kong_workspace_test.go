package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccKongWorkspace(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testWorkspaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongWorkspaceExists("kong_workspace.workspace"),
					resource.TestCheckResourceAttr("kong_workspace.workspace", "name", "myworkspace"),
				),
			},
			{
				Config: testUpdateWorkspaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongWorkspaceExists("kong_workspace.workspace"),
					resource.TestCheckResourceAttr("kong_workspace.workspace", "name", "yourworkspace"),
				),
			},
		},
	})
}

func TestAccKongWorkspaceImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testWorkspaceConfig,
			},
			{
				ResourceName:      "kong_workspace.workspace",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongWorkspaceDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.Workspaces

	workspaces := getResourcesByType("kong_workspace", state)

	if len(workspaces) > 1 {
		return fmt.Errorf("expecting max 1 workspace resource, found %v", len(workspaces))
	}

	if len(workspaces) == 0 {
		return nil
	}

	response, err := client.ListAll(context.Background())

	if err != nil {
		return fmt.Errorf("error thrown when trying to list workspaces: %v", err)
	}

	if response != nil {
		for _, element := range response {
			if *element.ID == workspaces[0].Primary.ID {
				return fmt.Errorf("workspace %s still exists, %+v", workspaces[0].Primary.ID, response)
			}
		}
	}

	return nil
}

func testAccCheckKongWorkspaceExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.Workspaces
		workspaces, err := client.ListAll(context.Background())

		if !kong.IsNotFoundErr(err) && err != nil {
			return err
		}

		if workspaces == nil {
			return fmt.Errorf("workspace with id %v not found", rs.Primary.ID)
		}

		if len(workspaces) != 2 {
			return fmt.Errorf("expected two workspaces (default & just-created), found %v", len(workspaces))
		}

		return nil
	}
}

const testWorkspaceConfig = `
resource "kong_workspace" "workspace" {
	name			= "myworkspace"
}
`
const testUpdateWorkspaceConfig = `
resource "kong_workspace" "workspace" {
	name			= "yourworkspace"
}
`
