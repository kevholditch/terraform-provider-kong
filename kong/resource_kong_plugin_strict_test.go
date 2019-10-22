package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKongGlobalPluginStrict(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateExplicitStrictGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
			{
				Config: testUpdateExplicitStrictGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccKongGlobalPluginImplicitStrict(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateImplicitStrictGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
			{
				Config: testUpdateImplicitStrictGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccKongPluginImportConfigJsonStrict(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateExplicitStrictGlobalPluginConfig,
			},
			{
				ResourceName:      "kong_plugin.hmac_auth",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testCreateExplicitStrictGlobalPluginConfig = `
resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	enabled = "true"
    strict_match = true
	config_json = <<EOT
	{
   	"algorithms": [
   	    "hmac-sha1",
   	    "hmac-sha256",
   	    "hmac-sha384",
   	    "hmac-sha512"
   	],
    "anonymous": null,
   	"clock_skew": 300,
   	"enforce_headers": [],
   	"hide_credentials": true,
   	"validate_request_body": false
	}
EOT
}`

const testCreateImplicitStrictGlobalPluginConfig = `
provider "kong" {
    strict_plugins_match = "true"
}

resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	enabled = "true"
	config_json = <<EOT
	{
   	"algorithms": [
   	    "hmac-sha1",
   	    "hmac-sha256",
   	    "hmac-sha384",
   	    "hmac-sha512"
   	],
    "anonymous": null,
   	"clock_skew": 300,
   	"enforce_headers": [],
   	"hide_credentials": true,
   	"validate_request_body": false
	}
EOT
}`

const testUpdateExplicitStrictGlobalPluginConfig = `
resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
    strict_match = true
	config_json = <<EOT
	{
   	"algorithms": [
   	    "hmac-sha1",
   	    "hmac-sha256",
   	    "hmac-sha384",
   	    "hmac-sha512"
   	],
    "anonymous": null,
   	"clock_skew": 300,
   	"enforce_headers": [],
   	"hide_credentials": false,
   	"validate_request_body": false
	}
EOT
}`

const testUpdateImplicitStrictGlobalPluginConfig = `
provider "kong" {
    strict_plugins_match = "true"
}

resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	config_json = <<EOT
	{
   	"algorithms": [
   	    "hmac-sha1",
   	    "hmac-sha256",
   	    "hmac-sha384",
   	    "hmac-sha512"
   	],
    "anonymous": null,
   	"clock_skew": 300,
   	"enforce_headers": [],
   	"hide_credentials": false,
   	"validate_request_body": false
	}
EOT
}
`
