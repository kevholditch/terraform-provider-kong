package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hbagdi/go-kong/kong"
)

func resourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongUpstreamCreate,
		Read:   resourceKongUpstreamRead,
		Delete: resourceKongUpstreamDelete,
		Update: resourceKongUpstreamUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"slots": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  10000,
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
			"hash_on_cookie": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"hash_on_cookie_path": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "/",
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
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "http",
									},
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
									"https_verify_certificate": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"https_sni": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  nil,
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
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "http",
									},
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
		},
	}
}

func resourceKongUpstreamCreate(d *schema.ResourceData, meta interface{}) error {

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	client := meta.(*config).adminClient.Upstreams
	upstream, err := client.Create(context.Background(), upstreamRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong upstream: %v error: %v", upstreamRequest, err)
	}

	d.SetId(*upstream.ID)

	return resourceKongUpstreamRead(d, meta)
}

func resourceKongUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	client := meta.(*config).adminClient.Upstreams
	_, err := client.Update(context.Background(), upstreamRequest)

	if err != nil {
		return fmt.Errorf("error updating kong upstream: %s", err)
	}

	return resourceKongUpstreamRead(d, meta)
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Upstreams
	upstream, err := client.Get(context.Background(), kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("could not find kong upstream: %v", err)
	}

	if upstream == nil {
		d.SetId("")
	} else {
		d.Set("name", upstream.Name)
		d.Set("slots", upstream.Slots)
		d.Set("hash_on", upstream.HashOn)
		d.Set("hash_fallback", upstream.HashFallback)
		d.Set("hash_on_header", upstream.HashOnHeader)
		d.Set("hash_fallback_header", upstream.HashFallbackHeader)
		d.Set("hash_on_cookie", upstream.HashOnCookie)
		d.Set("hash_on_cookie_path", upstream.HashOnCookiePath)
		if err := d.Set("healthchecks", flattenHealthCheck(upstream.Healthchecks)); err != nil {
			return err
		}
	}

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Upstreams
	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong upstream: %v", err)
	}

	return nil
}

func createKongUpstreamRequestFromResourceData(d *schema.ResourceData) *kong.Upstream {

	upstreamRequest := &kong.Upstream{}

	upstreamRequest.Name = readStringPtrFromResource(d, "name")
	upstreamRequest.Slots = readIntPtrFromResource(d, "slots")
	upstreamRequest.HashOn = readStringPtrFromResource(d, "hash_on")
	upstreamRequest.HashFallback = readStringPtrFromResource(d, "hash_fallback")
	upstreamRequest.HashOnHeader = readStringPtrFromResource(d, "hash_on_header")
	upstreamRequest.HashFallbackHeader = readStringPtrFromResource(d, "hash_fallback_header")
	upstreamRequest.HashOnCookie = readStringPtrFromResource(d, "hash_on_cookie")
	upstreamRequest.HashOnCookiePath = readStringPtrFromResource(d, "hash_on_cookie_path")

	if healthChecksArray := readArrayFromResource(d, "healthchecks"); healthChecksArray != nil && len(healthChecksArray) > 0 {
		healthChecksMap := healthChecksArray[0].(map[string]interface{})
		upstreamRequest.Healthchecks = createKongHealthCheckFromMap(&healthChecksMap)
	}

	return upstreamRequest
}

func createKongHealthCheckFromMap(data *map[string]interface{}) *kong.Healthcheck {
	if data != nil {
		dataMap := *data
		healthCheck := &kong.Healthcheck{}

		if dataMap["active"] != nil {
			if activeArray := dataMap["active"].([]interface{}); activeArray != nil && len(activeArray) > 0 {
				activeMap := activeArray[0].(map[string]interface{})
				healthCheck.Active = createKongHealthCheckActiveFromMap(&activeMap)
			}
		}

		if dataMap["passive"] != nil {
			if passiveArray := dataMap["passive"].([]interface{}); passiveArray != nil && len(passiveArray) > 0 {
				passiveMap := passiveArray[0].(map[string]interface{})
				healthCheck.Passive = createKongHealthCheckPassiveFromMap(&passiveMap)
			}
		}

		return healthCheck
	}
	return nil
}

