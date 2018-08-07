package gokong

import (
	"encoding/json"
	"fmt"
)

type PluginClient struct {
	config *Config
}

type PluginRequest struct {
	Name       string                 `json:"name"`
	ApiId      string                 `json:"api_id,omitempty"`
	ConsumerId string                 `json:"consumer_id,omitempty"`
	ServiceId  string                 `json:"service_id,omitempty"`
	RouteId    string                 `json:"route_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

type Plugin struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	ApiId      string                 `json:"api_id,omitempty"`
	ConsumerId string                 `json:"consumer_id,omitempty"`
	ServiceId  string                 `json:"service_id,omitempty"`
	RouteId    string                 `json:"route_id,omitempty"`
	Config     map[string]interface{} `json:"config,omitempty"`
	Enabled    bool                   `json:"enabled,omitempty"`
}

type Plugins struct {
	Results []*Plugin `json:"data,omitempty"`
	Total   int       `json:"total,omitempty"`
	Next    string    `json:"next,omitempty"`
}

type PluginFilter struct {
	Id         string `url:"id,omitempty"`
	Name       string `url:"name,omitempty"`
	ApiId      string `url:"api_id,omitempty"`
	ConsumerId string `url:"consumer_id,omitempty"`
	ServiceId  string `url:"service_id,omitempty"`
	RouteId    string `url:"route_id,omitempty"`
	Size       int    `url:"size,omitempty"`
	Offset     int    `url:"offset,omitempty"`
}

const PluginsPath = "/plugins/"

func (pluginClient *PluginClient) GetById(id string) (*Plugin, error) {

	r, body, errs := newGet(pluginClient.config, pluginClient.config.HostAddress+PluginsPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugin := &Plugin{}
	err := json.Unmarshal([]byte(body), plugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin plugin response, error: %v", err)
	}

	if plugin.Id == "" {
		return nil, nil
	}

	return plugin, nil
}

func (pluginClient *PluginClient) List() (*Plugins, error) {
	return pluginClient.ListFiltered(nil)
}

func (pluginClient *PluginClient) ListFiltered(filter *PluginFilter) (*Plugins, error) {

	address, err := addQueryString(pluginClient.config.HostAddress+PluginsPath, filter)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for plugins filter, error: %v", err)
	}

	r, body, errs := newGet(pluginClient.config, address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get plugins, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	plugins := &Plugins{}
	err = json.Unmarshal([]byte(body), plugins)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugins list response, error: %v", err)
	}

	return plugins, nil
}

func (pluginClient *PluginClient) Create(pluginRequest *PluginRequest) (*Plugin, error) {

	r, body, errs := newPost(pluginClient.config, pluginClient.config.HostAddress+PluginsPath).Send(pluginRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	createdPlugin := &Plugin{}
	err := json.Unmarshal([]byte(body), createdPlugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin creation response, error: %v kong response: %s", err, body)
	}

	if createdPlugin.Id == "" {
		return nil, fmt.Errorf("could not create plugin, err: %v", body)
	}

	return createdPlugin, nil
}

func (pluginClient *PluginClient) UpdateById(id string, pluginRequest *PluginRequest) (*Plugin, error) {

	r, body, errs := newPatch(pluginClient.config, pluginClient.config.HostAddress+PluginsPath+id).Send(pluginRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update plugin, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	updatedPlugin := &Plugin{}
	err := json.Unmarshal([]byte(body), updatedPlugin)
	if err != nil {
		return nil, fmt.Errorf("could not parse plugin update response, error: %v kong response: %s", err, body)
	}

	if updatedPlugin.Id == "" {
		return nil, fmt.Errorf("could not update plugin, error: %v", body)
	}

	return updatedPlugin, nil
}

func (pluginClient *PluginClient) DeleteById(id string) error {

	r, body, errs := newDelete(pluginClient.config, pluginClient.config.HostAddress+PluginsPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete plugin, result: %v error: %v", r, errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return fmt.Errorf("not authorised, message from kong: %s", body)
	}

	return nil
}
