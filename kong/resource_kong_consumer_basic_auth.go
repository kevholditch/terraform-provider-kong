package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerBasicAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerBasicAuthCreate,
		ReadContext:   resourceKongConsumerBasicAuthRead,
		DeleteContext: resourceKongConsumerBasicAuthDelete,
		UpdateContext: resourceKongConsumerBasicAuthUpdate,
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKongConsumerBasicAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	BasicAuthRequest := &kong.BasicAuth{
		Username: kong.String(d.Get("username").(string)),
		Password: kong.String(d.Get("password").(string)),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.BasicAuths
	basicAuth, err := client.Create(ctx, consumerId, BasicAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong basic auth: %v error: %v", BasicAuthRequest, err))
	}

	d.SetId(buildConsumerPairID(*basicAuth.ID, *consumerId))

	return resourceKongConsumerBasicAuthRead(ctx, d, meta)
}

func resourceKongConsumerBasicAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := splitConsumerID(d.Id())

	BasicAuthRequest := &kong.BasicAuth{
		ID:       kong.String(id.ID),
		Username: kong.String(d.Get("username").(string)),
		Password: kong.String(d.Get("password").(string)),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.BasicAuths
	_, err = client.Update(ctx, consumerId, BasicAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong basic auth: %s", err))
	}

	return resourceKongConsumerBasicAuthRead(ctx, d, meta)
}

func resourceKongConsumerBasicAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.BasicAuths
	basicAuth, err := client.Get(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong basic auth with id: %s error: %v", id, err))
	}

	if basicAuth == nil {
		d.SetId("")
	} else {
		err = d.Set("consumer_id", basicAuth.Consumer.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("username", basicAuth.Username)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", basicAuth.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerBasicAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := meta.(*config).adminClient.BasicAuths
	err = client.Delete(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong basic auth: %v", err))
	}

	return diags
}
