# kong_target

## Example Usage

```hcl
resource "kong_target" "target" {
    target  		= "sample_target:80"
    weight 	  	= 10
    upstream_id = "${kong_upstream.upstream.id}"
}
```

## Argument Reference

* `target` - (Required) is the target address (IP or hostname) and port. If omitted the port defaults to 8000.
* `weight` - (Required) is the weight this target gets within the upstream load balancer (0-1000, defaults to 100).
* `upstream_id` - (Required) is the id of the upstream to apply this target to.

## Import

To import a target use a combination of the upstream id and the target id as follows:

```shell
terraform import kong_target.<target_identifier> <upstream_id>/<target_id>
```
