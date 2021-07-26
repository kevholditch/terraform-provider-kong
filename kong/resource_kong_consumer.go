package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerCreate,
		Read:   resourceKongConsumerRead,
		Delete: resourceKongConsumerDelete,
		Update: resourceKongConsumerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"custom_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceKongConsumerCreate(d *schema.ResourceData, meta interface{}) error {

	consumerRequest := &kong.Consumer{
		Username: kong.String(d.Get("username").(string)),
		CustomID: kong.String(d.Get("custom_id").(string)),
	}

	client := meta.(*config).adminClient.Consumers
	consumer, err := client.Create(context.Background(), consumerRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong consumer: %v error: %v", consumerRequest, err)
	}

	d.SetId(*consumer.ID)

	return resourceKongConsumerRead(d, meta)
}

func resourceKongConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	consumerRequest := &kong.Consumer{
		ID:       kong.String(d.Id()),
		Username: kong.String(d.Get("username").(string)),
		CustomID: kong.String(d.Get("custom_id").(string)),
	}

	client := meta.(*config).adminClient.Consumers
	_, err := client.Update(context.Background(), consumerRequest)

	if err != nil {
		return fmt.Errorf("error updating kong consumer: %s", err)
	}

	return resourceKongConsumerRead(d, meta)
}

func resourceKongConsumerRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	client := meta.(*config).adminClient.Consumers
	consumer, err := client.Get(context.Background(), kong.String(id))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return fmt.Errorf("could not find kong consumer with id: %s error: %v", id, err)
	}

	if consumer == nil {
		d.SetId("")
	} else {
		d.Set("username", consumer.Username)
		d.Set("custom_id", consumer.CustomID)
	}

	return nil
}

func resourceKongConsumerDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Consumers
	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong consumer: %v", err)
	}

	return nil
}
