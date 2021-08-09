package kong

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccKongUpstream(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongUpstreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateUpstreamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "name", "MyUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "slots", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on", "none"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_fallback", "none"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_header", ""),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_fallback_header", ""),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_cookie", ""),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_cookie_path", "/"),

					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.type", "http"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.timeout", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.concurrency", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.http_path", "/"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.https_verify_certificate", "true"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.https_sni", ""),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.interval", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.successes", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.http_statuses.1", "302"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.interval", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.tcp_failures", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_failures", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.timeouts", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.0", "429"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.1", "404"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.2", "500"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.3", "501"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.4", "502"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.5", "503"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.6", "504"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.7", "505"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.type", "http"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.successes", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.1", "201"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.2", "202"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.3", "203"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.4", "204"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.5", "205"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.6", "206"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.7", "207"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.8", "208"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.9", "226"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.10", "300"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.11", "301"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.12", "302"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.13", "303"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.14", "304"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.15", "305"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.16", "306"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.17", "307"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.18", "308"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.tcp_failures", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_failures", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.timeouts", "0"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.0", "429"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.1", "500"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.2", "503"),
				),
			},
			{
				Config: fmt.Sprintf(testUpdateUpstreamConfig, testCert1, testKey1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongUpstreamExists("kong_upstream.upstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "name", "MyUpstream"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "slots", "20"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on", "header"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_fallback", "cookie"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_header", "HeaderName"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_fallback_header", "FallbackHeaderName"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_cookie", "CookieName"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "hash_on_cookie_path", "/path"),

					resource.TestCheckResourceAttr("kong_upstream.upstream", "host_header", "x-host"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "tags.0", "a"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "tags.1", "b"),
					func(s *terraform.State) error {
						module := s.RootModule()
						cert, ok := module.Resources["kong_certificate.certificate"]
						if !ok {
							return fmt.Errorf("could not find certificate resource")
						}

						service, ok := module.Resources["kong_upstream.upstream"]
						if !ok {
							return fmt.Errorf("could not find upstream resource")
						}

						v, ok := service.Primary.Attributes["client_certificate_id"]
						if !ok {
							return fmt.Errorf("could not find client_certificate_id property")
						}

						if v != cert.Primary.ID {
							return fmt.Errorf("client_certificate_id does not match certificate id")
						}
						return nil
					},

					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.type", "https"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.timeout", "10"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.concurrency", "20"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.http_path", "/status"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.https_verify_certificate", "false"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.https_sni", "some.domain.com"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.interval", "5"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.successes", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.healthy.0.http_statuses.1", "201"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.interval", "3"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.tcp_failures", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_failures", "2"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.timeouts", "7"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.0", "500"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.active.0.unhealthy.0.http_statuses.1", "501"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.type", "https"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.successes", "1"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.1", "201"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.healthy.0.http_statuses.2", "202"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.tcp_failures", "5"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_failures", "6"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.timeouts", "3"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.0", "500"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.1", "501"),
					resource.TestCheckResourceAttr("kong_upstream.upstream", "healthchecks.0.passive.0.unhealthy.0.http_statuses.2", "502"),
				),
			},
		},
	})
}

