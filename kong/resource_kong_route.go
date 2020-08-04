package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
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
			"preserve_host": &schema.Schema{
				Type:     schema.TypeBool,
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

	route, err := meta.(*config).adminClient.Routes().AddRoute(routeRequest)
	if err != nil {
		return fmt.Errorf("failed to create kong route: %v error: %v", routeRequest, err)
	}

	d.SetId(*route.Id)

	return resourceKongRouteRead(d, meta)
}

func resourceKongRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	routeRequest := createKongRouteRequestFromResourceData(d)

	_, err := meta.(*config).adminClient.Routes().UpdateRoute(d.Id(), routeRequest)

	if err != nil {
		return fmt.Errorf("error updating kong route: %s", err)
	}

	return resourceKongRouteRead(d, meta)
}

func resourceKongRouteRead(d *schema.ResourceData, meta interface{}) error {

	route, err := meta.(*config).adminClient.Routes().GetRoute(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong route: %v", err)
	}

	if route == nil {
		d.SetId("")
	} else {
		if &route.Protocols != nil {
			d.Set("protocols", gokong.StringValueSlice(route.Protocols))
		}

		if &route.Methods != nil {
			d.Set("methods", gokong.StringValueSlice(route.Methods))
		}

		if &route.Hosts != nil {
			d.Set("hosts", gokong.StringValueSlice(route.Hosts))
		}

		if &route.Paths != nil {
			d.Set("paths", gokong.StringValueSlice(route.Paths))
		}

		if &route.StripPath != nil {
			d.Set("strip_path", route.StripPath)
		}

		if &route.PreserveHost != nil {
			d.Set("preserve_host", route.PreserveHost)
		}

		if route.Service != nil {
			d.Set("service_id", route.Service.Id)
		}
	}

	return nil
}

func resourceKongRouteDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*config).adminClient.Routes().DeleteRoute(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong route: %v", err)
	}

	return nil
}

func createKongRouteRequestFromResourceData(d *schema.ResourceData) *gokong.RouteRequest {
	service := gokong.RouteServiceObject{
		Id: readStringFromResource(d, "service_id"),
	}

	return &gokong.RouteRequest{
		Protocols:    readStringArrayPtrFromResource(d, "protocols"),
		Methods:      readStringArrayPtrFromResource(d, "methods"),
		Hosts:        readStringArrayPtrFromResource(d, "hosts"),
		Paths:        readStringArrayPtrFromResource(d, "paths"),
		StripPath:    readBoolPtrFromResource(d, "strip_path"),
		PreserveHost: readBoolPtrFromResource(d, "preserve_host"),
		Service:      &service,
	}
}
