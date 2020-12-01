# kong_upstream

## Example Usage

```hcl
resource "kong_upstream" "upstream" {
    name                 = "sample_upstream"
    slots                = 10
    hash_on              = "header"
    hash_fallback        = "cookie"
    hash_on_header       = "HeaderName"
    hash_fallback_header = "FallbackHeaderName"
    hash_on_cookie       = "CookieName"
    hash_on_cookie_path  = "/path"
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
```

## Argument Reference

* `name` - (Required) is a hostname, which must be equal to the host of a Service.
* `slots` - (Optional) is the number of slots in the load balancer algorithm (10*65536, defaults to 10000).
* `hash_on` - (Optional) is a hashing input type: `none `(resulting in a weighted*round*robin scheme with no hashing), `consumer`, `ip`, `header`, or `cookie`. Defaults to `none`.
* `hash_fallback` - (Optional) is a hashing input type if the primary `hash_on` does not return a hash (eg. header is missing, or no consumer identified). One of: `none`, `consumer`, `ip`, `header`, or `cookie`. Not available if `hash_on` is set to `cookie`. Defaults to `none`.
* `hash_on_header` - (Optional) is a header name to take the value from as hash input. Only required when `hash_on` is set to `header`. Default `nil`.
* `hash_fallback_header` - (Optional) is a header name to take the value from as hash input. Only required when `hash_fallback` is set to `header`. Default `nil`.
* `hash_on_cookie` - (Optional) is a cookie name to take the value from as hash input. Only required when `hash_on` or `hash_fallback` is set to `cookie`. If the specified cookie is not in the request, Kong will generate a value and set the cookie in the response. Default `nil`.
* `hash_on_cookie_path` - (Optional) is a cookie path to set in the response headers. Only required when `hash_on` or `hash_fallback` is set to `cookie`. Defaults to `/`.
* `healthchecks.active.type` - (Optional) is a active health check type. HTTP or HTTPS, or just attempt a TCP connection. Possible values are `tcp`, `http` or `https`. Defaults to `http`.
* `healthchecks.active.timeout` - (Optional) is a socket timeout for active health checks (in seconds). Defaults to `1`.
* `healthchecks.active.concurrency` - (Optional) is a number of targets to check concurrently in active health checks. Defaults to `10`.
* `healthchecks.active.http_path` - (Optional) is a path to use in GET HTTP request to run as a probe on active health checks. Defaults to `/`.
* `healthchecks.active.https_verify_certificate` - (Optional) check the validity of the SSL certificate of the remote host when performing active health checks using HTTPS. Defaults to `true`.
* `healthchecks.active.https_sni` - (Optional) is the hostname to use as an SNI (Server Name Identification) when performing active health checks using HTTPS. This is particularly useful when Targets are configured using IPs, so that the target host’s certificate can be verified with the proper SNI. Default `nil`.
* `healthchecks.active.healthy.interval` - (Optional) is an interval between active health checks for healthy targets (in seconds). A value of zero indicates that active probes for healthy targets should not be performed. Defaults to `0`.
* `healthchecks.active.healthy.successes` - (Optional) is a number of successes in active probes (as defined by `healthchecks.active.healthy.http_statuses`) to consider a target healthy. Defaults to `0`.
* `healthchecks.active.healthy.http_statuses` - (Optional) is an array of HTTP statuses to consider a success, indicating healthiness, when returned by a probe in active health checks. Defaults to `[200, 302]`.
* `healthchecks.active.unhealthy.interval` - (Optional) is an interval between active health checks for unhealthy targets (in seconds). A value of zero indicates that active probes for unhealthy targets should not be performed. Defaults to `0`.
* `healthchecks.active.unhealthy.tcp_failures` - (Optional) is a number of TCP failures in active probes to consider a target unhealthy. Defaults to `0`.
* `healthchecks.active.unhealthy.http_failures` - (Optional) is a number of HTTP failures in active probes (as defined by `healthchecks.active.unhealthy.http_statuses`) to consider a target unhealthy. Defaults to `0`.
* `healthchecks.active.unhealthy.timeouts` - (Optional) is a number of timeouts in active probes to consider a target unhealthy. Defaults to `0`.
* `healthchecks.active.unhealthy.http_statuses` - (Optional) is an array of HTTP statuses to consider a failure, indicating unhealthiness, when returned by a probe in active health checks. Defaults to `[429, 404, 500, 501, 502, 503, 504, 505]`.
* `healthchecks.passive.type` - (Optional) is a passive health check type. Interpreting HTTP/HTTPS statuses, or just check for TCP connection success. Possible values are `tcp`, `http` or `https` (in passive checks, `http` and `https` options are equivalent.). Defaults to `http`.
* `healthchecks.passive.healthy.successes` - (Optional) is a Number of successes in proxied traffic (as defined by `healthchecks.passive.healthy.http_statuses`) to consider a target healthy, as observed by passive health checks. Defaults to `0`.
* `healthchecks.passive.healthy.http_statuses` - (Optional) is an array of HTTP statuses which represent healthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308]`.
* `healthchecks.passive.unhealthy.tcp_failures` - (Optional) is a number of TCP failures in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
* `healthchecks.passive.unhealthy.http_failures` - (Optional) is a number of HTTP failures in proxied traffic (as defined by `healthchecks.passive.unhealthy.http_statuses`) to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
* `healthchecks.passive.unhealthy.timeouts` - (Optional) is a number of timeouts in proxied traffic to consider a target unhealthy, as observed by passive health checks. Defaults to `0`.
* `healthchecks.passive.unhealthy.http_statuses` - (Optional) is an array of HTTP statuses which represent unhealthiness when produced by proxied traffic, as observed by passive health checks. Defaults to `[429, 500, 503]`.

## Import

To import an upstream:

```shell
terraform import kong_upstream.<upstream_identifier> <upstream_id>
```
