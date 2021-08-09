package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongUpstream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongUpstreamCreate,
		ReadContext:   resourceKongUpstreamRead,
		DeleteContext: resourceKongUpstreamDelete,
		UpdateContext: resourceKongUpstreamUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"slots": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  10000,
			},
			"hash_on": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "none",
			},
			"hash_fallback": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "none",
			},
			"hash_on_header": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"host_header": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"hash_fallback_header": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"hash_on_cookie": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"hash_on_cookie_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "/",
			},
			"healthchecks": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "http",
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  1,
									},
									"concurrency": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  10,
									},
									"http_path": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "/",
									},
									"https_verify_certificate": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"https_sni": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  nil,
									},
									"healthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"unhealthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"interval": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"http_failures": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"timeouts": {
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
						"passive": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "http",
									},
									"healthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"successes": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"unhealthy": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_statuses": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"tcp_failures": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"http_failures": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"timeouts": {
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

func resourceKongUpstreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	client := meta.(*config).adminClient.Upstreams
	upstream, err := client.Create(ctx, upstreamRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong upstream: %v error: %v", upstreamRequest, err))
	}

	d.SetId(*upstream.ID)

	return resourceKongUpstreamRead(ctx, d, meta)
}

func resourceKongUpstreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	upstreamRequest := createKongUpstreamRequestFromResourceData(d)

	client := meta.(*config).adminClient.Upstreams
	_, err := client.Update(ctx, upstreamRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong upstream: %s", err))
	}

	return resourceKongUpstreamRead(ctx, d, meta)
}

func resourceKongUpstreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Upstreams
	upstream, err := client.Get(ctx, kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong upstream: %v", err))
	}

	if upstream == nil {
		d.SetId("")
	} else {
		d.SetId(*upstream.ID)
		err := d.Set("name", upstream.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("slots", upstream.Slots)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_on", upstream.HashOn)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_fallback", upstream.HashFallback)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_on_header", upstream.HashOnHeader)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_fallback_header", upstream.HashFallbackHeader)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_on_cookie", upstream.HashOnCookie)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_on_cookie_path", upstream.HashOnCookiePath)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("healthchecks", flattenHealthCheck(upstream.Healthchecks)); err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", upstream.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
		if upstream.HostHeader != nil {
			err = d.Set("host_header", upstream.HostHeader)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if upstream.ClientCertificate != nil {
			err = d.Set("client_certificate_id", upstream.ClientCertificate.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func resourceKongUpstreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Upstreams
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong upstream: %v", err))
	}

	return diags
}

func createKongUpstreamRequestFromResourceData(d *schema.ResourceData) *kong.Upstream {

	upstreamRequest := &kong.Upstream{}

	if d.Id() != "" {
		upstreamRequest.ID = kong.String(d.Id())
	}
	upstreamRequest.Name = readStringPtrFromResource(d, "name")
	upstreamRequest.Slots = readIntPtrFromResource(d, "slots")
	upstreamRequest.HashOn = readStringPtrFromResource(d, "hash_on")
	upstreamRequest.HashFallback = readStringPtrFromResource(d, "hash_fallback")
	upstreamRequest.HashOnHeader = readStringPtrFromResource(d, "hash_on_header")
	upstreamRequest.HashFallbackHeader = readStringPtrFromResource(d, "hash_fallback_header")
	upstreamRequest.HashOnCookie = readStringPtrFromResource(d, "hash_on_cookie")
	upstreamRequest.HashOnCookiePath = readStringPtrFromResource(d, "hash_on_cookie_path")
	upstreamRequest.HostHeader = readStringPtrFromResource(d, "host_header")
	upstreamRequest.Tags = readStringArrayPtrFromResource(d, "tags")

	clientCertificateID := readIdPtrFromResource(d, "client_certificate_id")
	if clientCertificateID != nil {
		upstreamRequest.ClientCertificate = &kong.Certificate{
			ID: clientCertificateID,
		}
	}

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

		if dataMap["type"] != nil {
			passive.Type = kong.String(dataMap["type"].(string))
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

	if in.Type != nil {
		m["type"] = *in.Type
	}
	if in.Timeout != nil {
		m["timeout"] = *in.Timeout
	}
	if in.Concurrency != nil {
		m["concurrency"] = *in.Concurrency
	}
	if in.HTTPPath != nil {
		m["http_path"] = *in.HTTPPath
	}
	if in.HTTPSVerifyCertificate != nil {
		m["https_verify_certificate"] = *in.HTTPSVerifyCertificate
	}

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

	if in.Type != nil {
		m["type"] = *in.Type
	}

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

	if in.Interval != nil {
		m["interval"] = *in.Interval
	}
	m["http_statuses"] = in.HTTPStatuses
	if in.Successes != nil {
		m["successes"] = *in.Successes
	}

	return []map[string]interface{}{m}
}

func flattenActiveUnhealthy(in *kong.Unhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	if in.Interval != nil {
		m["interval"] = *in.Interval
	}
	m["http_statuses"] = in.HTTPStatuses
	if in.TCPFailures != nil {
		m["tcp_failures"] = *in.TCPFailures
	}
	if in.HTTPFailures != nil {
		m["http_failures"] = *in.HTTPFailures
	}
	if in.Timeouts != nil {
		m["timeouts"] = *in.Timeouts
	}

	return []map[string]interface{}{m}
}

func flattenPassiveHealthy(in *kong.Healthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HTTPStatuses
	if in.Successes != nil {
		m["successes"] = *in.Successes
	}

	return []map[string]interface{}{m}
}

func flattenPassiveUnhealthy(in *kong.Unhealthy) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	m := make(map[string]interface{})

	m["http_statuses"] = in.HTTPStatuses
	if in.TCPFailures != nil {
		m["tcp_failures"] = *in.TCPFailures
	}
	if in.HTTPFailures != nil {
		m["http_failures"] = *in.HTTPFailures
	}
	if in.Timeouts != nil {
		m["timeouts"] = *in.Timeouts
	}

	return []map[string]interface{}{m}
}
