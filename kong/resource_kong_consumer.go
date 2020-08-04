package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
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

	consumerRequest := createKongConsumerRequestFromResourceData(d)

	consumer, err := meta.(*config).adminClient.Consumers().Create(consumerRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong consumer: %v error: %v", consumerRequest, err)
	}

	d.SetId(consumer.Id)

	return resourceKongConsumerRead(d, meta)
}

func resourceKongConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	consumerRequest := createKongConsumerRequestFromResourceData(d)

	_, err := meta.(*config).adminClient.Consumers().UpdateById(d.Id(), consumerRequest)

	if err != nil {
		return fmt.Errorf("error updating kong consumer: %s", err)
	}

	return resourceKongConsumerRead(d, meta)
}

func resourceKongConsumerRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	consumer, err := meta.(*config).adminClient.Consumers().GetById(id)

	if err != nil {
		return fmt.Errorf("could not find kong consumer with id: %s error: %v", id, err)
	}

	if consumer == nil {
		d.SetId("")
	} else {
		d.Set("username", consumer.Username)
		d.Set("custom_id", consumer.CustomId)
	}

	return nil
}

func resourceKongConsumerDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*config).adminClient.Consumers().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong consumer: %v", err)
	}

	return nil
}

func createKongConsumerRequestFromResourceData(d *schema.ResourceData) *gokong.ConsumerRequest {

	consumerRequest := &gokong.ConsumerRequest{}

	consumerRequest.Username = readStringFromResource(d, "username")
	consumerRequest.CustomId = readStringFromResource(d, "custom_id")

	return consumerRequest
}
