package kong

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
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
			"kong_admin_username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_USERNAME", ""),
				Description: "An basic auth user for kong admin",
			},
			"kong_admin_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_PASSWORD", ""),
				Description: "An basic auth password for kong admin",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_api":                    resourceKongApi(),
			"kong_certificate":            resourceKongCertificate(),
			"kong_consumer":               resourceKongConsumer(),
			"kong_consumer_plugin_config": resourceKongConsumerPluginConfig(),
			"kong_plugin":                 resourceKongPlugin(),
			"kong_sni":                    resourceKongSni(),
			"kong_upstream":               resourceKongUpstream(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"kong_api":         dataSourceKongApi(),
			"kong_certificate": dataSourceKongCertificate(),
			"kong_consumer":    dataSourceKongConsumer(),
			"kong_plugin":      dataSourceKongPlugin(),
			"kong_upstream":    dataSourceKongUpstream(),
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
		Username:    d.Get("kong_admin_username").(string),
		Password:    d.Get("kong_admin_password").(string),
	}

	return gokong.NewClient(config), nil
}
