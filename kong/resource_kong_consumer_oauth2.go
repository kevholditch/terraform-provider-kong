package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerOAuth2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerOAuth2Create,
		ReadContext:   resourceKongConsumerOAuth2Read,
		DeleteContext: resourceKongConsumerOAuth2Delete,
		UpdateContext: resourceKongConsumerOAuth2Update,
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"hash_secret": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceKongConsumerOAuth2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	OAuth2CredentialRequest := &kong.Oauth2Credential{
		Name:         readStringPtrFromResource(d, "name"),
		ClientID:     readStringPtrFromResource(d, "client_id"),
		ClientSecret: readStringPtrFromResource(d, "client_secret"),
		HashSecret:   readBoolPtrFromResource(d, "hash_secret"),
		RedirectURIs: readStringArrayPtrFromResource(d, "redirect_uris"),
		Tags:         readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.Oauth2Credentials
	oAuth2Credentials, err := client.Create(ctx, consumerId, OAuth2CredentialRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create oauth2 credentials: %v error: %v", OAuth2CredentialRequest, err))
	}

	d.SetId(buildConsumerPairID(*oAuth2Credentials.ID, *consumerId))

	return resourceKongConsumerOAuth2Read(ctx, d, meta)
}

func resourceKongConsumerOAuth2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, _ := splitConsumerID(d.Id())

	OAuth2CredentialRequest := &kong.Oauth2Credential{
		ID:           kong.String(id.ID),
		Name:         readStringPtrFromResource(d, "name"),
		ClientID:     readStringPtrFromResource(d, "client_id"),
		ClientSecret: readStringPtrFromResource(d, "client_secret"),
		HashSecret:   readBoolPtrFromResource(d, "hash_secret"),
		RedirectURIs: readStringArrayPtrFromResource(d, "redirect_uris"),
		Tags:         readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.Oauth2Credentials
	_, err := client.Update(ctx, consumerId, OAuth2CredentialRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong oauth2 credentials: %s", err))
	}

	return resourceKongConsumerOAuth2Read(ctx, d, meta)
}

func resourceKongConsumerOAuth2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.Oauth2Credentials
	oAuth2Credentials, err := client.Get(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong oauth2 credentials with id: %s error: %v", id, err))
	}

	if oAuth2Credentials == nil {
		d.SetId("")
	} else {
		err = d.Set("consumer_id", oAuth2Credentials.Consumer.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("name", oAuth2Credentials.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("client_id", oAuth2Credentials.ClientID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("client_secret", oAuth2Credentials.ClientSecret)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("hash_secret", oAuth2Credentials.HashSecret)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("redirect_uris", oAuth2Credentials.RedirectURIs)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", oAuth2Credentials.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerOAuth2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := meta.(*config).adminClient.Oauth2Credentials
	err = client.Delete(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong oauth2 credentials: %v", err))
	}

	return diags
}
