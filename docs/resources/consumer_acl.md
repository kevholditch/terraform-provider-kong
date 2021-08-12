# kong_consumer_acl

Consumer ACL is a resource that allows you to configure the acl plugin for a consumer.

## Example Usage

```hcl
resource "kong_consumer" "my_consumer" {
	username  = "User1"
	custom_id = "123"
}

resource "kong_plugin" "acl_plugin" {
	name        = "acl"
	config_json = <<EOT
	{
		"allow": ["group1", "group2"]
	}
EOT
}

resource "kong_consumer_acl" "consumer_acl" {
	consumer_id    = "${kong_consumer.my_consumer.id}"
	group          = "group2"
	tags           = ["myTag", "otherTag"]
}
```

## Argument Reference

* `consumer_id` - (Required) the id of the consumer to be configured
* `group` - (Required) the acl group
* `tags` - (Optional) A list of strings associated with the consumer acl for grouping and filtering
