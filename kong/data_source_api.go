package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func dataSourceKongApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKongApiRead,
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
						"upstream_url": {
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
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"uris": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"methods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"upstream_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"strip_uri": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"preserve_host": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"retries": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upstream_connect_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upstream_send_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upstream_read_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"https_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"http_if_terminated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceKongApiRead(d *schema.ResourceData, meta interface{}) error {

	filter := &gokong.ApiFilter{}

	if v, _ := d.GetOk("filter"); v != nil {
		filterSet := v.(*schema.Set).List()
		if len(filterSet) == 1 {
			filterMap := filterSet[0].(map[string]interface{})
			filter.Id = filterMap["id"].(string)
			filter.Name = filterMap["name"].(string)
			filter.UpstreamUrl = filterMap["upstream_url"].(string)
		}
	}

	results, err := meta.(*gokong.KongAdminClient).Apis().ListFiltered(filter)

	if err != nil {
		return fmt.Errorf("could not find api, error: %v", err)
	}

	if len(results.Results) == 0 {
		return fmt.Errorf("could not find api using filter: %v", filter)
	}

	if len(results.Results) > 1 {
		return fmt.Errorf("found more than 1 api make filter more restrictive")
	}

	api := results.Results[0]

	if api == nil {
		return nil
	}

	d.SetId(*api.Id)

	if api.Id != nil {
		d.Set("id", api.Id)
	}

	if api.Name != nil {
		d.Set("name", api.Name)
	}

	if api.Hosts != nil {
		d.Set("hosts", api.Hosts)
	}

	if api.Uris != nil {
		d.Set("uris", api.Uris)
	}

	if api.Methods != nil {
		d.Set("methods", api.Methods)
	}

	if api.UpstreamUrl != nil {
		d.Set("upstream_url", api.UpstreamUrl)
	}

	if api.StripUri != nil {
		d.Set("strip_uri", api.StripUri)
	}

	if api.PreserveHost != nil {
		d.Set("preserve_host", api.PreserveHost)
	}

	if api.Retries != nil {
		d.Set("retries", api.Retries)
	}

	if api.UpstreamConnectTimeout != nil {
		d.Set("upstream_connect_timeout", api.UpstreamConnectTimeout)
	}

	if api.UpstreamSendTimeout != nil {
		d.Set("upstream_send_timeout", api.UpstreamSendTimeout)
	}

	if api.UpstreamReadTimeout != nil {
		d.Set("upstream_read_timeout", api.UpstreamReadTimeout)
	}

	if api.HttpsOnly != nil {
		d.Set("https_only", api.HttpsOnly)
	}

	if api.HttpIfTerminated != nil {
		d.Set("http_if_terminated", api.HttpIfTerminated)
	}

	return nil
}
