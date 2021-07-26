package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
	"strings"
)

func resourceKongJWTAuth() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongJWTAuthCreate,
		Read:   resourceKongJWTAuthRead,
		Delete: resourceKongJWTAuthDelete,
		Update: resourceKongJWTAuthUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"rsa_public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func resourceKongJWTAuthCreate(d *schema.ResourceData, meta interface{}) error {

	JWTAuthRequest := &kong.JWTAuth{
		Algorithm:    kong.String(d.Get("algorithm").(string)),
		Key:          kong.String(d.Get("key").(string)),
		RSAPublicKey: kong.String(d.Get("rsa_public_key").(string)),
		Secret:       kong.String(d.Get("secret").(string)),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.JWTAuths
	JWTAuth, err := client.Create(context.Background(), consumerId, JWTAuthRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong JWTAuth: %v error: %v", JWTAuthRequest, err)
	}

	d.SetId(buildConsumerPairID(*JWTAuth.ID, *consumerId))

	return resourceKongJWTAuthRead(d, meta)
}

func buildConsumerPairID(ID, consumerID string) string {
	return ID + "|" + consumerID
}

type ConsumerIDPair struct {
	ID         string
	ConsumerID string
}

func splitConsumerID(value string) (*ConsumerIDPair, error) {
	v := strings.Split(value, "|")
	if len(v) != 2 {
		return nil, fmt.Errorf("expecting there to be exactly 2 strings in ID but found %d", len(v))
	}
	return &ConsumerIDPair{ID: v[0], ConsumerID: v[1]}, nil
}

func resourceKongJWTAuthUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	id, err := splitConsumerID(d.Id())

	JWTAuthRequest := &kong.JWTAuth{
		ID:           kong.String(id.ID),
		Algorithm:    kong.String(d.Get("algorithm").(string)),
		Key:          kong.String(d.Get("key").(string)),
		RSAPublicKey: kong.String(d.Get("rsa_public_key").(string)),
		Secret:       kong.String(d.Get("secret").(string)),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.JWTAuths
	_, err = client.Update(context.Background(), consumerId, JWTAuthRequest)

	if err != nil {
		return fmt.Errorf("error updating kong JWTAuth: %s", err)
	}

	return resourceKongJWTAuthRead(d, meta)
}

func resourceKongJWTAuthRead(d *schema.ResourceData, meta interface{}) error {

	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.JWTAuths
	JWTAuth, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return fmt.Errorf("could not find kong JWTAuth with id: %s error: %v", id, err)
	}

	if JWTAuth == nil {
		d.SetId("")
	} else {
		d.Set("consumer_id", JWTAuth.Consumer.ID)
		d.Set("key", JWTAuth.Key)
		d.Set("secret", JWTAuth.Secret)
		d.Set("rsa_public_key", JWTAuth.RSAPublicKey)
		d.Set("algorithm", JWTAuth.Algorithm)
	}

	return nil
}

func resourceKongJWTAuthDelete(d *schema.ResourceData, meta interface{}) error {

	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}
	client := meta.(*config).adminClient.JWTAuths
	err = client.Delete(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return fmt.Errorf("could not delete kong JWTAuth: %v", err)
	}

	return nil
}
