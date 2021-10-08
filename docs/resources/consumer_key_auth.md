# kong_consumer_key_auth

Resource that allows you to configure the [Key Authentication](https://docs.konghq.com/hub/kong-inc/key-auth/) plugin for a consumer.

## Example Usage

```hcl
resource "kong_consumer" "my_consumer" {
  username  = "User1"
  custom_id = "123"
}

resource "kong_plugin" "key_auth_plugin" {
  name = "key-auth"
}

resource "kong_consumer_key_auth" "consumer_key_auth" {
  consumer_id = kong_consumer.my_consumer.id
  key         = "secret"
  tags        = ["myTag", "anotherTag"]
}
```

## Argument Reference

* `consumer_id` - (Required) the id of the consumer to associate the credentials to
* `key` - (Optional) Unique key to authenticate the client; if omitted the plugin will generate one
* `tags` - (Optional) A list of strings associated with the consumer key auth for grouping and filtering
