# kong_route

The route resource maps directly onto the json for the route endpoint in Kong. For more information on the parameters [see the Kong Route create documentation](https://docs.konghq.com/gateway-oss/2.5.x/admin-api/#route-object).

To create a tcp/tls route you set `sources` and `destinations` by repeating the corresponding element (`source` or `destination`) for each source or destination you want.

## Example Usage

```hcl
resource "kong_route" "route" {
    name            = "MyRoute"
	protocols 	    = [ "http", "https" ]
	methods 	    = [ "GET", "POST" ]
	hosts 		    = [ "example2.com" ]
	paths 		    = [ "/test" ]
	strip_path 	    = false
	preserve_host 	= true
	regex_priority 	= 1
	service_id 	     = kong_service.service.id
    header {
      name   = "x-test-1"
      values = ["a", "b"]
    }
}

```

To create a tcp/tls route you set `sources` and `destinations` by repeating the corresponding element (`source` or `destination`) for each source or destination you want, for example:

```hcl

resource "kong_route" "route" {
	protocols 		= [ "tcp" ]
	strip_path 		= true
	preserve_host 	= false
	source {
		ip   = "192.168.1.1"
		port = 80
	}
	source {
		ip   = "192.168.1.2"
	}
	destination {
		ip   = "172.10.1.1"
		port = 81
	}
	snis		= ["foo.com"]
	service_id  = kong_service.service.id
}
```

## Argument Reference

* `name` - (Optional) The name of the route
* `protocols` - (Required) The list of protocols to use
* `methods` - (Optional) A list of HTTP methods that match this Route
* `hosts` - (Optional) A list of domain names that match this Route  
* `paths` - (Optional) A list of paths that match this Route
* `header` - (Optional) One or more blocks of `name` to set name of header and `values` which is a list of `string` for the header values to match on.  See above example of how to set.  These headers will cause this Route to match if present in the request. The Host header cannot be used with this attribute: hosts should be specified using the hosts attribute.
* `https_redirect_status_code` - (Optional) The status code Kong responds with when all properties of a Route match except the protocol i.e. if the protocol of the request is HTTP instead of HTTPS. Location header is injected by Kong if the field is set to `301`, `302`, `307` or `308`. Accepted values are: `426`, `301`, `302`, `307`, `308`. Default: `426`.  
* `strip_path` - (Optional) When matching a Route via one of the paths, strip the matching prefix from the upstream request URL. Default: true.
* `regex_priority` - (Optional) A number used to choose which route resolves a given request when several routes match it using regexes simultaneously.
* `path_handling` - (Optional) Controls how the Service path, Route path and requested path are combined when sending a request to the upstream.
* `preserve_host` - (Optional) When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers. If set to false, the upstream Host header will be that of the Service’s host.
* `request_buffering` - (Optional) Whether to enable request body buffering or not. With HTTP 1.1, it may make sense to turn this off on services that receive data with chunked transfer encoding. Default: true.
* `response_buffering` - (Optional) Whether to enable response body buffering or not. With HTTP 1.1, it may make sense to turn this off on services that send data with chunked transfer encoding. Default: true.  
* `source` - (Required) A list of source `ip` and `port`
* `destination` - (Required) A list of destination `ip` and `port`
* `snis` - (Optional) A list of SNIs that match this Route when using stream routing.
* `service_id` - (Required) Service ID to map to
* `tags` - (Optional) A list of strings associated with the Route for grouping and filtering.


## Import

To import a route:

```shell
terraform import kong_route.<route_identifier> <route_id>
```
