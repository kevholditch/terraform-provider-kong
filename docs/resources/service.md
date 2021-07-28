# kong_service

The service resource maps directly onto the json for the service endpoint in Kong.  For more information on the parameters [see the Kong Service create documentation](https://docs.konghq.com/gateway-oss/2.5.x/admin-api/#service-object).

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
* `tags` - (Optional) A list of strings associated with the Service for grouping and filtering.
* `client_certificate` - (Optional) Certificate to be used as client certificate while TLS handshaking to the upstream server.
* `tls_verify` - (Optional) Whether to enable verification of upstream server TLS certificate.
* `tls_verify_depth` - (Optional) Maximum depth of chain while verifying Upstream server’s TLS certificate.
* `ca_certificates` - (Optional) A of CA Certificate IDs (created from the certificate resource). that are used to build the trust store while verifying upstream server’s TLS certificate.


## Import

To import a service:

```shell
terraform import kong_service.<service_identifier> <service_id>
```
