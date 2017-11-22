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
				ForceNew: true,
			},
			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				ForceNew: false,
			},
			"uris": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				ForceNew: false,
			},
			"methods": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				ForceNew: false,
			},
			"upstream_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"strip_uri": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
			},
			"preserve_host": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
			},
			"retries": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				ForceNew: false,
			},
			"upstream_connect_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				ForceNew: false,
			},
			"upstream_send_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				ForceNew: false,
			},
			"upstream_read_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				ForceNew: false,
			},
			"https_only": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
			},
			"http_if_terminated": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
			},
		},
	}
}

func resourceKongApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	apiRequest := createKongApiRequestFromResourceData(d)

	api, err := client.Apis().Create(apiRequest)

	if err != nil {
		return fmt.Errorf("Failed to create kong api: %v error: %v", apiRequest, err)
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
		return fmt.Errorf("Error updating kong api: %s", err)
	}
	
	return resourceKongApiRead(d, meta)
}

func resourceKongApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	id := d.Id()

	api, err := client.Apis().GetById(id)

	if err != nil {
		fmt.Errorf("Could not find kong api: %v", err)
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
		fmt.Errorf("Could not delete kong api: %v", err)
	}

	return nil
}

func createKongApiRequestFromResourceData(d *schema.ResourceData) *gokong.ApiRequest {

	apiRequest := &gokong.ApiRequest{}

	if attr, ok := d.GetOk("name"); ok {
		apiRequest.Name = attr.(string)
	}

	if attr, ok := d.GetOk("hosts"); ok {
		apiRequest.Hosts = attr.([]string)
	}

	if attr, ok := d.GetOk("uris"); ok {
		apiRequest.Uris = attr.([]string)
	}

	if attr, ok := d.GetOk("methods"); ok {
		apiRequest.Methods = attr.([]string)
	}

	if attr, ok := d.GetOk("upstream_url"); ok {
		apiRequest.UpstreamUrl = attr.(string)
	}

	if attr, ok := d.GetOk("strip_uri"); ok {
		apiRequest.StripUri = attr.(bool)
	}

	if attr, ok := d.GetOk("preserve_host"); ok {
		apiRequest.PreserveHost = attr.(bool)
	}

	if attr, ok := d.GetOk("retries"); ok {
		apiRequest.Retries = attr.(int)
	}

	if attr, ok := d.GetOk("upstream_connect_timeout"); ok {
		apiRequest.UpstreamConnectTimeout = attr.(int)
	}

	if attr, ok := d.GetOk("upstream_send_timeout"); ok {
		apiRequest.UpstreamSendTimeout = attr.(int)
	}

	if attr, ok := d.GetOk("upstream_read_timeout"); ok {
		apiRequest.UpstreamReadTimeout = attr.(int)
	}

	if attr, ok := d.GetOk("https_only"); ok {
		apiRequest.HttpsOnly = attr.(bool)
	}

	if attr, ok := d.GetOk("http_if_terminated"); ok {
		apiRequest.HttpIfTerminated = attr.(bool)
	}

	return apiRequest
}
