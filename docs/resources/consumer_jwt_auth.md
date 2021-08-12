# kong_consumer_jwt_auth

Consumer jwt auth is a resource that allows you to configure the jwt auth plugin for a consumer.

## Example Usage

```hcl
resource "kong_consumer" "my_consumer" {
  username  = "User1"
  custom_id = "123"
}

resource "kong_plugin" "jwt_plugin" {
  name        = "jwt"
  config_json = <<EOT
	{
		"claims_to_verify": ["exp"]
	}
EOT
}

resource "kong_consumer_jwt_auth" "consumer_jwt_config" {
  consumer_id    = "${kong_consumer.my_consumer.id}"
  algorithm      = "HS256"
  key            = "my_key"
  rsa_public_key = "foo"
  secret         = "my_secret"
}
```

## Argument Reference

* `consumer_id` - (Required) the id of the consumer to be configured with jwt auth
* `algorithm` - (Optional) The algorithm used to verify the token’s signature. Can be HS256, HS384, HS512, RS256, or ES256, Default is `HS256`
* `key` - (Optional) A unique string identifying the credential. If left out, it will be auto-generated.
* `rsa_public_key` - (Optional) If algorithm is `RS256` or `ES256`, the public key (in PEM format) to use to verify the token’s signature
* `secret` - (Optional) If algorithm is `HS256` or `ES256`, the secret used to sign JWTs for this credential. If left out, will be auto-generated
* `tags` - (Optional) A list of strings associated with the consumer JWT auth for grouping and filtering
