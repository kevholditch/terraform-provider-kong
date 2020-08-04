package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceKongCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKongCertificateRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKongCertificateRead(d *schema.ResourceData, meta interface{}) error {

	var filterId string

	if v, _ := d.GetOk("filter"); v != nil {
		filterSet := v.(*schema.Set).List()
		if len(filterSet) == 1 {
			filterMap := filterSet[0].(map[string]interface{})
			filterId = filterMap["id"].(string)
		}
	}

	result, err := meta.(*config).adminClient.Certificates().GetById(filterId)

	if err != nil {
		return fmt.Errorf("could not find certificate, error: %v", err)
	}

	if result == nil {
		return fmt.Errorf("could not find certificate by id: %v", filterId)
	}

	d.SetId(*result.Id)

	if result.Cert != nil {
		d.Set("certificate", result.Cert)
	}

	if result.Key != nil {
		d.Set("private_key", result.Key)
	}

	return nil
}
