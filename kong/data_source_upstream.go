package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func dataSourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKongUpstreamRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"order_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
		},
	}
}

func dataSourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {

	filter := &gokong.UpstreamFilter{}

	if v, _ := d.GetOk("filter"); v != nil {
		filterSet := v.(*schema.Set).List()
		if len(filterSet) == 1 {
			filterMap := filterSet[0].(map[string]interface{})
			filter.Id = filterMap["id"].(string)
			filter.Name = filterMap["name"].(string)
		}
	}

	results, err := meta.(*gokong.KongAdminClient).Upstreams().ListFiltered(filter)

	if err != nil {
		return fmt.Errorf("could not find upstream, error: %v", err)
	}

	if len(results.Results) == 0 {
		return fmt.Errorf("could not find upstream using filter: %v", filter)
	}

	if len(results.Results) > 1 {
		return fmt.Errorf("found more than 1 upstream make filter more restrictive")
	}

	upstream := results.Results[0]

	d.SetId(upstream.Id)
	d.Set("id", upstream.Id)
	d.Set("name", upstream.Name)
	d.Set("slots", upstream.Slots)
	d.Set("order_list", upstream.OrderList)

	return nil
}
