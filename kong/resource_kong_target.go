package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongTargetCreate,
		ReadContext:   resourceKongTargetRead,
		DeleteContext: resourceKongTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"upstream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKongTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	targetRequest := createKongTargetRequestFromResourceData(d)

	client := meta.(*config).adminClient.Targets
	target, err := client.Create(ctx, readStringPtrFromResource(d, "upstream_id"), targetRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong target: %v error: %v", targetRequest, err))
	}

	d.SetId(IDToString(target.Upstream.ID) + "/" + *target.ID)

	return resourceKongTargetRead(ctx, d, meta)
}

func resourceKongTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var ids = strings.Split(d.Id(), "/")

	upstreamClient := meta.(*config).adminClient.Upstreams
	// First check if the upstream exists. If it does not then the target no longer exists either.
	if upstream, _ := upstreamClient.Get(ctx, kong.String(ids[0])); upstream == nil {
		d.SetId("")
		return diags
	}

	// TODO: Support paging
	client := meta.(*config).adminClient.Targets
	targets, _, err := client.List(ctx, kong.String(ids[0]), nil)

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong target: %v", err))
	}

	if targets == nil {
		d.SetId("")
	} else {
		for _, element := range targets {
			if *element.ID == ids[1] {
				err := d.Set("target", element.Target)
				if err != nil {
					return diag.FromErr(err)
				}
				err = d.Set("weight", element.Weight)
				if err != nil {
					return diag.FromErr(err)
				}
				err = d.Set("upstream_id", element.Upstream.ID)
				if err != nil {
					return diag.FromErr(err)
				}
				err = d.Set("tags", element.Tags)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return diags
}

func resourceKongTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var ids = strings.Split(d.Id(), "/")
	client := meta.(*config).adminClient.Targets
	if err := client.Delete(ctx, kong.String(ids[0]), kong.String(ids[1])); err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong target: %v", err))
	}

	return diags
}

func createKongTargetRequestFromResourceData(d *schema.ResourceData) *kong.Target {
	upstream := kong.Upstream{
		ID: readStringPtrFromResource(d, "upstream_id"),
	}
	return &kong.Target{
		Target:   readStringPtrFromResource(d, "target"),
		Weight:   readIntPtrFromResource(d, "weight"),
		Upstream: &upstream,
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}
}
