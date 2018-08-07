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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

	apiRequest := createKongApiRequestFromResourceData(d)

	api, err := meta.(*gokong.KongAdminClient).Apis().Create(apiRequest)

	if err != nil || api == nil {
		return fmt.Errorf("failed to create kong api: %v error: %v", apiRequest, err)
	}

	d.SetId(*api.Id)

	return resourceKongApiRead(d, meta)
}

func resourceKongApiUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	apiRequest := createKongApiRequestFromResourceData(d)

	_, err := meta.(*gokong.KongAdminClient).Apis().UpdateById(d.Id(), apiRequest)

	if err != nil {
		return fmt.Errorf("error updating kong api: %s", err)
	}

	return resourceKongApiRead(d, meta)
}

func resourceKongApiRead(d *schema.ResourceData, meta interface{}) error {

	api, err := meta.(*gokong.KongAdminClient).Apis().GetById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong api: %v", err)
	}

	if api == nil {
		d.SetId("")
	} else {
		if api.Name != nil {
			d.Set("name", api.Name)
		}

		if api.Hosts != nil {
			d.Set("hosts", gokong.StringValueSlice(api.Hosts))
		}

		if api.Uris != nil {
			d.Set("uris", gokong.StringValueSlice(api.Uris))
		}

		if api.Methods != nil {
			d.Set("methods", gokong.StringValueSlice(api.Methods))
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
	}

	return nil
}

func resourceKongApiDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*gokong.KongAdminClient).Apis().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong api: %v", err)
	}

	return nil
}

func createKongApiRequestFromResourceData(d *schema.ResourceData) *gokong.ApiRequest {

	return &gokong.ApiRequest{
		Name:                   readStringPtrFromResource(d, "name"),
		Hosts:                  readStringArrayPtrFromResource(d, "hosts"),
		Uris:                   readStringArrayPtrFromResource(d, "uris"),
		Methods:                readStringArrayPtrFromResource(d, "methods"),
		UpstreamUrl:            readStringPtrFromResource(d, "upstream_url"),
		StripUri:               readBoolPtrFromResource(d, "strip_uri"),
		PreserveHost:           readBoolPtrFromResource(d, "preserve_host"),
		Retries:                readIntPtrFromResource(d, "retries"),
		UpstreamConnectTimeout: readIntPtrFromResource(d, "upstream_connect_timeout"),
		UpstreamSendTimeout:    readIntPtrFromResource(d, "upstream_send_timeout"),
		UpstreamReadTimeout:    readIntPtrFromResource(d, "upstream_read_timeout"),
		HttpsOnly:              readBoolPtrFromResource(d, "https_only"),
		HttpIfTerminated:       readBoolPtrFromResource(d, "http_if_terminated"),
	}
}
