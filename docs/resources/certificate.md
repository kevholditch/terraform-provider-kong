# kong_certificate

For more information on creating certificates in Kong [see their documentation](https://docs.konghq.com/gateway-oss/2.5.x/admin-api/#certificate-object)

## Example Usage

```hcl
resource "kong_certificate" "certificate" {
    certificate  = "public key --- 123 ----"
    private_key  = "private key --- 456 ----"
    snis         = ["foo.com", "bar.com"]
    tags         = ["myTag"]
}
```

## Argument Reference

* `certificate` - (Required) should be the public key of your certificate it is mapped to the `Cert` parameter on the Kong API.
* `private_key` - (Required) should be the private key of your certificate it is mapped to the `Key` parameter on the Kong API.
* `snis` - (Optional) a list of SNIs (alternative hosts on the certificate), under the bonnet this will create an SNI object in kong
* `snis` - (Optional) A list of strings associated with the Certificate for grouping and filtering

## Import

To import a certificate:

```shell
terraform import kong_certificate.<certifcate_identifier> <certificate_id>
```
