package kong

import (
	"os"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
	"log"
	"github.com/kevholditch/gokong/containers"
	"github.com/kevholditch/gokong"
)

const defaultKongVersion = "0.11"

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"kong":      testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestMain(m *testing.M) {

	testContext := containers.StartKong(GetEnvVarOrDefault("KONG_VERSION", defaultKongVersion))

	err := os.Setenv(gokong.EnvKongAdminHostAddress, testContext.KongHostAddress)
	if err != nil {
		log.Fatalf("Could not set kong host address env variable: %v", err)
	}

	code := m.Run()

	containers.StopKong(testContext)

	os.Exit(code)

}

