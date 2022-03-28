# kong_consumer_oauth2

Resource that allows you to configure the OAuth2 plugin credentials for a consumer.

## Example Usage

```hcl
resource "kong_consumer" "my_consumer" {
  username  = "User1"
  custom_id = "123"
}

resource "kong_plugin" "oauth2_plugin" {
	name = "oauth2"
	config_json = <<EOT
	{
		"global_credentials": true,
		"enable_password_grant": true,
		"token_expiration": 180,
		"refresh_token_ttl": 180,
		"provision_key": "testprovisionkey"
	}
EOT
}

resource "kong_consumer_oauth2" "consumer_oauth2" {
	name          = "test_application"
	consumer_id   = "${kong_consumer.my_consumer.id}"
	client_id     = "client_id"
	client_secret = "client_secret"
	redirect_uris = ["https://asdf.com/callback", "https://test.cl/callback"]
	tags          = ["myTag"]
}
```

## Argument Reference

* `name` - (Required) The name associated with the credential.
* `consumer_id` - (Required) The id of the consumer to be configured with oauth2.
* `client_id` - (Optional) Unique oauth2 client id. If not set, the oauth2 plugin will generate one
* `client_secret` - (Optional) Unique oauth2 client secret. If not set, the oauth2 plugin will generate one
* `hash_secret` - (Optional) A boolean flag that indicates whether the client_secret field will be stored in hashed form. If enabled on existing plugin instances, client secrets are hashed on the fly upon first usage. Default: `false`.
* `redirect_uris` - (Required) An array with one or more URLs in your app where users will be sent after authorization ([RFC 6742 Section 3.1.2](https://tools.ietf.org/html/rfc6749#section-3.1.2)).
* `tags` - (Optional) A list of strings associated with the consumer for grouping and filtering.
