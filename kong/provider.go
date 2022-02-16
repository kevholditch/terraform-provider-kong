package kong

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

type config struct {
	adminClient           *kong.Client
	strictPlugins         bool
	strictConsumerPlugins bool
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kong_admin_uri": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_ADDR", "http://localhost:8001"),
				Description: "The address of the kong admin url e.g. http://localhost:8001",
			},
			"kong_admin_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_USERNAME", ""),
				Description: "An basic auth user for kong admin",
			},
			"kong_admin_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_PASSWORD", ""),
				Description: "An basic auth password for kong admin",
			},
			"tls_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("TLS_SKIP_VERIFY", "false"),
				Description: "Whether to skip tls verify for https kong api endpoint using self signed or untrusted certs",
			},
			"kong_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_API_KEY", ""),
				Description: "API key for the kong api (if you have locked it down)",
			},
			"kong_admin_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("KONG_ADMIN_TOKEN", ""),
				Description: "API key for the kong api (Enterprise Edition)",
			},
			"kong_workspace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workspace context (Enterprise Edition)",
			},
			"strict_plugins_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				DefaultFunc: envDefaultFuncWithDefault("STRICT_PLUGINS_MATCH", "false"),
				Description: "Should plugins `config_json` field strictly match plugin configuration",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"kong_certificate":         resourceKongCertificate(),
			"kong_consumer":            resourceKongConsumer(),
			"kong_consumer_acl":        resourceKongConsumerACL(),
			"kong_consumer_basic_auth": resourceKongConsumerBasicAuth(),
			"kong_consumer_key_auth":   resourceKongConsumerKeyAuth(),
			"kong_consumer_oauth2":     resourceKongConsumerOAuth2(),
			"kong_plugin":              resourceKongPlugin(),
			"kong_upstream":            resourceKongUpstream(),
			"kong_target":              resourceKongTarget(),
			"kong_service":             resourceKongService(),
			"kong_route":               resourceKongRoute(),
			"kong_consumer_jwt_auth":   resourceKongConsumerJWTAuth(),
		},

		//DataSourcesMap: map[string]*schema.Resource{
		//	"kong_api":         dataSourceKongApi(),
		//	"kong_certificate": dataSourceKongCertificate(),
		//	"kong_consumer":    dataSourceKongConsumer(),
		//	"kong_plugin":      dataSourceKongPlugin(),
		//	"kong_upstream":    dataSourceKongUpstream(),
		//},
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

	kongConfig := &Config{
		Address:            d.Get("kong_admin_uri").(string),
		Username:           d.Get("kong_admin_username").(string),
		Password:           d.Get("kong_admin_password").(string),
		InsecureSkipVerify: d.Get("tls_skip_verify").(bool),
		APIKey:             d.Get("kong_api_key").(string),
		AdminToken:         d.Get("kong_admin_token").(string),
		Workspace:          d.Get("kong_workspace").(string),
	}

	client, err := GetKongClient(*kongConfig)
	if err != nil {
		return nil, err
	}

	config := &config{
		adminClient:   client,
		strictPlugins: d.Get("strict_plugins_match").(bool),
	}

	return config, nil
}
