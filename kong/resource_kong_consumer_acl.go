package kong

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongConsumerACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerACLCreate,
		Read:   resourceKongConsumerACLRead,
		Delete: resourceKongConsumerACLDelete,
		Update: resourceKongConsumerACLUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"group": &schema.Schema{
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

func resourceKongConsumerACLCreate(d *schema.ResourceData, meta interface{}) error {
	ACLGroupRequest := &kong.ACLGroup{
		Group: kong.String(d.Get("group").(string)),
		Tags:  readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.ACLs
	aclGroup, err := client.Create(context.Background(), consumerId, ACLGroupRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong ACL Group: %v error: %v", ACLGroupRequest, err)
	}

	d.SetId(buildJWTID(*aclGroup.ID, *consumerId))

	return resourceKongConsumerACLRead(d, meta)
}

func resourceKongConsumerACLUpdate(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())

	ACLGroupRequest := &kong.ACLGroup{
		ID:    kong.String(id.ID),
		Group: kong.String(d.Get("group").(string)),
		Tags:  readStringArrayPtrFromResource(d, "tags"),
	}

	consumerId := kong.String(d.Get("consumer_id").(string))

	client := meta.(*config).adminClient.ACLs
	_, err = client.Update(context.Background(), consumerId, ACLGroupRequest)

	if err != nil {
		return fmt.Errorf("error updating kong ACL Group: %s", err)
	}

	return resourceKongConsumerACLRead(d, meta)
}

func resourceKongConsumerACLRead(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.ACLs
	ACLGroup, err := client.Get(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if kong.IsNotFoundErr(err) {
		d.SetId("")
	} else if err != nil {
		return fmt.Errorf("could not find kong ACLGroup with id: %s error: %v", id, err)
	}

	if ACLGroup == nil {
		d.SetId("")
	} else {
		d.Set("consumer_id", ACLGroup.Consumer.ID)
		d.Set("group", ACLGroup.Group)
		d.Set("tags", ACLGroup.Tags)
	}

	return nil
}

func resourceKongConsumerACLDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := splitConsumerID(d.Id())
	if err != nil {
		return err
	}
	client := meta.(*config).adminClient.ACLs
	err = client.Delete(context.Background(), kong.String(id.ConsumerID), kong.String(id.ID))

	if err != nil {
		return fmt.Errorf("could not delete kong ACL Group: %v", err)
	}

	return nil
}
