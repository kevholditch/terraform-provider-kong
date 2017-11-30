package kong

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongPluginCreate,
		Read:   resourceKongPluginRead,
		Delete: resourceKongPluginDelete,
		Update: resourceKongPluginUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"api_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				Default:  nil,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {

	pluginRequest := createKongPluginRequestFromResourceData(d)

	plugin, err := meta.(*gokong.KongAdminClient).Plugins().Create(pluginRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong plugin: %v error: %v", pluginRequest, err)
	}

	d.SetId(plugin.Id)

	return resourceKongPluginRead(d, meta)
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	pluginRequest := createKongPluginRequestFromResourceData(d)

	_, err := meta.(*gokong.KongAdminClient).Plugins().UpdateById(d.Id(), pluginRequest)

	if err != nil {
		return fmt.Errorf("error updating kong plugin: %s", err)
	}

	return resourceKongPluginRead(d, meta)
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {

	plugin, err := meta.(*gokong.KongAdminClient).Plugins().GetById(d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong plugin: %v", err)
	}

	d.Set("name", plugin.Name)

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*gokong.KongAdminClient).Plugins().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong plugin: %v", err)
	}

	return nil
}

func createKongPluginRequestFromResourceData(d *schema.ResourceData) *gokong.PluginRequest {

	pluginRequest := &gokong.PluginRequest{}

	pluginRequest.Name = readStringFromResource(d, "name")
	pluginRequest.ApiId = readStringFromResource(d, "api_id")
	pluginRequest.ConsumerId = readStringFromResource(d, "consumer_id")
	pluginRequest.Config = readMapFromResource(d, "config")

	return pluginRequest
}
