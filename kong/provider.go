package kong

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
	"os"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kong_admin_uri": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_ADDR", "http://localhost:8001"),
				Description: "The address of the kong admin url e.g. http://localhost:8001",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_api": resourceKongApi(),
		},

		ConfigureFunc: providerConfigure,
	}
}


func envDefaultFuncWithDefault(key string, defaultValue string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(key); v != "" {
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			return v, nil
		}
		return defaultValue, nil
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &gokong.Config{
		HostAddress: d.Get("kong_admin_uri").(string),
	}

	return gokong.NewClient(config), nil
}
