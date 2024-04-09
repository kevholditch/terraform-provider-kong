package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerKeyAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerKeyAuthCreate,
		ReadContext:   resourceKongConsumerKeyAuthRead,
		DeleteContext: resourceKongConsumerKeyAuthDelete,
		UpdateContext: resourceKongConsumerKeyAuthUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				ForceNew:  false,
				Sensitive: true,
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

func resourceKongConsumerKeyAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	KeyAuthRequest := &kong.KeyAuth{
		Key:  readStringPtrFromResource(d, "key"),
		Tags: readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.KeyAuths
	keyAuth, err := client.Create(ctx, consumerId, KeyAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong key auth: %v error: %v", KeyAuthRequest, err))
	}

	d.SetId(buildConsumerPairID(*keyAuth.ID, *consumerId))

	return resourceKongConsumerKeyAuthRead(ctx, d, meta)
}

func resourceKongConsumerKeyAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := splitConsumerID(d.Id())

	KeyAuthRequest := &kong.KeyAuth{
		ID:   kong.String(id.ID),
		Key:  readStringPtrFromResource(d, "key"),
		Tags: readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.KeyAuths
	_, err = client.Update(ctx, consumerId, KeyAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong key auth: %s", err))
	}

	return resourceKongConsumerKeyAuthRead(ctx, d, meta)
}

func resourceKongConsumerKeyAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.KeyAuths
	keyAuth, err := client.Get(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong key auth with id: %s error: %v", id, err))
	}

	if keyAuth == nil {
		d.SetId("")
	} else {
		err = d.Set("consumer_id", keyAuth.Consumer.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("key", keyAuth.Key)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", keyAuth.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerKeyAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := meta.(*config).adminClient.KeyAuths
	err = client.Delete(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong key auth: %v", err))
	}

	return diags
}
