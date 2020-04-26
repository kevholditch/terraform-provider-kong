package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hbagdi/go-kong/kong"
)

func resourceKongConsumerPluginConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceKongConsumerPluginConfigCreate,
		Read:   resourceKongConsumerPluginConfigRead,
		Delete: resourceKongConsumerPluginConfigDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"consumer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plugin_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// Suppress diff when config is empty so we can sync with upstream always
			// The ForceNew property is what makes this work.
			"config_json": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				StateFunc:    normalizeDataJSON,
				ValidateFunc: validateDataJSON,
				Description:  "JSON format of plugin config",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"computed_config": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type idFields struct {
	consumerID string
	pluginName string
	id         string
}

func validateDataJSON(configI interface{}, k string) ([]string, []error) {
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

func buildID(consumerID, pluginName, configID string) string {
	return consumerID + "|" + pluginName + "|" + configID
}

func splitIDIntoFields(id string) (*idFields, error) {
	idSplit := strings.Split(id, "|")

	if len(idSplit) != 3 {
		return nil, fmt.Errorf("failed to calculate consumer plugin config id, should be pipe separated as consumerId|pluginName|id found: %v", id)
	}

	return &idFields{
		consumerID: idSplit[0],
		pluginName: idSplit[1],
		id:         idSplit[2],
	}, nil
}

func resourceKongConsumerPluginConfigCreate(d *schema.ResourceData, meta interface{}) error {

	consumerID := readStringFromResource(d, "consumer_id")
	pluginName := readStringFromResource(d, "plugin_name")
	configJSON := readStringFromResource(d, "config_json")

	client := meta.(*config).adminClient.Plugins

	config := kong.Configuration{}
	err := json.Unmarshal([]byte(configJSON), &config)
	plugin := &kong.Plugin{
		Consumer: &kong.Consumer{
			ID: kong.String(consumerID),
		},
		Name:   kong.String(pluginName),
		Config: config,
	}

	consumerPluginConfig, err := client.Create(context.Background(), plugin)

	if err != nil {
		return fmt.Errorf("failed to create kong consumer plugin config, error: %v", err)
	}

	if consumerPluginConfig == nil {
		d.SetId("")
	} else {
		d.SetId(buildID(consumerID, pluginName, *consumerPluginConfig.ID))
	}

	return resourceKongConsumerPluginConfigRead(d, meta)
}

func resourceKongConsumerPluginConfigRead(d *schema.ResourceData, meta interface{}) error {

	idFields, err := splitIDIntoFields(d.Id())

	if err != nil {
		return err
	}

	// First check if the consumer exists. If it does not then the consumer plugin no longer exists either.
	consumerClient := meta.(*config).adminClient.Consumers
	if consumer, _ := consumerClient.Get(context.Background(), kong.String(idFields.consumerID)); consumer == nil {
		d.SetId("")
		return nil
	}

	client := meta.(*config).adminClient.Plugins

	plugin, err := client.Get(context.Background(), kong.String(idFields.id))
	if err != nil {
		return fmt.Errorf("could not find kong consumer plugin config with id: %s error: %v", d.Id(), err)
	}

	if plugin.Config == nil {
		return fmt.Errorf("could not configure plugin for kong consumer")
	}

	d.Set("consumer_id", idFields.consumerID)
	d.Set("plugin_name", idFields.pluginName)

	// We sync this property from upstream as a method to allow you to import a resource with the config tracked in
	// terraform state. We do not track `config` as it will be a source of a perpetual diff.
	// https://www.terraform.io/docs/extend/best-practices/detecting-drift.html#capture-all-state-in-read
	upstreamJSON, err := consumerPluginConfigJSONToString(plugin.Config)
	if err != nil {
		return fmt.Errorf("could not read in consumer plugin config body: %s error: %v", d.Id(), err)
	}

	d.Set("computed_config", upstreamJSON)

	return nil
}

func resourceKongConsumerPluginConfigDelete(d *schema.ResourceData, meta interface{}) error {

	idFields, err := splitIDIntoFields(d.Id())

	if err != nil {
		return err
	}

	client := meta.(*config).adminClient.Plugins

	client.Delete(context.Background(), kong.String(idFields.id))

	if err != nil {
		return fmt.Errorf("could not delete kong consumer plugin config: %v", err)
	}

	return nil
}

// consumerPluginConfigJsonToString removes computed properties
func consumerPluginConfigJSONToString(config kong.Configuration) (string, error) {
	marshalledData := map[string]interface{}{}

	for key, val := range config {
		if !contains(computedPluginProperties, key) {
			marshalledData[key] = val
		}
	}
	rawJSON, _ := json.Marshal(marshalledData)

	return string(rawJSON), nil
}
