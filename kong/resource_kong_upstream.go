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

	upstream, err := meta.(*config).adminClient.Upstreams().Create(upstreamRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong upstream: %v error: %v", upstreamRequest, err)
	}

	d.SetId(upstream.Id)

	return resourceKongUpstreamRead(d, meta)
}

func resourceKongUpstreamUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	_, err := meta.(*config).adminClient.Upstreams().UpdateById(d.Id(), upstreamRequest)

	if err != nil {
		return fmt.Errorf("error updating kong upstream: %s", err)
	}

	return resourceKongUpstreamRead(d, meta)
}

func resourceKongUpstreamRead(d *schema.ResourceData, meta interface{}) error {

	upstream, err := meta.(*config).adminClient.Upstreams().GetById(d.Id())

	if err != nil {
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
		if err := d.Set("healthchecks", flattenHealthCheck(upstream.HealthChecks)); err != nil {
			return err
		}
	}

	return nil
}

func resourceKongUpstreamDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*config).adminClient.Upstreams().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong upstream: %v", err)
	}

	return nil
}

func createKongUpstreamRequestFromResourceData(d *schema.ResourceData) *gokong.UpstreamRequest {

	upstreamRequest := &gokong.UpstreamRequest{}

	upstreamRequest.Name = readStringFromResource(d, "name")
	upstreamRequest.Slots = readIntFromResource(d, "slots")
	upstreamRequest.HashOn = readStringFromResource(d, "hash_on")
	upstreamRequest.HashFallback = readStringFromResource(d, "hash_fallback")
	upstreamRequest.HashOnHeader = readStringFromResource(d, "hash_on_header")
	upstreamRequest.HashFallbackHeader = readStringFromResource(d, "hash_fallback_header")
	upstreamRequest.HashOnCookie = readStringFromResource(d, "hash_on_cookie")
	upstreamRequest.HashOnCookiePath = readStringFromResource(d, "hash_on_cookie_path")

	if healthChecksArray := readArrayFromResource(d, "healthchecks"); healthChecksArray != nil && len(healthChecksArray) > 0 {
		healthChecksMap := healthChecksArray[0].(map[string]interface{})
		upstreamRequest.HealthChecks = createKongHealthCheckFromMap(&healthChecksMap)
	}

	return upstreamRequest
}

