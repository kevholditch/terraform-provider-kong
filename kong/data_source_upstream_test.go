package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKongUpstream(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testUpstreamDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "name", "TestUpstream"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "slots", "10"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "hash_on", "header"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "hash_fallback", "consumer"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "hash_on_header", "HeaderName"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "hash_fallback_header", "FallbackHeaderName"),

					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.timeout", "10"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.concurrency", "20"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.http_path", "/status"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.healthy.0.interval", "5"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.healthy.0.successes", "1"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.healthy.0.http_statuses.1", "201"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.interval", "3"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.tcp_failures", "1"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.http_failures", "2"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.timeouts", "7"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.http_statuses.0", "500"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.active.0.unhealthy.0.http_statuses.1", "501"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.healthy.0.successes", "1"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.healthy.0.http_statuses.0", "200"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.healthy.0.http_statuses.1", "201"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.healthy.0.http_statuses.2", "202"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.tcp_failures", "5"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.http_failures", "6"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.timeouts", "3"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.http_statuses.0", "500"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.http_statuses.1", "501"),
					resource.TestCheckResourceAttr("data.kong_upstream.upstream_data_source", "healthchecks.0.passive.0.unhealthy.0.http_statuses.2", "502"),
				),
			},
		},
	})
}

const testUpstreamDataSourceConfig = `
resource "kong_upstream" "upstream" {
	name  		= "TestUpstream"
	slots 		= 10
	hash_on              = "header"
	hash_fallback        = "consumer"
	hash_on_header       = "HeaderName"
	hash_fallback_header = "FallbackHeaderName"
	healthchecks         = {
		active = {
			http_path    = "/status"
			timeout      = 10
			concurrency  = 20
			healthy = {
				successes = 1
				interval  = 5
				http_statuses = [200, 201]
			}
			unhealthy = {
				timeouts      = 7
				interval      = 3
				tcp_failures  = 1
				http_failures = 2
				http_statuses = [500, 501]
			}
		}
		passive = {
			healthy = {
				successes = 1
				http_statuses = [200, 201, 202]
			}
			unhealthy = {
				timeouts      = 3
				tcp_failures  = 5
				http_failures = 6
				http_statuses = [500, 501, 502]
			}
		}
	}
}

data "kong_upstream" "upstream_data_source" {
	filter = {
		id   = "${kong_upstream.upstream.id}"
		name = "TestUpstream"
	}
}
`
