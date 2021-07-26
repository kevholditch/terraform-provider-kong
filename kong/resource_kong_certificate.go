package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongCertificateCreate,
		ReadContext:   resourceKongCertificateRead,
		DeleteContext: resourceKongCertificateDelete,
		UpdateContext: resourceKongCertificateUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"certificate": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  false,
				Sensitive: false,
			},
			"private_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  false,
				Sensitive: true,
			},
			"snis": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKongCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	certificateRequest := &kong.Certificate{
		Cert: kong.String(d.Get("certificate").(string)),
		Key:  kong.String(d.Get("private_key").(string)),
		SNIs: readStringArrayPtrFromResource(d, "snis"),
	}

	client := meta.(*config).adminClient.Certificates

	certificate, err := client.Create(ctx, certificateRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong certificate: %v error: %v", certificateRequest, err))
	}

	d.SetId(*certificate.ID)

	return resourceKongCertificateRead(ctx, d, meta)
}

func resourceKongCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	certificateRequest := &kong.Certificate{
		ID:   kong.String(d.Id()),
		Cert: kong.String(d.Get("certificate").(string)),
		Key:  kong.String(d.Get("private_key").(string)),
		SNIs: readStringArrayPtrFromResource(d, "snis"),
	}

	client := meta.(*config).adminClient.Certificates

	_, err := client.Update(ctx, certificateRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong certificate: %s", err))
	}

	return resourceKongCertificateRead(ctx, d, meta)
}

func resourceKongCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Certificates

	certificate, err := client.Get(ctx, kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong certificate: %v", err))
	}

	if certificate == nil {
		d.SetId("")
	} else {
		if certificate.Cert != nil {
			d.Set("certificate", &certificate.Cert)
		}

		if certificate.Key != nil {
			d.Set("private_key", &certificate.Key)
		}

		if certificate.SNIs != nil {
			d.Set("snis", StringValueSlice(certificate.SNIs))
		}
	}

	return diags
}

func resourceKongCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Certificates

	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong certificate: %v", err))
	}

	return diags
}