func createKongHealthCheckFromMap(data *map[string]interface{}) *gokong.UpstreamHealthCheck {
	if data != nil {
		dataMap := *data
		healthCheck := &gokong.UpstreamHealthCheck{}

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

func createKongHealthCheckActiveFromMap(data *map[string]interface{}) *gokong.UpstreamHealthCheckActive {
	if data != nil {
		dataMap := *data
		active := &gokong.UpstreamHealthCheckActive{}

		if dataMap["type"] != nil {
			active.Type = dataMap["type"].(string)
		}
		if dataMap["timeout"] != nil {
			active.Timeout = dataMap["timeout"].(int)
		}
		if dataMap["concurrency"] != nil {
			active.Concurrency = dataMap["concurrency"].(int)
		}
		if dataMap["http_path"] != nil {
			active.HttpPath = dataMap["http_path"].(string)
		}
		if dataMap["https_verify_certificate"] != nil {
			active.HttpsVerifyCertificate = dataMap["https_verify_certificate"].(bool)
		}
		if dataMap["https_sni"] != nil {
			if len(dataMap["https_sni"].(string)) != 0 {
				active.HttpsSni = String(dataMap["https_sni"].(string))
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

func createKongHealthCheckPassiveFromMap(data *map[string]interface{}) *gokong.UpstreamHealthCheckPassive {
	if data != nil {
		dataMap := *data
		passive := &gokong.UpstreamHealthCheckPassive{}

		if dataMap["type"] != nil {
			passive.Type = dataMap["type"].(string)
		}

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

func createKongActiveHealthyFromMap(data *map[string]interface{}) *gokong.ActiveHealthy {
	if data != nil {
		dataMap := *data
		healthy := &gokong.ActiveHealthy{}

		if dataMap["interval"] != nil {
			healthy.Interval = dataMap["interval"].(int)
		}
		if dataMap["http_statuses"] != nil {
			healthy.HttpStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["successes"] != nil {
			healthy.Successes = dataMap["successes"].(int)
		}

		return healthy
	}
	return nil
}

func createKongActiveUnhealthyFromMap(data *map[string]interface{}) *gokong.ActiveUnhealthy {
	if data != nil {
		dataMap := *data
		unhealthy := &gokong.ActiveUnhealthy{}

		if dataMap["interval"] != nil {
			unhealthy.Interval = dataMap["interval"].(int)
		}
		if dataMap["http_statuses"] != nil {
			unhealthy.HttpStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["tcp_failures"] != nil {
			unhealthy.TcpFailures = dataMap["tcp_failures"].(int)
		}
		if dataMap["http_failures"] != nil {
			unhealthy.HttpFailures = dataMap["http_failures"].(int)
		}
		if dataMap["timeouts"] != nil {
			unhealthy.Timeouts = dataMap["timeouts"].(int)
		}

		return unhealthy
	}
	return nil
}

func createKongPassiveHealthyFromMap(data *map[string]interface{}) *gokong.PassiveHealthy {
	if data != nil {
		dataMap := *data
		healthy := &gokong.PassiveHealthy{}

		if dataMap["http_statuses"] != nil {
			healthy.HttpStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["successes"] != nil {
			healthy.Successes = dataMap["successes"].(int)
		}

		return healthy
	}
	return nil
}

func createKongPassiveUnhealthyFromMap(data *map[string]interface{}) *gokong.PassiveUnhealthy {
	if data != nil {
		dataMap := *data
		unhealthy := &gokong.PassiveUnhealthy{}

		if dataMap["http_statuses"] != nil {
			unhealthy.HttpStatuses = readIntArrayFromInterface(dataMap["http_statuses"])
		}
		if dataMap["tcp_failures"] != nil {
			unhealthy.TcpFailures = dataMap["tcp_failures"].(int)
		}
		if dataMap["http_failures"] != nil {
			unhealthy.HttpFailures = dataMap["http_failures"].(int)
		}
		if dataMap["timeouts"] != nil {
			unhealthy.Timeouts = dataMap["timeouts"].(int)
		}

		return unhealthy
	}
	return nil
}

func flattenHealthCheck(in *gokong.UpstreamHealthCheck) []interface{} {
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

func flattenHealthCheckActive(in *gokong.UpstreamHealthCheckActive) []interface{} {
	if in == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})

	m["type"] = in.Type
	m["timeout"] = in.Timeout
	m["concurrency"] = in.Concurrency
	m["http_path"] = in.HttpPath
	m["https_verify_certificate"] = in.HttpsVerifyCertificate

	if in.HttpsSni != nil {
		m["https_sni"] = *in.HttpsSni
	}
	if in.Healthy != nil {
		m["healthy"] = flattenActiveHealthy(in.Healthy)
	}
	if in.Unhealthy != nil {
		m["unhealthy"] = flattenActiveUnhealthy(in.Unhealthy)
	}

	return []interface{}{m}
}

func flattenHealthCheckPassive(in *gokong.UpstreamHealthCheckPassive) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	m["type"] = in.Type

	if in.Healthy != nil {
		m["healthy"] = flattenPassiveHealthy(in.Healthy)
	}
	if in.Unhealthy != nil {
		m["unhealthy"] = flattenPassiveUnhealthy(in.Unhealthy)
	}

	return []interface{}{m}
}

func flattenActiveHealthy(in *gokong.ActiveHealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = in.Interval
	m["http_statuses"] = in.HttpStatuses
	m["successes"] = in.Successes

	return []map[string]interface{}{m}
}

func flattenActiveUnhealthy(in *gokong.ActiveUnhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["interval"] = in.Interval
	m["http_statuses"] = in.HttpStatuses
	m["tcp_failures"] = in.TcpFailures
	m["http_failures"] = in.HttpFailures
	m["timeouts"] = in.Timeouts

	return []map[string]interface{}{m}
}

func flattenPassiveHealthy(in *gokong.PassiveHealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HttpStatuses
	m["successes"] = in.Successes

	return []map[string]interface{}{m}
}

func flattenPassiveUnhealthy(in *gokong.PassiveUnhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HttpStatuses
	m["tcp_failures"] = in.TcpFailures
	m["http_failures"] = in.HttpFailures
	m["timeouts"] = in.Timeouts

	return []map[string]interface{}{m}
}
