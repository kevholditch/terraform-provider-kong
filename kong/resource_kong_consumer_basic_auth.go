package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerBasicAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerBasicAuthCreate,
		Read:   resourceKongConsumerBasicAuthRead,
		Delete: resourceKongConsumerBasicAuthDelete,
		Update: resourceKongConsumerBasicAuthUpdate,
		Schema: map[string]*schema.Schema{
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceKongConsumerBasicAuthCreate(d *schema.ResourceData, meta interface{}) error {
	BasicAuthRequest := &kong.BasicAuth{
		Username: kong.String(d.Get("username").(string)),
		Password: kong.String(d.Get("password").(string)),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.BasicAuths
	basicAuth, err := client.Create(context.Background(), consumerId, BasicAuthRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong basic auth: %v error: %v", BasicAuthRequest, err)
	}

	d.SetId(buildConsumerPairID(*basicAuth.ID, *consumerId))

	return resourceKongConsumerBasicAuthRead(d, meta)
}

func resourceKongConsumerBasicAuthUpdate(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())

	BasicAuthRequest := &kong.BasicAuth{
		ID:       kong.String(id.ID),
		Username: kong.String(d.Get("username").(string)),
		Password: kong.String(d.Get("password").(string)),
		Tags:     readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.BasicAuths
	_, err = client.Update(context.Background(), consumerId, BasicAuthRequest)

	if err != nil {
		return fmt.Errorf("error updating kong basic auth: %s", err)
	}

	return resourceKongConsumerBasicAuthRead(d, meta)
}

func resourceKongConsumerBasicAuthRead(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.BasicAuths
	basicAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return fmt.Errorf("could not find kong ACLGroup with id: %s error: %v", id, err)
	}

	if basicAuth == nil {
		d.SetId("")
	} else {
		d.Set("consumer_id", basicAuth.Consumer.ID)
		d.Set("username", basicAuth.Username)
		d.Set("tags", basicAuth.Tags)
	}

	return nil
}

func resourceKongConsumerBasicAuthDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}
	client := meta.(*config).adminClient.BasicAuths
	err = client.Delete(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return fmt.Errorf("could not delete kong basic auth: %v", err)
	}

	return nil
}
