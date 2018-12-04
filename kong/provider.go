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
			"tls_skip_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("TLS_SKIP_VERIFY", "false"),
				Description: "Whether to skip tls verify for https kong api endpoint using self signed or untrusted certs",
			},
			"kong_api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_API_KEY", ""),
				Description: "API key for the kong api (if you have locked it down)",
			},
			"kong_admin_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_TOKEN", ""),
				Description: "API key for the kong api (Enterprise Edition)",
			},
			"kong_max_retries": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_MAX_RETRIES", "10"),
				Description: "Max retries if kong is having trouble",
			},
			"kong_retry_interval": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_RETRY_INTERVAL", "2"),
				Description: "Interval in seconds to wait before retrying on kong errors",
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
			"kong_service":                resourceKongService(),
			"kong_route":                  resourceKongRoute(),
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
		if v, exists := os.LookupEnv(key); exists && v != "" {
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
		HostAddress:        d.Get("kong_admin_uri").(string),
		Username:           d.Get("kong_admin_username").(string),
		Password:           d.Get("kong_admin_password").(string),
		InsecureSkipVerify: d.Get("tls_skip_verify").(bool),
		ApiKey:             d.Get("kong_api_key").(string),
		AdminToken:         d.Get("kong_admin_token").(string),
		MaxRetries:         d.Get("kong_max_retries").(int),
		RetryInterval:      d.Get("kong_retry_interval").(int),
	}

	return gokong.NewClient(config), nil
}
