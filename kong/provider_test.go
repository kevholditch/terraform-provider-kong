package kong

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kevholditch/terraform-provider-kong/kong/containers"
)

const defaultKongVersion = "2.5.0-ubuntu"
const EnvKongAdminHostAddress = "KONG_ADMIN_ADDR"
const EnvKongAdminUsername = "KONG_ADMIN_USERNAME"
const EnvKongAdminPassword = "KONG_ADMIN_PASSWORD"
const defaultKongRepository = "kong"
const defaultKongLicense = ""
const providerNameKong = "kong"

var (
	testAccProviders         map[string]*schema.Provider
	testAccProvider          *schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		providerNameKong: testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		providerNameKong: func() (*schema.Provider, error) { return Provider(), nil }, //nolint:unparam
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func TestProvider_configure(t *testing.T) {

	rc := terraform.NewResourceConfigRaw(map[string]interface{}{})
	p := Provider()
	err := p.Configure(context.Background(), rc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProvider_configure_strict(t *testing.T) {

	rc := terraform.NewResourceConfigRaw(map[string]interface{}{
		"strict_plugins_match": "true",
	})
	p := Provider()
	err := p.Configure(context.Background(), rc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {

	testContext := containers.StartKong(defaultKongRepository, GetEnvVarOrDefault("KONG_VERSION", defaultKongVersion), defaultKongLicense)

	err := os.Setenv(EnvKongAdminHostAddress, testContext.KongHostAddress)
	if err != nil {
		log.Fatalf("Could not set kong host address env variable: %v", err)
	}
	err = os.Setenv(EnvKongAdminPassword, "AnUsername")
	if err != nil {
		log.Fatalf("Could not set kong admin username env variable: %v", err)
	}
	err = os.Setenv(EnvKongAdminPassword, "AnyPassword")
	if err != nil {
		log.Fatalf("Could not set kong admin password env variable: %v", err)
	}

	code := m.Run()

	containers.StopKong(testContext)

	os.Exit(code)

}
