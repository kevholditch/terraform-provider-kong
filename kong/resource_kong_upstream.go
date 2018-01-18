package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongUpstreamCreate,
		Read:   resourceKongUpstreamRead,
		Delete: resourceKongUpstreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slots": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"order_list": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Default:  nil,
				ForceNew: true,
			},
		},
	}
}

func resourceKongUpstreamCreate(d *schema.ResourceData, meta interface{}) error {

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	upstream, err := meta.(*gokong.KongAdminClient).Upstreams().Create(upstreamRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong upstream: %v error: %v", upstreamRequest, err)
	}

	d.SetId(upstream.Id)

	return resourceKongUpstreamRead(d, meta)
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {

	upstream, err := meta.(*gokong.KongAdminClient).Upstreams().GetById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong upstream: %v", err)
	}

	d.Set("name", upstream.Name)
	d.Set("slots", upstream.Slots)

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*gokong.KongAdminClient).Upstreams().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong upstream: %v", err)
	}

	return nil
}

func createKongUpstreamRequestFromResourceData(d *schema.ResourceData) *gokong.UpstreamRequest {

	upstreamRequest := &gokong.UpstreamRequest{}

	upstreamRequest.Name = readStringFromResource(d, "name")
	upstreamRequest.Slots = readIntFromResource(d, "slots")
	upstreamRequest.OrderList = readIntArrayFromResource(d, "order_list")

	return upstreamRequest
}
