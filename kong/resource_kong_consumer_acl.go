package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongConsumerACLCreate,
		ReadContext:   resourceKongConsumerACLRead,
		DeleteContext: resourceKongConsumerACLDelete,
		UpdateContext: resourceKongConsumerACLUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"group": {
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

func resourceKongConsumerACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ACLGroupRequest := &kong.ACLGroup{
		Group: kong.String(d.Get("group").(string)),
		Tags:  readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.ACLs
	aclGroup, err := client.Create(ctx, consumerId, ACLGroupRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong ACL Group: %v error: %v", ACLGroupRequest, err))
	}

	d.SetId(buildConsumerPairID(*aclGroup.ID, *consumerId))

	return resourceKongConsumerACLRead(ctx, d, meta)
}

func resourceKongConsumerACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id, err := splitConsumerID(d.Id())

	ACLGroupRequest := &kong.ACLGroup{
		ID:    kong.String(id.ID),
		Group: kong.String(d.Get("group").(string)),
		Tags:  readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.ACLs
	_, err = client.Update(ctx, consumerId, ACLGroupRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong ACL Group: %s", err))
	}

	return resourceKongConsumerACLRead(ctx, d, meta)
}

func resourceKongConsumerACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.ACLs
	ACLGroup, err := client.Get(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong ACLGroup with id: %s error: %v", id, err))
	}

	if ACLGroup == nil {
		d.SetId("")
	} else {
		err := d.Set("consumer_id", ACLGroup.Consumer.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("group", ACLGroup.Group)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", ACLGroup.Tags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceKongConsumerACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := meta.(*config).adminClient.ACLs
	err = client.Delete(ctx, kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong ACL Group: %v", err))
	}

	return diags
}
