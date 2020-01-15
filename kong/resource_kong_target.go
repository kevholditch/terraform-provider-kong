package kong

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
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

	target, err := meta.(*config).adminClient.Targets().CreateFromUpstreamId(readStringFromResource(d, "upstream_id"), targetRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong target: %v error: %v", targetRequest, err)
	}

	d.SetId(gokong.IdToString(target.Upstream) + "/" + *target.Id)

	return resourceKongTargetRead(d, meta)
}

func resourceKongTargetRead(d *schema.ResourceData, meta interface{}) error {

	var ids = strings.Split(d.Id(), "/")

	// First check if the upstream exists. If it does not then the target no longer exists either.
	if upstream, _ := meta.(*config).adminClient.Upstreams().GetById(ids[0]); upstream == nil {
		d.SetId("")
		return nil
	}

	targets, err := meta.(*config).adminClient.Targets().GetTargetsFromUpstreamId(ids[0])

	if err != nil {
		return fmt.Errorf("could not find kong target: %v", err)
	}

	if targets == nil {
		d.SetId("")
	} else {
		for _, element := range targets {
			if *element.Id == ids[1] {
				d.Set("target", element.Target)
				d.Set("weight", element.Weight)
				d.Set("upstream_id", element.Upstream)
			}
		}
	}

	return nil
}

func resourceKongTargetDelete(d *schema.ResourceData, meta interface{}) error {

	var ids = strings.Split(d.Id(), "/")
	if err := meta.(*config).adminClient.Targets().DeleteFromUpstreamById(ids[0], ids[1]); err != nil {
		return fmt.Errorf("could not delete kong target: %v", err)
	}

	return nil
}

func createKongTargetRequestFromResourceData(d *schema.ResourceData) *gokong.TargetRequest {
	return &gokong.TargetRequest{
		Target: readStringFromResource(d, "target"),
		Weight: readIntFromResource(d, "weight"),
	}
}
