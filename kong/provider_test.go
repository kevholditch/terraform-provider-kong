package kong

import (
	"log"
	"os"
	"testing"

	terraformconfig "github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kevholditch/gokong"
	"github.com/kevholditch/gokong/containers"
)

const defaultKongVersion = "1.0.2"

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"kong": testAccProvider,
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

func TestProvider_configure(t *testing.T) {
	c, err := terraformconfig.NewRawConfig(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	rc := terraform.NewResourceConfig(c)
	p := Provider()
	err = p.Configure(rc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProvider_configure_strict(t *testing.T) {
	c, err := terraformconfig.NewRawConfig(map[string]interface{}{
		"strict_plugins_match": "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	rc := terraform.NewResourceConfig(c)
	p := Provider()
	err = p.Configure(rc)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {

	testContext := containers.StartKong(GetEnvVarOrDefault("KONG_VERSION", defaultKongVersion))

	err := os.Setenv(gokong.EnvKongAdminHostAddress, testContext.KongHostAddress)
	if err != nil {
		log.Fatalf("Could not set kong host address env variable: %v", err)
	}
	err = os.Setenv(gokong.EnvKongAdminPassword, "AnUsername")
	if err != nil {
		log.Fatalf("Could not set kong admin username env variable: %v", err)
	}
	err = os.Setenv(gokong.EnvKongAdminPassword, "AnyPassword")
	if err != nil {
		log.Fatalf("Could not set kong admin password env variable: %v", err)
	}

	code := m.Run()

	containers.StopKong(testContext)

	os.Exit(code)

}
