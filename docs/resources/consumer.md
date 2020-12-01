# kong_consumer

The consumer resource maps directly onto the json for creating an Consumer in Kong.  For more information on the parameters [see the Kong Consumer create documentation](https://getkong.org/docs/1.0.x/admin-api/#consumer-object).

## Example Usage

```hcl
resource "kong_consumer" "consumer" {
    username  = "User1"
    custom_id = "123"
}
```

## Argument Reference

* `username` - (Required) The usernamae to use
* `custom_id` - (Required) A custom id for the consumer

## Import

To import a consumer:

```shell
terraform import kong_consumer.<consumer_identifier> <consumer_id>
```
