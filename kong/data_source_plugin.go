package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func dataSourceKongPlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKongPluginRead,
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
						"api_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"consumer_id": {
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
			"api_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {

	filter := &gokong.PluginFilter{}

	if v, _ := d.GetOk("filter"); v != nil {
		filterSet := v.(*schema.Set).List()
		if len(filterSet) == 1 {
			filterMap := filterSet[0].(map[string]interface{})
			filter.Id = filterMap["id"].(string)
			filter.Name = filterMap["name"].(string)
			filter.ApiId = filterMap["api_id"].(string)
			filter.ConsumerId = filterMap["consumer_id"].(string)
		}
	}

	results, err := meta.(*gokong.KongAdminClient).Plugins().ListFiltered(filter)

	if err != nil {
		return fmt.Errorf("could not find plugin, error: %v", err)
	}

	if len(results.Results) == 0 {
		return fmt.Errorf("could not find plugin using filter: %v", filter)
	}

	if len(results.Results) > 1 {
		return fmt.Errorf("found more than 1 plugin make filter more restrictive")
	}

	plugin := results.Results[0]

	d.SetId(plugin.Id)
	d.Set("id", plugin.Id)
	d.Set("name", plugin.Name)
	d.Set("api_id", plugin.ApiId)
	d.Set("consumer_id", plugin.ConsumerId)
	d.Set("enabled", plugin.Enabled)

	return nil
}
