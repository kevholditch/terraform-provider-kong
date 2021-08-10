# kong_consumer_basic_auth

Consumer basic auth is a resource that allows you to configure the basic auth plugin for a consumer.

## Example Usage

```hcl
resource "kong_consumer" "my_consumer" {
  username  = "User1"
  custom_id = "123"
}

resource "kong_plugin" "basic_auth_plugin" {
  name        = "basic-auth"
}

resource "kong_consumer_basic_auth" "consumer_basic_auth" {
  consumer_id    = "${kong_consumer.my_consumer.id}"
  username       = "foo_updated"
  password       = "bar_updated"
  tags           = ["myTag", "anotherTag"]
}
```

## Argument Reference

* `consumer_id` - (Required) the id of the consumer to be configured with basic auth
* `username` - (Required) username to be used for basic auth
* `password` - (Required) password to be used for basic auth
* `tags` - (Optional) A list of strings associated with the consumer basic auth for grouping and filtering
