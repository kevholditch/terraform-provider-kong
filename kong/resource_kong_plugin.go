package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongPluginCreate,
		ReadContext:   resourceKongPluginRead,
		DeleteContext: resourceKongPluginDelete,
		UpdateContext: resourceKongPluginUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"consumer_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"route_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"config_json": {
				Type:         schema.TypeString,
				Optional:     true,
				StateFunc:    normalizeDataJSON,
				ValidateFunc: validateDataJSON,
				Description:  "plugin configuration in JSON format, configuration must be a valid JSON object.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"strict_match": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"computed_config": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceKongPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.Plugins
	plugin, err := client.Create(ctx, pluginRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong plugin: %v error: %v", pluginRequest, err))
	}

	d.SetId(*plugin.ID)

	return resourceKongPluginRead(ctx, d, meta)
}

func resourceKongPluginUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.Partial(false)

	pluginRequest, err := createKongPluginRequestFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*config).adminClient.Plugins
	_, err = client.Update(ctx, pluginRequest)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating kong plugin: %s", err))
	}

	return resourceKongPluginRead(ctx, d, meta)
}

func resourceKongPluginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Plugins
	plugin, err := client.Get(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong plugin: %v", err))
	}

	if plugin == nil {
		d.SetId("")
	} else {
		d.SetId(*plugin.ID)
		err = d.Set("name", plugin.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		if plugin.Service != nil {
			err = d.Set("service_id", plugin.Service.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if plugin.Route != nil {
			err = d.Set("route_id", plugin.Route.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if plugin.Consumer != nil {
			err = d.Set("consumer_id", plugin.Consumer.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		err = d.Set("enabled", plugin.Enabled)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("tags", plugin.Tags)
		if err != nil {
			return diag.FromErr(err)
		}

		// We sync this property from upstream as a method to allow you to import a resource with the config tracked in
		// terraform state. We do not track `config` as it will be a source of a perpetual diff.
		// https://www.terraform.io/docs/extend/best-practices/detecting-drift.html#capture-all-state-in-read
		upstreamJSON := pluginConfigJSONToString(plugin.Config)
		setConfig := func(strict bool) error {
			if strict {
				err := d.Set("config_json", upstreamJSON)
				if err != nil {
					return err
				}
			} else {
				err := d.Set("computed_config", upstreamJSON)
				if err != nil {
					return err
				}
			}
			return nil
		}
		if value, ok := d.GetOk("strict_match"); ok {
			err := setConfig(value.(bool))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := setConfig(meta.(*config).strictPlugins)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

func resourceKongPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.Plugins
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong plugin: %v", err))
	}

	return diags
}

func createKongPluginRequestFromResourceData(d *schema.ResourceData) (*kong.Plugin, error) {

	pluginRequest := &kong.Plugin{}
	// Build Consumer Configuration
	consumerID := readIdPtrFromResource(d, "consumer_id")
	if consumerID != nil {
		pluginRequest.Consumer = &kong.Consumer{
			ID: consumerID,
		}
	}
	// Build Service Configuration
	serviceID := readIdPtrFromResource(d, "service_id")
	if serviceID != nil {
		pluginRequest.Service = &kong.Service{
			ID: serviceID,
		}
	}
	// Build Route Configuration
	routeID := readIdPtrFromResource(d, "route_id")
	if routeID != nil {
		pluginRequest.Route = &kong.Route{
			ID: routeID,
		}
	}
	if d.Id() != "" {
		pluginRequest.ID = kong.String(d.Id())
	}

	pluginRequest.Name = readStringPtrFromResource(d, "name")
	pluginRequest.Enabled = readBoolPtrFromResource(d, "enabled")
	pluginRequest.Tags = readStringArrayPtrFromResource(d, "tags")

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

func validateDataJSON(configI interface{}, _ string) ([]string, []error) {
	dataJSON := configI.(string)
	dataMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(dataJSON), &dataMap)
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
}

func normalizeDataJSON(configI interface{}) string {
	dataJSON := configI.(string)

	dataMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(dataJSON), &dataMap)
	if err != nil {
		// The validate function should've taken care of this.
		log.Printf("[ERROR] Invalid JSON data in config_json: %s", err)
		return ""
	}

	ret, err := json.Marshal(dataMap)
	if err != nil {
		// Should never happen.
		log.Printf("[ERROR] Problem normalizing JSON for config_json: %s", err)
		return dataJSON
	}

	return string(ret)
}
