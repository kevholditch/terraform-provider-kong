package kong

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hbagdi/go-kong/kong"
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
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"config_json": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				StateFunc:    normalizeDataJSON,
				ValidateFunc: validateDataJSON,
				Description:  "plugin configuration in JSON format, configuration must be a valid JSON object.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"strict_match": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"computed_config": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKongPluginCreate(d *schema.ResourceData, meta interface{}) error {

	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.Plugins
	plugin, err := client.Create(context.Background(), pluginRequest)

	if err != nil {
		return fmt.Errorf("failed to create kong plugin: %v error: %v", pluginRequest, err)
	}

	d.SetId(*plugin.ID)

	return resourceKongPluginRead(d, meta)
}

func resourceKongPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(false)

	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.Plugins
	_, err = client.Update(context.Background(), pluginRequest)

	if err != nil {
		return fmt.Errorf("error updating kong plugin: %s", err)
	}

	return resourceKongPluginRead(d, meta)
}

func resourceKongPluginRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Plugins
	plugin, err := client.Get(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not find kong plugin: %v", err)
	}

	if plugin == nil {
		d.SetId("")
	} else {
		d.Set("name", plugin.Name)
		d.Set("service_id", plugin.Service.ID)
		d.Set("route_id", plugin.Route.ID)
		d.Set("consumer_id", plugin.Consumer.ID)
		d.Set("enabled", plugin.Enabled)

		// We sync this property from upstream as a method to allow you to import a resource with the config tracked in
		// terraform state. We do not track `config` as it will be a source of a perpetual diff.
		// https://www.terraform.io/docs/extend/best-practices/detecting-drift.html#capture-all-state-in-read
		upstreamJSON := pluginConfigJSONToString(plugin.Config)
		setConfig := func(strict bool) {
			if strict {
				d.Set("config_json", upstreamJSON)
			} else {
				d.Set("computed_config", upstreamJSON)
			}
		}
		if value, ok := d.GetOk("strict_match"); ok {
			setConfig(value.(bool))
		} else {
			setConfig(meta.(*config).strictPlugins)
		}
	}

	return nil
}

func resourceKongPluginDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*config).adminClient.Plugins
	err := client.Delete(context.Background(), kong.String(d.Id()))

	if err != nil {
		return fmt.Errorf("could not delete kong plugin: %v", err)
	}

	return nil
}

func createKongPluginRequestFromResourceData(d *schema.ResourceData) (*kong.Plugin, error) {

	pluginRequest := &kong.Plugin{}

	pluginRequest.Name = readStringPtrFromResource(d, "name")
	pluginRequest.Consumer.ID = readIdPtrFromResource(d, "consumer_id")
	pluginRequest.Service.ID = readIdPtrFromResource(d, "service_id")
	pluginRequest.Route.ID = readIdPtrFromResource(d, "route_id")
	pluginRequest.Enabled = readBoolPtrFromResource(d, "enabled")

	if data, ok := d.GetOk("config_json"); ok {
		var configJSON map[string]interface{}

		err := json.Unmarshal([]byte(data.(string)), &configJSON)
		if err != nil {
			return pluginRequest, fmt.Errorf("failed to unmarshal config_json, err: %v", err)
		}

		pluginRequest.Config = configJSON
	}

	return pluginRequest, nil
}

// Since this config is a schemaless "blob" we have to remove computed properties
func pluginConfigJSONToString(data map[string]interface{}) string {
	marshalledData := map[string]interface{}{}
	for key, val := range data {
		if !contains(computedPluginProperties, key) {
			marshalledData[key] = val
		}
	}
	// We know it is valid JSON at this point
	rawJSON, _ := json.Marshal(marshalledData)

	return string(rawJSON)
}
