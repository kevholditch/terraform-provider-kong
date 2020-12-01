# kong_sni

For more information on creating SNIs in Kong [see their documentaton](https://getkong.org/docs/1.0.x/admin-api/#sni-objects)

## Example Usage

```hcl
resource "kong_certificate" "certificate" {
    certificate  = "public key --- 123 ----"
    private_key  = "private key --- 456 ----"
}

resource "kong_sni" "sni" {
    name  	   = "www.example.com"
    certificate_id = "${kong_certificate.certificate.id}"
}
```

## Argument Reference

* `name` - (Required) is your domain you want to assign to the certificate
* `certificate_id` - (Required) is the id of a certificate

## Import

To import a SNI:

```shell
terraform import kong_sni.<sni_identifier> <sni_id>
```
