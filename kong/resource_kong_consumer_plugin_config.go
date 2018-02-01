package kong

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func resourceKongConsumerPluginConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerPluginConfigCreate,
		Read:   resourceKongConsumerPluginConfigRead,
		Delete: resourceKongConsumerPluginConfigDelete,
		Update: resourceKongConsumerPluginConfigUpdate,

		Schema: map[string]*schema.Schema{
			"plugin_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"config_json": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceKongConsumerPluginConfigCreate(d *schema.ResourceData, meta interface{}) error {
	plugin, err := createKongPluginConfig(d, meta)
	if err != nil {
		return fmt.Errorf("failed to configure kong plugin: %v error: %v", plugin, err)
	}

	d.SetId(plugin.Id)

	return resourceKongConsumerPluginConfigRead(d, meta)
}

func resourceKongConsumerPluginConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	_, err := createKongPluginConfig(d, meta)

	if err != nil {
		return fmt.Errorf("error updating kong plugin: %s", err)
	}

	return resourceKongConsumerPluginConfigRead(d, meta)
}

func resourceKongConsumerPluginConfigRead(d *schema.ResourceData, meta interface{}) error {
	pluginName := readStringFromResource(d, "plugin_name")
	consumerID := readStringFromResource(d, "consumer_id")
	plugin, err := meta.(*gokong.KongAdminClient).Consumers().GetPluginConfig(consumerID, pluginName, d.Id())

	if err != nil {
		return fmt.Errorf("could not find kong plugin: %v", err)
	}

	if plugin == nil {
		d.SetId("")
	} else {
		d.Set("id", plugin.Id)
		d.Set("config", plugin.Body)
	}

	return nil
}

func resourceKongConsumerPluginConfigDelete(d *schema.ResourceData, meta interface{}) error {
	pluginName := readStringFromResource(d, "plugin_name")
	consumerID := readStringFromResource(d, "consumer_id")
	err := meta.(*gokong.KongAdminClient).Consumers().DeletePluginConfig(consumerID, pluginName, d.Id())
	if err != nil {
		return fmt.Errorf("could not delete kong plugin: %v", err)
	}

	return nil
}

func createKongPluginConfig(d *schema.ResourceData, meta interface{}) (*gokong.ConsumerPluginConfig, error) {
	pluginName := readStringFromResource(d, "plugin_name")
	pluginConfig := readStringFromResource(d, "config_json")
	consumerID := readStringFromResource(d, "consumer_id")
	plugin, err := meta.(*gokong.KongAdminClient).Consumers().CreatePluginConfig(consumerID, pluginName, pluginConfig)
	return plugin, err

}
