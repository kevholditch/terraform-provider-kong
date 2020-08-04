package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func dataSourceKongConsumer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKongConsumerRead,
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
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_id": {
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
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKongConsumerRead(d *schema.ResourceData, meta interface{}) error {

	filter := &gokong.ConsumerFilter{}

	if v, _ := d.GetOk("filter"); v != nil {
		filterSet := v.(*schema.Set).List()
		if len(filterSet) == 1 {
			filterMap := filterSet[0].(map[string]interface{})
			filter.Id = filterMap["id"].(string)
			filter.Username = filterMap["username"].(string)
			filter.CustomId = filterMap["custom_id"].(string)
		}
	}

	results, err := meta.(*config).adminClient.Consumers().ListFiltered(filter)

	if err != nil {
		return fmt.Errorf("could not find consumer, error: %v", err)
	}

	if len(results.Results) == 0 {
		return fmt.Errorf("could not find consumer using filter: %v", filter)
	}

	if len(results.Results) > 1 {
		return fmt.Errorf("found more than 1 consumer make filter more restrictive")
	}

	consumer := results.Results[0]

	d.SetId(consumer.Id)
	d.Set("id", consumer.Id)
	d.Set("username", consumer.Username)
	d.Set("custom_id", consumer.CustomId)

	return nil
}
