package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongWorkspaceCreate,
		ReadContext:   resourceKongWorkspaceRead,
		DeleteContext: resourceKongWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceKongWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	workspaceRequest := createKongWorkspaceRequestFromResourceData(d)

	client := meta.(*config).adminClient.Workspaces
	workspace, err := client.Create(ctx, workspaceRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong workspace: %v error: %v", workspaceRequest, err))
	}

	d.SetId(*workspace.ID)

	return resourceKongWorkspaceRead(ctx, d, meta)
}

func resourceKongWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Workspaces
	workspace, err := client.Get(ctx, kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong workspace: %v", err))
	}

	if workspace == nil {
		d.SetId("")
	} else {
		if workspace.Name != nil {
			err := d.Set("name", workspace.Name)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func resourceKongWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Workspaces
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong workspace: %v", err))
	}

	return diags
}

func createKongWorkspaceRequestFromResourceData(d *schema.ResourceData) *kong.Workspace {

	workspace := &kong.Workspace{
		Name: readStringPtrFromResource(d, "name"),
	}
	if d.Id() != "" {
		workspace.ID = kong.String(d.Id())
	}
	return workspace
}
