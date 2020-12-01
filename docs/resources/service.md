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

`name` - (Required) Service name
`protocol` - (Required) Protocol to use
`host` - (Required) Host to map to
`port` - (Required) Port to map to
`path` - (Optional) Path to map to
`retries` - (Optional) Number of retries
`connect_timeout` - (Optional) Connection timeout
`write_timeout` - (Optional) Write timout
`read_timeout` - (Optional) Read timeout

## Import

To import a service:

```shell
terraform import kong_service.<service_identifier> <service_id>
```
