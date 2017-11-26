package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongApiCreate,
		Read:   resourceKongApiRead,
		Delete: resourceKongApiDelete,
		Update: resourceKongApiUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"uris": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"methods": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"upstream_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"strip_uri": &schema.Schema{
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
			"retries": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  5,
			},
			"upstream_connect_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"upstream_send_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"upstream_read_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"https_only": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"http_if_terminated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceKongApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	apiRequest := createKongApiRequestFromResourceData(d)

	api, err := client.Apis().Create(apiRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong api: %v error: %v", apiRequest, err)
	}

	d.SetId(api.Id)

	return resourceKongApiRead(d, meta)
}

func resourceKongApiUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	client := meta.(*gokong.KongAdminClient)

	apiRequest := createKongApiRequestFromResourceData(d)

	id := d.Id()

	_, err := client.Apis().UpdateById(id, apiRequest)

	if err != nil {
		return fmt.Errorf("error updating kong api: %s", err)
	}

	return resourceKongApiRead(d, meta)
}

func resourceKongApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	id := d.Id()

	api, err := client.Apis().GetById(id)

	if err != nil {
		return fmt.Errorf("could not find kong api: %v", err)
	}

	d.Set("name", api.Name)
	d.Set("hosts", api.Hosts)
	d.Set("uris", api.Uris)
	d.Set("methods", api.Methods)
	d.Set("upstream_url", api.UpstreamUrl)
	d.Set("strip_uri", api.StripUri)
	d.Set("preserve_host", api.PreserveHost)
	d.Set("retries", api.Retries)
	d.Set("upstream_connect_timeout", api.UpstreamConnectTimeout)
	d.Set("upstream_send_timeout", api.UpstreamSendTimeout)
	d.Set("upstream_read_timeout", api.UpstreamReadTimeout)
	d.Set("https_only", api.HttpsOnly)
	d.Set("http_if_terminated", api.HttpIfTerminated)

	return nil
}

func resourceKongApiDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*gokong.KongAdminClient)

	id := d.Id()

	err := client.Apis().DeleteById(id)

	if err != nil {
		return fmt.Errorf("could not delete kong api: %v", err)
	}

	return nil
}

func createKongApiRequestFromResourceData(d *schema.ResourceData) *gokong.ApiRequest {

	apiRequest := &gokong.ApiRequest{}

	apiRequest.Name = readStringFromResource(d, "name")
	apiRequest.Hosts = readArrayFromResource(d, "hosts")
	apiRequest.Uris = readArrayFromResource(d, "uris")
	apiRequest.Methods = readArrayFromResource(d, "methods")
	apiRequest.UpstreamUrl = readStringFromResource(d, "upstream_url")
	apiRequest.StripUri = readBoolFromResource(d, "strip_uri")
	apiRequest.PreserveHost = readBoolFromResource(d, "preserve_host")
	apiRequest.Retries = readIntFromResource(d, "retries")
	apiRequest.UpstreamConnectTimeout = readIntFromResource(d, "upstream_connect_timeout")
	apiRequest.UpstreamSendTimeout = readIntFromResource(d, "upstream_send_timeout")
	apiRequest.UpstreamReadTimeout = readIntFromResource(d, "upstream_read_timeout")
	apiRequest.HttpsOnly = readBoolFromResource(d, "https_only")
	apiRequest.HttpIfTerminated = readBoolFromResource(d, "http_if_terminated")

	return apiRequest
}
