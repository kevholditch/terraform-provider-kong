# kong_route

The route resource maps directly onto the json for the route endpoint in Kong. For more information on the parameters [see the Kong Route create documentation](https://getkong.org/docs/1.0.x/admin-api/#route-object).

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
	service_id 	= kong_service.service.id
}

```

To create a tcp/tls route you set `sources` and `destinations` by repeating the corresponding element (`source` or `destination`) for each source or destination you want, for example:

```hcl

resource "kong_route" "route" {
	protocols 		= [ "tcp" ]
	strip_path 		= true
	preserve_host 		= false
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
	service_id  	= kong_service.service.id
}
```

## Argument Reference

* `protocols` - (Required) The list of protocols to use
* `strip_path` - (Optional) When set to true strips the path
* `preserve_host` - (Optional) When set to true preserves the host header
* `source` - (Required) Source `ip` and `port`
* `destination` - (Required) Destination `ip` and `port`
* `snis` - (Optional) List of SNIs to use
* `service_id` - (Required) Service ID to map to

## Import

To import a route:

```shell
terraform import kong_route.<route_identifier> <route_id>
```
