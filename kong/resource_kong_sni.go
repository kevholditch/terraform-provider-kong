package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongSni() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongSniCreate,
		Read:   resourceKongSniRead,
		Delete: resourceKongSniDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"certificate_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKongSniCreate(d *schema.ResourceData, meta interface{}) error {

	sniRequest := createKongSniRequestFromResourceData(d)

	client := meta.(*config).adminClient.SNIs
	sni, err := client.Create(context.Background(), sniRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong sni: %v error: %v", sniRequest, err)
	}

	d.SetId(*sni.Name)

	return resourceKongSniRead(d, meta)
}

func resourceKongSniRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.SNIs
	sni, err := client.Get(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not find kong sni: %v", err)
	}

	if sni == nil {
		d.SetId("")
	} else {
		d.Set("name", sni.Name)
		d.Set("certificate_id", sni.Certificate.ID)
	}

	return nil
}

func resourceKongSniDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.SNIs
	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong sni: %v", err)
	}

	return nil
}

func createKongSniRequestFromResourceData(d *schema.ResourceData) *kong.SNI {

	sniRequest := &kong.SNI{}

	sniRequest.Name = readStringPtrFromResource(d, "name")
	sniRequest.Certificate = &kong.Certificate{
		ID: readIdPtrFromResource(d, "certificate_id"),
	}
	if d.Id() != "" {
		sniRequest.ID = kong.String(d.Id())
	}
	return sniRequest
}
