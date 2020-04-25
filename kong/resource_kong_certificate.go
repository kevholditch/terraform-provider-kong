package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hbagdi/go-kong/kong"
)

func resourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongCertificateCreate,
		Read:   resourceKongCertificateRead,
		Delete: resourceKongCertificateDelete,
		Update: resourceKongCertificateUpdate,
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
		},
	}
}

func resourceKongCertificateCreate(d *schema.ResourceData, meta interface{}) error {

	certificateRequest := createKongCertificateFromResourceData(d)

	client := meta.(*config).adminClient.Certificates

	certificate, err := client.Create(context.Background(), certificateRequest)
	// certificate, err := meta.(*config).adminClient.Certificates().Create(certificateRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong certificate: %v error: %v", certificateRequest, err)
	}

	d.SetId(*certificate.ID)

	return resourceKongCertificateRead(d, meta)
}

func resourceKongCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	certificateRequest := createKongCertificateFromResourceData(d)

	client := meta.(*config).adminClient.Certificates

	_, err := client.Update(context.Background(), certificateRequest)

	if err != nil {
		return fmt.Errorf("error updating kong certificate: %s", err)
	}

	return resourceKongCertificateRead(d, meta)
}

func resourceKongCertificateRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Certificates

	certificate, err := client.Get(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not find kong certificate: %v", err)
	}

	if certificate == nil {
		d.SetId("")
	} else {
		if certificate.Cert != nil {
			d.Set("certificate", certificate.Cert)
		}

		if certificate.Key != nil {
			d.Set("private_key", certificate.Key)
		}
	}

	return nil
}

func resourceKongCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Certificates

	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong certificate: %v", err)
	}

	return nil
}

func createKongCertificateFromResourceData(d *schema.ResourceData) *kong.Certificate {

	certificate := &kong.Certificate{}

	certificate.Cert = readStringPtrFromResource(d, "certificate")
	certificate.Key = readStringPtrFromResource(d, "private_key")

	return certificate
}
