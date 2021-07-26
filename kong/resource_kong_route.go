package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongRouteCreate,
		Read:   resourceKongRouteRead,
		Delete: resourceKongRouteDelete,
		Update: resourceKongRouteUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"protocols": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"methods": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"paths": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"strip_path": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"source": &schema.Schema{
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
			"destination": &schema.Schema{
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
			"snis": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"preserve_host": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"regex_priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceKongRouteCreate(d *schema.ResourceData, meta interface{}) error {

	routeRequest := createKongRouteRequestFromResourceData(d)

	client := meta.(*config).adminClient.Routes
	route, err := client.Create(context.Background(), routeRequest)
	if err != nil {
		return fmt.Errorf("failed to create kong route: %v error: %v", routeRequest, err)
	}

	d.SetId(*route.ID)

	return resourceKongRouteRead(d, meta)
}

func resourceKongRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	routeRequest := createKongRouteRequestFromResourceData(d)

	client := meta.(*config).adminClient.Routes

	_, err := client.Update(context.Background(), routeRequest)

	if err != nil {
		return fmt.Errorf("error updating kong route: %s", err)
	}

	return resourceKongRouteRead(d, meta)
}

func resourceKongRouteRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Routes
	route, err := client.Get(context.Background(), kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("could not find kong route: %v", err)
	}

	if route == nil {
		d.SetId("")
	} else {
		if route.Name != nil {
			d.Set("name", route.Name)
		}
		if route.Protocols != nil {
			d.Set("protocols", StringValueSlice(route.Protocols))
		}

		if route.Methods != nil {
			d.Set("methods", StringValueSlice(route.Methods))
		}

		if route.Hosts != nil {
			d.Set("hosts", StringValueSlice(route.Hosts))
		}

		if route.Paths != nil {
			d.Set("paths", StringValueSlice(route.Paths))
		}

		if route.StripPath != nil {
			d.Set("strip_path", route.StripPath)
		}

		if route.Sources != nil {
			d.Set("source", flattenIpCidrArray(route.Sources))
		}

		if route.Destinations != nil {
			d.Set("destination", flattenIpCidrArray(route.Destinations))
		}

		if route.PreserveHost != nil {
			d.Set("preserve_host", route.PreserveHost)
		}

		if route.RegexPriority != nil {
			d.Set("regex_priority", route.RegexPriority)
		}

		if route.SNIs != nil {
			d.Set("snis", StringValueSlice(route.SNIs))
		}

		if route.Service != nil {
			d.Set("service_id", route.Service.ID)
		}

	}

	return nil
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

func resourceKongRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Routes
	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong route: %v", err)
	}

	return nil
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
	}
	if d.Id() != "" {
		route.ID = kong.String(d.Id())
	}
	return route
}
