package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongRouteCreate,
		ReadContext:   resourceKongRouteRead,
		DeleteContext: resourceKongRouteDelete,
		UpdateContext: resourceKongRouteUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"protocols": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"methods": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"paths": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"strip_path": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"source": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"destination": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"snis": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"preserve_host": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"regex_priority": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"path_handling": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "v0",
			},
			"https_redirect_status_code": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  426,
			},
			"request_buffering": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"response_buffering": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"header": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceKongRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	routeRequest := createKongRouteRequestFromResourceData(d)

	client := meta.(*config).adminClient.Routes
	route, err := client.Create(ctx, routeRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong route: %v error: %v", routeRequest, err))
	}

	d.SetId(*route.ID)

	return resourceKongRouteRead(ctx, d, meta)
}

func resourceKongRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	routeRequest := createKongRouteRequestFromResourceData(d)

	client := meta.(*config).adminClient.Routes

	_, err := client.Update(ctx, routeRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong route: %s", err))
	}

	return resourceKongRouteRead(ctx, d, meta)
}

func resourceKongRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Routes
	route, err := client.Get(ctx, kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong route: %v", err))
	}

	if route == nil {
		d.SetId("")
	} else {
		if route.Name != nil {
			err := d.Set("name", route.Name)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if route.Protocols != nil {
			err := d.Set("protocols", StringValueSlice(route.Protocols))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Methods != nil {
			err := d.Set("methods", StringValueSlice(route.Methods))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Hosts != nil {
			err := d.Set("hosts", StringValueSlice(route.Hosts))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Paths != nil {
			err := d.Set("paths", StringValueSlice(route.Paths))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.StripPath != nil {
			err := d.Set("strip_path", route.StripPath)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Sources != nil {
			err := d.Set("source", flattenIpCidrArray(route.Sources))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Destinations != nil {
			err := d.Set("destination", flattenIpCidrArray(route.Destinations))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.PreserveHost != nil {
			err := d.Set("preserve_host", route.PreserveHost)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.RegexPriority != nil {
			err := d.Set("regex_priority", route.RegexPriority)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.SNIs != nil {
			err := d.Set("snis", StringValueSlice(route.SNIs))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.Service != nil {
			err := d.Set("service_id", route.Service.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.PathHandling != nil {
			err := d.Set("path_handling", route.PathHandling)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.HTTPSRedirectStatusCode != nil {
			err := d.Set("https_redirect_status_code", route.HTTPSRedirectStatusCode)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.RequestBuffering != nil {
			err := d.Set("request_buffering", route.RequestBuffering)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if route.ResponseBuffering != nil {
			err := d.Set("response_buffering", route.ResponseBuffering)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		err = d.Set("tags", route.Tags)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	return diags
}
func flattenIpCidrArray(addresses []*kong.CIDRPort) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(addresses), len(addresses))
	for i, v := range addresses {
		m := make(map[string]interface{})
		if v.IP != nil {
			m["ip"] = v.IP
		}
		if v.Port != nil {
			m["port"] = v.Port
		}
		out[i] = m
	}
	return out
}

func resourceKongRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Routes
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong route: %v", err))
	}

	return diags
}

func createKongRouteRequestFromResourceData(d *schema.ResourceData) *kong.Route {

	route := &kong.Route{
		Name:          readStringPtrFromResource(d, "name"),
		Protocols:     readStringArrayPtrFromResource(d, "protocols"),
		Methods:       readStringArrayPtrFromResource(d, "methods"),
		Hosts:         readStringArrayPtrFromResource(d, "hosts"),
		Paths:         readStringArrayPtrFromResource(d, "paths"),
		StripPath:     readBoolPtrFromResource(d, "strip_path"),
		Sources:       readIpPortArrayFromResource(d, "source"),
		Destinations:  readIpPortArrayFromResource(d, "destination"),
		PreserveHost:  readBoolPtrFromResource(d, "preserve_host"),
		RegexPriority: readIntPtrFromResource(d, "regex_priority"),
		SNIs:          readStringArrayPtrFromResource(d, "snis"),
		Service: &kong.Service{
			ID: readIdPtrFromResource(d, "service_id"),
		},
		PathHandling:            readStringPtrFromResource(d, "path_handling"),
		HTTPSRedirectStatusCode: readIntPtrFromResource(d, "https_redirect_status_code"),
		RequestBuffering:        readBoolPtrFromResource(d, "request_buffering"),
		ResponseBuffering:       readBoolPtrFromResource(d, "response_buffering"),
		Tags:                    readStringArrayPtrFromResource(d, "tags"),
		Headers:                 readMapStringArrayFromResource(d, "header"),
	}
	if d.Id() != "" {
		route.ID = kong.String(d.Id())
	}
	return route
}
