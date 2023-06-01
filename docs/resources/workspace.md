# kong_workspace

## Example Usage

```hcl
resource "kong_workspace" "workspace" {
    name  		= "myworkspace"
}
```

## Argument Reference

* `name` - (Required) The name of the Kong workspace.

## Import

Import workspaces using the workspace id:

```shell
terraform import kong_workspace.<workspace_identifier> <workspace_id>
```
