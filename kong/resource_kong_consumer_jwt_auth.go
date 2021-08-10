package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
	"strings"
)

func resourceKongConsumerJWTAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerJWTAuthCreate,
		ReadContext:   resourceKongConsumerJWTAuthRead,
		DeleteContext: resourceKongConsumerJWTAuthDelete,
		UpdateContext: resourceKongConsumerJWTAuthUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"rsa_public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"secret": {
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

func resourceKongConsumerJWTAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	JWTAuthRequest := &kong.JWTAuth{
		Algorithm:    kong.String(d.Get("algorithm").(string)),
		Key:          kong.String(d.Get("key").(string)),
		RSAPublicKey: kong.String(d.Get("rsa_public_key").(string)),
		Secret:       kong.String(d.Get("secret").(string)),
		Tags:         readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.JWTAuths
	JWTAuth, err := client.Create(ctx, consumerId, JWTAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong JWTAuth: %v error: %v", JWTAuthRequest, err))
	}

	d.SetId(buildConsumerPairID(*JWTAuth.ID, *consumerId))

	return resourceKongConsumerJWTAuthRead(ctx, d, meta)
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

func resourceKongConsumerJWTAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	id, err := splitConsumerID(d.Id())

	JWTAuthRequest := &kong.JWTAuth{
		ID:           kong.String(id.ID),
		Algorithm:    kong.String(d.Get("algorithm").(string)),
		Key:          kong.String(d.Get("key").(string)),
		RSAPublicKey: kong.String(d.Get("rsa_public_key").(string)),
		Secret:       kong.String(d.Get("secret").(string)),
		Tags:         readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.JWTAuths
	_, err = client.Update(ctx, consumerId, JWTAuthRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong JWTAuth: %s", err))
	}

	return resourceKongConsumerJWTAuthRead(ctx, d, meta)
}

func resourceKongConsumerJWTAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.JWTAuths
	JWTAuth, err := client.Get(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong JWTAuth with id: %s error: %v", id, err))
	}

	if JWTAuth == nil {
		d.SetId("")
	} else {
		err = d.Set("consumer_id", JWTAuth.Consumer.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("key", JWTAuth.Key)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("secret", JWTAuth.Secret)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("rsa_public_key", JWTAuth.RSAPublicKey)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("algorithm", JWTAuth.Algorithm)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", JWTAuth.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerJWTAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := meta.(*config).adminClient.JWTAuths
	err = client.Delete(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong JWTAuth: %v", err))
	}

	return diags
}