func createKongHealthCheckActiveFromMap(data *map[string]interface{}) *kong.ActiveHealthcheck {
	if data != nil {
		dataMap := *data
		active := &kong.ActiveHealthcheck{}

		if dataMap["type"] != nil {
			active.Type = kong.String(dataMap["type"].(string))
		}
		if dataMap["timeout"] != nil {
			active.Timeout = kong.Int(dataMap["timeout"].(int))
		}
		if dataMap["concurrency"] != nil {
			active.Concurrency = kong.Int(dataMap["concurrency"].(int))
		}
		if dataMap["http_path"] != nil {
			active.HTTPPath = kong.String(dataMap["http_path"].(string))
		}
		if dataMap["https_verify_certificate"] != nil {
			active.HTTPSVerifyCertificate = kong.Bool(dataMap["https_verify_certificate"].(bool))
		}
		if dataMap["https_sni"] != nil {
			if len(dataMap["https_sni"].(string)) != 0 {
				active.HTTPSSni = kong.String(dataMap["https_sni"].(string))
			}
		}

		if dataMap["healthy"] != nil {
			if healthyArray := dataMap["healthy"].([]interface{}); healthyArray != nil && len(healthyArray) > 0 {
				healthyMap := healthyArray[0].(map[string]interface{})
				active.Healthy = createKongActiveHealthyFromMap(&healthyMap)
			}
		}

		if dataMap["unhealthy"] != nil {
			if unhealthyArray := dataMap["unhealthy"].([]interface{}); unhealthyArray != nil && len(unhealthyArray) > 0 {
				unhealthyMap := unhealthyArray[0].(map[string]interface{})
				active.Unhealthy = createKongActiveUnhealthyFromMap(&unhealthyMap)
			}
		}

		return active
	}

	return nil
}

func createKongHealthCheckPassiveFromMap(data *map[string]interface{}) *kong.PassiveHealthcheck {
	if data != nil {
		dataMap := *data
		passive := &kong.PassiveHealthcheck{}

		// if dataMap["type"] != nil {
		// 	passive.Type = dataMap["type"].(string)
		// }

		if dataMap["healthy"] != nil {
			if healthyArray := dataMap["healthy"].([]interface{}); healthyArray != nil && len(healthyArray) > 0 {
				healthyMap := healthyArray[0].(map[string]interface{})
				passive.Healthy = createKongPassiveHealthyFromMap(&healthyMap)
			}
		}

		if dataMap["unhealthy"] != nil {
			if unhealthyArray := dataMap["unhealthy"].([]interface{}); unhealthyArray != nil && len(unhealthyArray) > 0 {
				unhealthyMap := unhealthyArray[0].(map[string]interface{})
				passive.Unhealthy = createKongPassiveUnhealthyFromMap(&unhealthyMap)
			}
		}

		return passive
	}

	return nil
}

