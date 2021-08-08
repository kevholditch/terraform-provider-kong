package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongServiceCreate,
		ReadContext:   resourceKongServiceRead,
		DeleteContext: resourceKongServiceDelete,
		UpdateContext: resourceKongServiceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  80,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"retries": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  5,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"write_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  60000,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tls_verify": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"tls_verify_depth": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"ca_certificate_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKongServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	serviceRequest := createKongServiceRequestFromResourceData(d)

	client := meta.(*config).adminClient.Services
	service, err := client.Create(ctx, serviceRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong service: %v error: %v", serviceRequest, err))
	}

	d.SetId(*service.ID)

	return resourceKongServiceRead(ctx, d, meta)
}

func resourceKongServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	serviceRequest := createKongServiceRequestFromResourceData(d)

	client := meta.(*config).adminClient.Services
	_, err := client.Update(ctx, serviceRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong service: %s", err))
	}

	return resourceKongServiceRead(ctx, d, meta)
}

func resourceKongServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*config).adminClient.Services
	service, err := client.Get(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong service: %v", err))
	}

	if service == nil {
		d.SetId("")
	} else {
		if service.Name != nil {
			err := d.Set("name", service.Name)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.Protocol != nil {
			err := d.Set("protocol", service.Protocol)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.Host != nil {
			err := d.Set("host", service.Host)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.Port != nil {
			err := d.Set("port", service.Port)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.Path != nil {
			err := d.Set("path", service.Path)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.Retries != nil {
			err := d.Set("retries", service.Retries)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.ConnectTimeout != nil {
			err := d.Set("connect_timeout", service.ConnectTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.WriteTimeout != nil {
			err := d.Set("write_timeout", service.WriteTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.ReadTimeout != nil {
			err := d.Set("read_timeout", service.ReadTimeout)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		err = d.Set("tags", service.Tags)
		if err != nil {
			return diag.FromErr(err)
		}

		if service.TLSVerify != nil {
			err = d.Set("tls_verify", service.TLSVerify)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if service.TLSVerifyDepth != nil {
			err = d.Set("tls_verify_depth", service.TLSVerifyDepth)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if service.ClientCertificate != nil {
			err = d.Set("client_certificate_id", service.ClientCertificate.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if service.CACertificates != nil {
			err = d.Set("ca_certificate_ids", service.CACertificates)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceKongServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Services
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong service: %v", err))
	}

	return diags
}

func createKongServiceRequestFromResourceData(d *schema.ResourceData) *kong.Service {
	service := &kong.Service{
		Name:           readStringPtrFromResource(d, "name"),
		Protocol:       readStringPtrFromResource(d, "protocol"),
		Host:           readStringPtrFromResource(d, "host"),
		Port:           readIntPtrFromResource(d, "port"),
		Path:           readStringPtrFromResource(d, "path"),
		Retries:        readIntWithZeroPtrFromResource(d, "retries"),
		ConnectTimeout: readIntPtrFromResource(d, "connect_timeout"),
		WriteTimeout:   readIntPtrFromResource(d, "write_timeout"),
		ReadTimeout:    readIntPtrFromResource(d, "read_timeout"),
		Tags:           readStringArrayPtrFromResource(d, "tags"),
		TLSVerify:      readBoolPtrFromResource(d, "tls_verify"),
		TLSVerifyDepth: readIntPtrFromResource(d, "tls_verify_depth"),
		CACertificates: readStringArrayPtrFromResource(d, "ca_certificate_ids"),
	}

	clientCertificateID := readIdPtrFromResource(d, "client_certificate_id")
	if clientCertificateID != nil {
		service.ClientCertificate = &kong.Certificate{
			ID: clientCertificateID,
		}
	}

	if d.Id() != "" {
		service.ID = kong.String(d.Id())
	}
	return service
}
