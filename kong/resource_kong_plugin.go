package kong

import (
	"encoding/json"
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

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"service_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"route_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"config_json": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "plugin configuration in JSON format, configuration must be a valid JSON object.",
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {

	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return err
	}

	plugin, err := meta.(*gokong.KongAdminClient).Plugins().Create(pluginRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong plugin: %v error: %v", pluginRequest, err)
	}

	d.SetId(plugin.Id)

	return resourceKongPluginRead(d, meta)
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return err
	}

	_, err = meta.(*gokong.KongAdminClient).Plugins().UpdateById(d.Id(), pluginRequest)

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

	if plugin == nil {
		d.SetId("")
	} else {
		d.Set("name", plugin.Name)
	}

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {

	err := meta.(*gokong.KongAdminClient).Plugins().DeleteById(d.Id())

	if err != nil {
		return fmt.Errorf("could not delete kong plugin: %v", err)
	}

	return nil
}

func createKongPluginRequestFromResourceData(d *schema.ResourceData) (*gokong.PluginRequest, error) {

	pluginRequest := &gokong.PluginRequest{}

	pluginRequest.Name = readStringFromResource(d, "name")
	pluginRequest.ApiId = readStringFromResource(d, "api_id")
	pluginRequest.ConsumerId = readStringFromResource(d, "consumer_id")
	pluginRequest.ServiceId = readStringFromResource(d, "service_id")
	pluginRequest.RouteId = readStringFromResource(d, "route_id")
	pluginRequest.Config = readMapFromResource(d, "config")

	if pluginRequest.Config == nil {
		if data, ok := d.GetOk("config_json"); ok {
			var configJson map[string]interface{}

			err := json.Unmarshal([]byte(data.(string)), &configJson)
			if err != nil {
				return pluginRequest, fmt.Errorf("failed to unmarshal config_json, err: %v", err)
			}

			pluginRequest.Config = configJson
		}
	}

	return pluginRequest, nil
}
