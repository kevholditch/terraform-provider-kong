package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerCreate,
		ReadContext:   resourceKongConsumerRead,
		DeleteContext: resourceKongConsumerDelete,
		UpdateContext: resourceKongConsumerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"custom_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceKongConsumerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	consumerRequest := &kong.Consumer{
		Username: readStringPtrFromResource(d, "username"),
		CustomID: readStringPtrFromResource(d, "custom_id"),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	client := meta.(*config).adminClient.Consumers
	consumer, err := client.Create(ctx, consumerRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong consumer: %v error: %v", consumerRequest, err))
	}

	d.SetId(*consumer.ID)

	return resourceKongConsumerRead(ctx, d, meta)
}

func resourceKongConsumerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	consumerRequest := &kong.Consumer{
		ID:       kong.String(d.Id()),
		Username: kong.String(d.Get("username").(string)),
		CustomID: kong.String(d.Get("custom_id").(string)),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	client := meta.(*config).adminClient.Consumers
	_, err := client.Update(ctx, consumerRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong consumer: %s", err))
	}

	return resourceKongConsumerRead(ctx, d, meta)
}

func resourceKongConsumerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	id := d.Id()

	client := meta.(*config).adminClient.Consumers
	consumer, err := client.Get(ctx, kong.String(id))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong consumer with id: %s error: %v", id, err))
	}

	if consumer == nil {
		d.SetId("")
	} else {
		err := d.Set("username", consumer.Username)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("custom_id", consumer.CustomID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", consumer.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Consumers
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong consumer: %v", err))
	}

	return diags
}
