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
			"hash_on": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "none",
			},
			"hash_fallback": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "none",
			},
			"hash_on_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"hash_fallback_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"healthchecks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"concurrency": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Default:  10,
									},
									"http_path": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "/",
									},
									"healthy": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"unhealthy": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_failures": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"timeouts": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"passive": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"healthy": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"unhealthy": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"http_failures": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
												"timeouts": &schema.Schema{
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
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

	results, err := meta.(*config).adminClient.Upstreams().ListFiltered(filter)

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
	d.Set("hash_on", upstream.HashOn)
	d.Set("hash_fallback", upstream.HashFallback)
	d.Set("hash_on_header", upstream.HashOnHeader)
	d.Set("hash_fallback_header", upstream.HashFallbackHeader)
	if err := d.Set("healthchecks", flattenHealthCheck(upstream.HealthChecks)); err != nil {
		return err
	}

	return nil
}
