# kong_certificate

For more information on creating certificates in Kong [see their documentation](https://getkong.org/docs/1.0.x/admin-api/#certificate-object)

## Example Usage

```hcl
resource "kong_certificate" "certificate" {
    certificate  = "public key --- 123 ----"
    private_key = "private key --- 456 ----"
}
```

## Argument Reference

* `certificate` - (Required) should be the public key of your certificate it is mapped to the `Cert` parameter on the Kong API.
* `private_key` - (Required) should be the private key of your certificate it is mapped to the `Key` parameter on the Kong API.

## Import

To import a certificate:

```shell
terraform import kong_certificate.<certifcate_identifier> <certificate_id>
```