func createKongActiveHealthyFromMap(data *map[string]interface{}) *kong.Healthy {
	if data != nil {
		dataMap := *data
		healthy := &kong.Healthy{}

		if dataMap["interval"] != nil {
			healthy.Interval = kong.Int(dataMap["interval"].(int))
		}
		if dataMap["http_statuses"] != nil {
			healthy.HTTPStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["successes"] != nil {
			healthy.Successes = kong.Int(dataMap["successes"].(int))
		}

		return healthy
	}
	return nil
}

func createKongActiveUnhealthyFromMap(data *map[string]interface{}) *kong.Unhealthy {
	if data != nil {
		dataMap := *data
		unhealthy := &kong.Unhealthy{}

		if dataMap["interval"] != nil {
			unhealthy.Interval = kong.Int(dataMap["interval"].(int))
		}
		if dataMap["http_statuses"] != nil {
			unhealthy.HTTPStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["tcp_failures"] != nil {
			unhealthy.TCPFailures = kong.Int(dataMap["tcp_failures"].(int))
		}
		if dataMap["http_failures"] != nil {
			unhealthy.HTTPFailures = kong.Int(dataMap["http_failures"].(int))
		}
		if dataMap["timeouts"] != nil {
			unhealthy.Timeouts = kong.Int(dataMap["timeouts"].(int))
		}

		return unhealthy
	}
	return nil
}

func createKongPassiveHealthyFromMap(data *map[string]interface{}) *kong.Healthy {
	if data != nil {
		dataMap := *data
		healthy := &kong.Healthy{}

		if dataMap["http_statuses"] != nil {
			healthy.HTTPStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["successes"] != nil {
			healthy.Successes = kong.Int(dataMap["successes"].(int))
		}

		return healthy
	}
	return nil
}

func createKongPassiveUnhealthyFromMap(data *map[string]interface{}) *kong.Unhealthy {
	if data != nil {
		dataMap := *data
		unhealthy := &kong.Unhealthy{}

		if dataMap["http_statuses"] != nil {
			unhealthy.HTTPStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["tcp_failures"] != nil {
			unhealthy.TCPFailures = kong.Int(dataMap["tcp_failures"].(int))
		}
		if dataMap["http_failures"] != nil {
			unhealthy.HTTPFailures = kong.Int(dataMap["http_failures"].(int))
		}
		if dataMap["timeouts"] != nil {
			unhealthy.Timeouts = kong.Int(dataMap["timeouts"].(int))
		}

		return unhealthy
	}
	return nil
}

func flattenHealthCheck(in *kong.Healthcheck) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if in.Active != nil {
		m["active"] = flattenHealthCheckActive(in.Active)
	}
	if in.Passive != nil {
		m["passive"] = flattenHealthCheckPassive(in.Passive)
	}

	return []interface{}{m}
}

func flattenHealthCheckActive(in *kong.ActiveHealthcheck) []interface{} {
	if in == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})

	m["type"] = in.Type
	m["timeout"] = in.Timeout
	m["concurrency"] = in.Concurrency
	m["http_path"] = in.HTTPPath
	m["https_verify_certificate"] = in.HTTPSVerifyCertificate

	if in.HTTPSSni != nil {
		m["https_sni"] = *in.HTTPSSni
	}
	if in.Healthy != nil {
		m["healthy"] = flattenActiveHealthy(in.Healthy)
	}
	if in.Unhealthy != nil {
		m["unhealthy"] = flattenActiveUnhealthy(in.Unhealthy)
	}

	return []interface{}{m}
}

func flattenHealthCheckPassive(in *kong.PassiveHealthcheck) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	// Requires the merge of https://github.com/hbagdi/go-kong/pull/22
	// m["type"] = in.Type

	if in.Healthy != nil {
		m["healthy"] = flattenPassiveHealthy(in.Healthy)
	}
	if in.Unhealthy != nil {
		m["unhealthy"] = flattenPassiveUnhealthy(in.Unhealthy)
	}

	return []interface{}{m}
}

func flattenActiveHealthy(in *kong.Healthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = in.Interval
	m["http_statuses"] = in.HTTPStatuses
	m["successes"] = in.Successes

	return []map[string]interface{}{m}
}

func flattenActiveUnhealthy(in *kong.Unhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = in.Interval
	m["http_statuses"] = in.HTTPStatuses
	m["tcp_failures"] = in.TCPFailures
	m["http_failures"] = in.HTTPFailures
	m["timeouts"] = in.Timeouts

	return []map[string]interface{}{m}
}

func flattenPassiveHealthy(in *kong.Healthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HTTPStatuses
	m["successes"] = in.Successes

	return []map[string]interface{}{m}
}

func flattenPassiveUnhealthy(in *kong.Unhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HTTPStatuses
	m["tcp_failures"] = in.TCPFailures
	m["http_failures"] = in.HTTPFailures
	m["timeouts"] = in.Timeouts

	return []map[string]interface{}{m}
}
