package kong

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongTargetCreate,
		Read:   resourceKongTargetRead,
		Delete: resourceKongTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"upstream_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKongTargetCreate(d *schema.ResourceData, meta interface{}) error {

	targetRequest := createKongTargetRequestFromResourceData(d)

	client := meta.(*config).adminClient.Targets
	target, err := client.Create(context.Background(), readStringPtrFromResource(d, "upstream_id"), targetRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong target: %v error: %v", targetRequest, err)
	}

	d.SetId(IDToString(target.Upstream.ID) + "/" + *target.ID)

	return resourceKongTargetRead(d, meta)
}

func resourceKongTargetRead(d *schema.ResourceData, meta interface{}) error {

	var ids = strings.Split(d.Id(), "/")

	upstreamClient := meta.(*config).adminClient.Upstreams
	// First check if the upstream exists. If it does not then the target no longer exists either.
	if upstream, _ := upstreamClient.Get(context.Background(), kong.String(ids[0])); upstream == nil {
		d.SetId("")
		return nil
	}

	// TODO: Support paging
	client := meta.(*config).adminClient.Targets
	targets, _, err := client.List(context.Background(), kong.String(ids[0]), nil)

	if err != nil {
		return fmt.Errorf("could not find kong target: %v", err)
	}

	if targets == nil {
		d.SetId("")
	} else {
		for _, element := range targets {
			if *element.ID == ids[1] {
				d.Set("target", element.Target)
				d.Set("weight", element.Weight)
				d.Set("upstream_id", element.Upstream.ID)
			}
		}
	}

	return nil
}

func resourceKongTargetDelete(d *schema.ResourceData, meta interface{}) error {

	var ids = strings.Split(d.Id(), "/")
	client := meta.(*config).adminClient.Targets
	if err := client.Delete(context.Background(), kong.String(ids[0]), kong.String(ids[1])); err != nil {
		return fmt.Errorf("could not delete kong target: %v", err)
	}

	return nil
}

func createKongTargetRequestFromResourceData(d *schema.ResourceData) *kong.Target {
	upstream := kong.Upstream{
		ID: readStringPtrFromResource(d, "upstream_id"),
	}
	return &kong.Target{
		Target:   readStringPtrFromResource(d, "target"),
		Weight:   readIntPtrFromResource(d, "weight"),
		Upstream: &upstream,
	}
}