func TestAccKongUpstreamImport(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongUpstreamDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testCreateUpstreamConfig,
			},

			resource.TestStep{
				ResourceName:      "kong_upstream.upstream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKongUpstreamDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.Upstreams

	upstreams := getResourcesByType("kong_upstream", state)

	if len(upstreams) != 1 {
		return fmt.Errorf("expecting only 1 upstream resource found %v", len(upstreams))
	}

	response, err := client.Get(context.Background(), kong.String(upstreams[0].Primary.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get upstream by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("upstream %s still exists, %+v", upstreams[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckKongUpstreamExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.Upstreams
		api, err := client.Get(context.Background(), kong.String(rs.Primary.ID))

		if err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("upstream with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

func TestCreateKongHealthCheckFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.Healthcheck
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.Healthcheck{},
		}, // Simple data
		{
			in: &map[string]interface{}{},
			expected: &kong.Healthcheck{
				Active:  nil,
				Passive: nil,
			},
		}, // All data
		{
			in: &map[string]interface{}{
				"active": []interface{}{
					map[string]interface{}{
						"type":                     "http",
						"concurrency":              12,
						"http_path":                "/health",
						"https_verify_certificate": true,
						"timeout":                  60,
					},
				},
				"passive": []interface{}{
					map[string]interface{}{
						"type": "https",
					},
				},
			},
			expected: &kong.Healthcheck{
				Active: &kong.ActiveHealthcheck{
					Type:                   kong.String("http"),
					Concurrency:            kong.Int(12),
					Healthy:                nil,
					HTTPPath:               kong.String("/health"),
					HTTPSVerifyCertificate: kong.Bool(true),
					HTTPSSni:               nil,
					Timeout:                kong.Int(60),
					Unhealthy:              nil,
				},
				Passive: &kong.PassiveHealthcheck{
					Type:      kong.String("https"),
					Healthy:   nil,
					Unhealthy: nil,
				},
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongHealthCheckFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongHealthCheckActiveFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.ActiveHealthcheck
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.ActiveHealthcheck{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"type":                     "http",
				"concurrency":              12,
				"http_path":                "/health",
				"https_verify_certificate": true,
				"timeout":                  60,
			},
			expected: &kong.ActiveHealthcheck{
				Type:                   kong.String("http"),
				Concurrency:            kong.Int(12),
				Healthy:                nil,
				HTTPPath:               kong.String("/health"),
				HTTPSVerifyCertificate: kong.Bool(true),
				HTTPSSni:               nil,
				Timeout:                kong.Int(60),
				Unhealthy:              nil,
			},
		}, // All data
		{
			in: &map[string]interface{}{
				"type":                     "http",
				"concurrency":              12,
				"http_path":                "/health",
				"https_verify_certificate": true,
				"timeout":                  60,
				"https_sni":                "some.domain.com",
				"healthy": []interface{}{
					map[string]interface{}{
						"successes":     3,
						"interval":      5,
						"http_statuses": []interface{}{200},
					},
				},
				"unhealthy": []interface{}{
					map[string]interface{}{
						"http_failures": 1,
						"http_statuses": []interface{}{500},
						"interval":      5,
						"tcp_failures":  2,
						"timeouts":      4,
					},
				},
			},
			expected: &kong.ActiveHealthcheck{
				Type:        kong.String("http"),
				Concurrency: kong.Int(12),
				Healthy: &kong.Healthy{
					Successes:    kong.Int(3),
					Interval:     kong.Int(5),
					HTTPStatuses: []int{200},
				},
				HTTPPath:               kong.String("/health"),
				HTTPSVerifyCertificate: kong.Bool(true),
				HTTPSSni:               kong.String("some.domain.com"),
				Timeout:                kong.Int(60),
				Unhealthy: &kong.Unhealthy{
					HTTPFailures: kong.Int(1),
					HTTPStatuses: []int{500},
					Interval:     kong.Int(5),
					TCPFailures:  kong.Int(2),
					Timeouts:     kong.Int(4),
				},
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongHealthCheckActiveFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongHealthCheckPassiveFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.PassiveHealthcheck
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.PassiveHealthcheck{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"type": "http",
			},
			expected: &kong.PassiveHealthcheck{
				Type:      kong.String("http"),
				Healthy:   nil,
				Unhealthy: nil,
			},
		}, // All data
		{
			in: &map[string]interface{}{
				"type": "https",
				"healthy": []interface{}{
					map[string]interface{}{
						"successes":     3,
						"http_statuses": []interface{}{200},
					},
				},
				"unhealthy": []interface{}{
					map[string]interface{}{
						"http_failures": 1,
						"http_statuses": []interface{}{500},
						"tcp_failures":  2,
						"timeouts":      4,
					},
				},
			},
			expected: &kong.PassiveHealthcheck{
				Type: kong.String("https"),
				Healthy: &kong.Healthy{
					Successes:    kong.Int(3),
					HTTPStatuses: []int{200},
				},
				Unhealthy: &kong.Unhealthy{
					HTTPFailures: kong.Int(1),
					HTTPStatuses: []int{500},
					TCPFailures:  kong.Int(2),
					Timeouts:     kong.Int(4),
				},
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongHealthCheckPassiveFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongActiveHealthyFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.Healthy
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.Healthy{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"interval":      3,
				"http_statuses": []interface{}{200},
				"successes":     2,
			},
			expected: &kong.Healthy{
				HTTPStatuses: []int{200},
				Interval:     kong.Int(3),
				Successes:    kong.Int(2),
			},
		}, // EmptyHTTPStatuses
		{
			in: &map[string]interface{}{
				"interval":      3,
				"http_statuses": []interface{}{},
				"successes":     2,
			},
			expected: &kong.Healthy{
				HTTPStatuses: []int{},
				Interval:     kong.Int(3),
				Successes:    kong.Int(2),
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongActiveHealthyFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongActiveUnhealthyFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.Unhealthy
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.Unhealthy{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"http_failures": 4,
				"http_statuses": []interface{}{200},
				"interval":      3,
				"tcp_failures":  5,
				"timeouts":      6,
			},
			expected: &kong.Unhealthy{
				HTTPFailures: kong.Int(4),
				HTTPStatuses: []int{200},
				Interval:     kong.Int(3),
				TCPFailures:  kong.Int(5),
				Timeouts:     kong.Int(6),
			},
		}, // EmptyHTTPStatuses
		{
			in: &map[string]interface{}{
				"http_failures": 4,
				"http_statuses": []interface{}{},
				"interval":      3,
				"tcp_failures":  5,
				"timeouts":      6,
			},
			expected: &kong.Unhealthy{
				HTTPFailures: kong.Int(4),
				HTTPStatuses: []int{},
				Interval:     kong.Int(3),
				TCPFailures:  kong.Int(5),
				Timeouts:     kong.Int(6),
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongActiveUnhealthyFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongPassiveHealthyFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.Healthy
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.Healthy{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"http_statuses": []interface{}{
					200,
				},
				"successes": 3,
			},
			expected: &kong.Healthy{
				HTTPStatuses: []int{200},
				Successes:    kong.Int(3),
			},
		}, // EmptyHTTPStatuses
		{
			in: &map[string]interface{}{
				"http_statuses": []interface{}{},
				"successes":     3,
			},
			expected: &kong.Healthy{
				HTTPStatuses: []int{},
				Successes:    kong.Int(3),
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongPassiveHealthyFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestCreateKongPassiveUnhealthyFromMap(t *testing.T) {
	cases := []struct {
		in       *map[string]interface{}
		expected *kong.Unhealthy
	}{
		// Empty data
		{
			in:       &map[string]interface{}{},
			expected: &kong.Unhealthy{},
		}, // Simple data
		{
			in: &map[string]interface{}{
				"http_statuses": []interface{}{200},
				"tcp_failures":  3,
				"http_failures": 4,
				"timeouts":      5,
			},
			expected: &kong.Unhealthy{
				HTTPStatuses: []int{200},
				TCPFailures:  kong.Int(3),
				HTTPFailures: kong.Int(4),
				Timeouts:     kong.Int(5),
			},
		}, // EmptyHTTPStatuses
		{
			in: &map[string]interface{}{
				"http_statuses": []interface{}{},
				"tcp_failures":  3,
				"http_failures": 4,
				"timeouts":      5,
			},
			expected: &kong.Unhealthy{
				HTTPStatuses: []int{},
				TCPFailures:  kong.Int(3),
				HTTPFailures: kong.Int(4),
				Timeouts:     kong.Int(5),
			},
		},
		{
			in:       nil,
			expected: nil,
		},
	}

	for _, c := range cases {
		out := createKongPassiveUnhealthyFromMap(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenHealthCheck(t *testing.T) {
	cases := []struct {
		in       *kong.Healthcheck
		expected []interface{}
	}{
		// Simple data
		{
			in: &kong.Healthcheck{
				Active:  nil,
				Passive: nil,
			},
			expected: []interface{}{
				map[string]interface{}{},
			},
		}, // All data
		{
			in: &kong.Healthcheck{
				Active: &kong.ActiveHealthcheck{
					Type:                   kong.String("http"),
					Concurrency:            kong.Int(12),
					Healthy:                nil,
					HTTPPath:               kong.String("/health"),
					HTTPSVerifyCertificate: kong.Bool(true),
					HTTPSSni:               nil,
					Timeout:                kong.Int(60),
					Unhealthy:              nil,
				},
				Passive: &kong.PassiveHealthcheck{
					Type:      kong.String("https"),
					Healthy:   nil,
					Unhealthy: nil,
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"active": []interface{}{
						map[string]interface{}{
							"type":                     "http",
							"concurrency":              12,
							"http_path":                "/health",
							"https_verify_certificate": true,
							"timeout":                  60,
						},
					},
					"passive": []interface{}{
						map[string]interface{}{
							"type": "https",
						},
					},
				},
			},
		}, // Nil object
		{
			in:       nil,
			expected: []interface{}{},
		},
	}

	for _, c := range cases {
		out := flattenHealthCheck(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenHealthCheckActive(t *testing.T) {
	cases := []struct {
		in       kong.ActiveHealthcheck
		expected []interface{}
	}{
		// Simple data
		{
			in: kong.ActiveHealthcheck{
				Type:                   kong.String("http"),
				Concurrency:            kong.Int(12),
				Healthy:                nil,
				HTTPPath:               kong.String("/health"),
				HTTPSVerifyCertificate: kong.Bool(true),
				HTTPSSni:               nil,
				Timeout:                kong.Int(60),
				Unhealthy:              nil,
			},
			expected: []interface{}{
				map[string]interface{}{
					"type":                     "http",
					"concurrency":              12,
					"http_path":                "/health",
					"https_verify_certificate": true,
					"timeout":                  60,
				},
			},
		}, // All data
		{
			in: kong.ActiveHealthcheck{
				Type:        kong.String("http"),
				Concurrency: kong.Int(12),
				Healthy: &kong.Healthy{
					Successes:    kong.Int(3),
					Interval:     kong.Int(5),
					HTTPStatuses: []int{200},
				},
				HTTPPath:               kong.String("/health"),
				HTTPSVerifyCertificate: kong.Bool(true),
				HTTPSSni:               kong.String("some.domain.com"),
				Timeout:                kong.Int(60),
				Unhealthy: &kong.Unhealthy{
					HTTPFailures: kong.Int(1),
					HTTPStatuses: []int{500},
					Interval:     kong.Int(5),
					TCPFailures:  kong.Int(2),
					Timeouts:     kong.Int(4),
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"type":                     "http",
					"concurrency":              12,
					"http_path":                "/health",
					"https_verify_certificate": true,
					"timeout":                  60,
					"https_sni":                "some.domain.com",
					"healthy": []map[string]interface{}{
						{
							"successes":     3,
							"interval":      5,
							"http_statuses": []int{200},
						},
					},
					"unhealthy": []map[string]interface{}{
						{
							"http_failures": 1,
							"http_statuses": []int{500},
							"interval":      5,
							"tcp_failures":  2,
							"timeouts":      4,
						},
					},
				},
			},
		}, // Nil object
		{
			in:       kong.ActiveHealthcheck{},
			expected: []interface{}{map[string]interface{}{}},
		},
	}

	for _, c := range cases {
		out := flattenHealthCheckActive(&c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenHealthCheckPassive(t *testing.T) {
	cases := []struct {
		in       kong.PassiveHealthcheck
		expected []interface{}
	}{
		// Simple data
		{
			in: kong.PassiveHealthcheck{
				Type:      kong.String("http"),
				Healthy:   nil,
				Unhealthy: nil,
			},
			expected: []interface{}{
				map[string]interface{}{
					"type": "http",
				},
			},
		}, // All data
		{
			in: kong.PassiveHealthcheck{
				Type: kong.String("https"),
				Healthy: &kong.Healthy{
					Successes:    kong.Int(3),
					HTTPStatuses: []int{200},
				},
				Unhealthy: &kong.Unhealthy{
					HTTPFailures: kong.Int(1),
					HTTPStatuses: []int{500},
					TCPFailures:  kong.Int(2),
					Timeouts:     kong.Int(4),
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"type": "https",
					"healthy": []map[string]interface{}{
						{
							"successes":     3,
							"http_statuses": []int{200},
						},
					},
					"unhealthy": []map[string]interface{}{
						{
							"http_failures": 1,
							"http_statuses": []int{500},
							"tcp_failures":  2,
							"timeouts":      4,
						},
					},
				},
			},
		}, // Nil object
		{
			in:       kong.PassiveHealthcheck{},
			expected: []interface{}{map[string]interface{}{}},
		},
	}

	for _, c := range cases {
		out := flattenHealthCheckPassive(&c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenActiveHealthy(t *testing.T) {
	cases := []struct {
		in       *kong.Healthy
		expected []map[string]interface{}
	}{
		// Simple, all data
		{
			in: &kong.Healthy{
				HTTPStatuses: []int{200},
				Interval:     kong.Int(3),
				Successes:    kong.Int(2),
			},
			expected: []map[string]interface{}{
				{
					"interval":      3,
					"http_statuses": []int{200},
					"successes":     2,
				},
			},
		}, // EmptyHTTPStatuses
		{
			in: &kong.Healthy{
				HTTPStatuses: []int{},
				Interval:     kong.Int(3),
				Successes:    kong.Int(2),
			},
			expected: []map[string]interface{}{
				{
					"interval":      3,
					"http_statuses": []int{},
					"successes":     2,
				},
			},
		}, // Nil object
		{
			in:       nil,
			expected: []map[string]interface{}{},
		},
	}

	for _, c := range cases {
		out := flattenActiveHealthy(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenActiveUnhealthy(t *testing.T) {
	cases := []struct {
		in       *kong.Unhealthy
		expected []map[string]interface{}
	}{
		// Simple, all data
		{
			in: &kong.Unhealthy{
				HTTPFailures: kong.Int(4),
				HTTPStatuses: []int{200},
				Interval:     kong.Int(3),
				TCPFailures:  kong.Int(5),
				Timeouts:     kong.Int(6),
			},
			expected: []map[string]interface{}{
				{
					"http_failures": 4,
					"http_statuses": []int{200},
					"interval":      3,
					"tcp_failures":  5,
					"timeouts":      6,
				},
			},
		}, // EmptyHTTPStatuses
		{
			in: &kong.Unhealthy{
				HTTPFailures: kong.Int(4),
				HTTPStatuses: []int{},
				Interval:     kong.Int(3),
				TCPFailures:  kong.Int(5),
				Timeouts:     kong.Int(6),
			},
			expected: []map[string]interface{}{
				{
					"http_failures": 4,
					"http_statuses": []int{},
					"interval":      3,
					"tcp_failures":  5,
					"timeouts":      6,
				},
			},
		}, // Nil object
		{
			in:       nil,
			expected: []map[string]interface{}{},
		},
	}

	for _, c := range cases {
		out := flattenActiveUnhealthy(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenPassiveHealthy(t *testing.T) {
	cases := []struct {
		in       *kong.Healthy
		expected []map[string]interface{}
	}{
		// Simple, all data
		{
			in: &kong.Healthy{
				HTTPStatuses: []int{200},
				Successes:    kong.Int(3),
			},
			expected: []map[string]interface{}{
				{
					"http_statuses": []int{200},
					"successes":     3,
				},
			},
		}, // EmptyHTTPStatuses
		{
			in: &kong.Healthy{
				HTTPStatuses: []int{},
				Successes:    kong.Int(3),
			},
			expected: []map[string]interface{}{
				{
					"http_statuses": []int{},
					"successes":     3,
				},
			},
		}, // Nil object
		{
			in:       nil,
			expected: []map[string]interface{}{},
		},
	}

	for _, c := range cases {
		out := flattenPassiveHealthy(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

func TestFlattenPassiveUnhealthy(t *testing.T) {
	cases := []struct {
		in       *kong.Unhealthy
		expected []map[string]interface{}
	}{
		// Simple, all data
		{
			in: &kong.Unhealthy{
				HTTPStatuses: []int{200},
				TCPFailures:  kong.Int(3),
				HTTPFailures: kong.Int(4),
				Timeouts:     kong.Int(5),
			},
			expected: []map[string]interface{}{
				{
					"http_statuses": []int{200},
					"tcp_failures":  3,
					"http_failures": 4,
					"timeouts":      5,
				},
			},
		}, // EmptyHTTPStatuses
		{
			in: &kong.Unhealthy{
				HTTPStatuses: []int{},
				TCPFailures:  kong.Int(3),
				HTTPFailures: kong.Int(4),
				Timeouts:     kong.Int(5),
			},
			expected: []map[string]interface{}{
				{
					"http_statuses": []int{},
					"tcp_failures":  3,
					"http_failures": 4,
					"timeouts":      5,
				},
			},
		}, // Nil object
		{
			in:       nil,
			expected: []map[string]interface{}{},
		},
	}

	for _, c := range cases {
		out := flattenPassiveUnhealthy(c.in)
		if !reflect.DeepEqual(out, c.expected) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", out, c.expected)
		}
	}
}

const testCreateUpstreamConfig = `
resource "kong_upstream" "upstream" {
	name  		= "MyUpstream"
	slots 		= 10
}
`
const testUpdateUpstreamConfig = `
resource "kong_certificate" "certificate" {
	certificate  = <<EOF
%s
EOF
	private_key =  <<EOF
%s
EOF
   snis			= ["foo.com"]
}

resource "kong_upstream" "upstream" {
	name  		         = "MyUpstream"
	slots 		         = 20
	hash_on              = "header"
	hash_fallback        = "cookie"
	hash_on_header       = "HeaderName"
	hash_fallback_header = "FallbackHeaderName"
	hash_on_cookie       = "CookieName"
	hash_on_cookie_path  = "/path"
	host_header          = "x-host"
    tags                 = ["a", "b"]
    client_certificate_id = kong_certificate.certificate.id

	healthchecks {
		active {
			type                     = "https"
			http_path                = "/status"
			timeout                  = 10
			concurrency              = 20
			https_verify_certificate = false
			https_sni                = "some.domain.com"
			healthy {
				successes = 1
				interval  = 5
				http_statuses = [200, 201]
			}
			unhealthy {
				timeouts      = 7
				interval      = 3
				tcp_failures  = 1
				http_failures = 2
				http_statuses = [500, 501]
			}
		}
		passive {
			type    = "https"
			healthy {
				successes = 1
				http_statuses = [200, 201, 202]
			}
			unhealthy {
				timeouts      = 3
				tcp_failures  = 5
				http_failures = 6
				http_statuses = [500, 501, 502]
			}
		}
	}
}
`
