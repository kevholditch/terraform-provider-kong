package kong

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func getResourcesByType(resourceType string, state *terraform.State) []*terraform.ResourceState {

	var result []*terraform.ResourceState

	for _, rs := range state.RootModule().Resources {
		if rs.Type == resourceType {
			result = append(result, rs)
		}

	}

	return result
}

func String(v string) *string {
	return &v
}
