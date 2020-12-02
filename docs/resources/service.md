# kong_service

The service resource maps directly onto the json for the service endpoint in Kong.  For more information on the parameters [see the Kong Service create documentation](https://getkong.org/docs/1.0.x/admin-api/#service-object).

## Example Usage

```hcl
resource "kong_service" "service" {
	name     	= "test"
	protocol 	= "http"
	host     	= "test.org"
	port     	= 8080
	path     	= "/mypath"
	retries  	= 5
	connect_timeout = 1000
	write_timeout 	= 2000
	read_timeout  	= 3000
}
```

## Argument reference

* `name` - (Required) Service name
* `protocol` - (Required) Protocol to use
* `host` - (Optional) Host to map to
* `port` - (Optional, int) Port to map to. Default: 80
* `path` - (Optional) Path to map to
* `retries` - (Optional, int) Number of retries. Default: 5
* `connect_timeout` - (Optional, int) Connection timeout. Default(ms): 60000
* `write_timeout` - (Optional, int) Write timout. Default(ms): 60000
* `read_timeout` - (Optional, int) Read timeout. Default(ms): 60000

## Import

To import a service:

```shell
terraform import kong_service.<service_identifier> <service_id>
```
