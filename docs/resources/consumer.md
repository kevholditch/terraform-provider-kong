# kong_consumer

The consumer resource maps directly onto the json for creating a Consumer in Kong.  For more information on the parameters [see the Kong Consumer create documentation](https://docs.konghq.com/gateway-oss/2.5.x/admin-api/#consumer-object).

## Example Usage

```hcl
resource "kong_consumer" "consumer" {
    username  = "User1"
    custom_id = "123"
    tags      = ["mySuperTag"]
}
```

## Argument Reference

* `username` - (Semi-optional) The username to use, you must set either the username or custom_id
* `custom_id` - (Semi-optional) A custom id for the consumer, you must set either the username or custom_id
* `tags` - (Optional) A list of strings associated with the Consumer for grouping and filtering

## Import

To import a consumer:

```shell
terraform import kong_consumer.<consumer_identifier> <consumer_id>
```
