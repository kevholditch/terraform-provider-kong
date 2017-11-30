package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongCertificateCreate,
		Read:   resourceKongCertificateRead,
		Delete: resourceKongCertificateDelete,
		Update: resourceKongCertificateUpdate,

		Schema: map[string]*schema.Schema{
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceKongCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	certificateRequest := createKongCertificateRequestFromResourceData(d)

	consumer, err := client.Certificates().Create(certificateRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong certificate: %v error: %v", certificateRequest, err)
	}

	d.SetId(consumer.Id)

	return resourceKongCertificateRead(d, meta)
}

func resourceKongCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	client := meta.(*gokong.KongAdminClient)

	certificateRequest := createKongCertificateRequestFromResourceData(d)

	id := d.Id()

	_, err := client.Certificates().UpdateById(id, certificateRequest)

	if err != nil {
		return fmt.Errorf("error updating kong certificate: %s", err)
	}

	return resourceKongCertificateRead(d, meta)
}

func resourceKongCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gokong.KongAdminClient)

	id := d.Id()

	certificate, err := client.Certificates().GetById(id)

	if err != nil {
		return fmt.Errorf("could not find kong certificate: %v", err)
	}

	d.Set("certificate", certificate.Cert)
	d.Set("private_key", certificate.Key)

	return nil
}

func resourceKongCertificateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*gokong.KongAdminClient)

	id := d.Id()

	err := client.Consumers().DeleteById(id)

	if err != nil {
		return fmt.Errorf("could not delete kong certificate: %v", err)
	}

	return nil
}

func createKongCertificateRequestFromResourceData(d *schema.ResourceData) *gokong.CertificateRequest {

	certificateRequest := &gokong.CertificateRequest{}

	certificateRequest.Cert = readStringFromResource(d, "certificate")
	certificateRequest.Key = readStringFromResource(d, "private_key")

	return certificateRequest
}
