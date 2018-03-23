package gokong

import (
	"encoding/json"
	"fmt"
)

type ApiClient struct {
	config *Config
}

type ApiRequest struct {
	Name                   string   `json:"name"`
	Hosts                  []string `json:"hosts,omitempty"`
	Uris                   []string `json:"uris,omitempty"`
	Methods                []string `json:"methods,omitempty"`
	UpstreamUrl            string   `json:"upstream_url"`
	StripUri               bool     `json:"strip_uri"`
	PreserveHost           bool     `json:"preserve_host"`
	Retries                string   `json:"retries,omitempty"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout,omitempty"`
	HttpsOnly              bool     `json:"https_only"`
	HttpIfTerminated       bool     `json:"http_if_terminated"`
}

type Api struct {
	Id                     string   `json:"id"`
	CreatedAt              int      `json:"created_at"`
	Name                   string   `json:"name"`
	Hosts                  []string `json:"hosts,omitempty"`
	Uris                   []string `json:"uris,omitempty"`
	Methods                []string `json:"methods,omitempty"`
	UpstreamUrl            string   `json:"upstream_url"`
	StripUri               bool     `json:"strip_uri,omitempty"`
	PreserveHost           bool     `json:"preserve_host,omitempty"`
	Retries                int      `json:"retries,omitempty"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout,omitempty"`
	HttpsOnly              bool     `json:"https_only,omitempty"`
	HttpIfTerminated       bool     `json:"http_if_terminated,omitempty"`
}

type Apis struct {
	Results []*Api `json:"data,omitempty"`
	Total   int    `json:"total,omitempty"`
	Next    string `json:"next,omitempty"`
	Offset  string `json:"offset,omitempty"`
}

type ApiFilter struct {
	Id          string `url:"id,omitempty"`
	Name        string `url:"name,omitempty"`
	UpstreamUrl string `url:"upstream_url,omitempty"`
	Retries     int    `url:"retries,omitempty"`
	Size        int    `url:"size,omitempty"`
	Offset      int    `url:"offset,omitempty"`
}

const ApisPath = "/apis/"

func (apiClient *ApiClient) GetByName(name string) (*Api, error) {
	return apiClient.GetById(name)
}

func (apiClient *ApiClient) GetById(id string) (*Api, error) {
	_, body, errs := NewRequest(apiClient.config).Get(apiClient.config.HostAddress+ApisPath+id).Set("If-None-Match", `W/"wyzzy"`).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get api, error: %v", errs)
	}

	api := &Api{}
	err := json.Unmarshal([]byte(body), api)
	if err != nil {
		return nil, fmt.Errorf("could not parse api get response, error: %v", err)
	}

	if api.Id == "" {
		return nil, nil
	}

	return api, nil
}

func (apiClient *ApiClient) List() (*Apis, error) {
	return apiClient.ListFiltered(nil)
}

func (apiClient *ApiClient) ListFiltered(filter *ApiFilter) (*Apis, error) {

	address, err := addQueryString(apiClient.config.HostAddress+ApisPath, filter)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for apis filter, error: %v", err)
	}

	_, body, errs := NewRequest(apiClient.config).Get(address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get apis, error: %v", errs)
	}

	apis := &Apis{}
	err = json.Unmarshal([]byte(body), apis)
	if err != nil {
		return nil, fmt.Errorf("could not parse apis list response, error: %v", err)
	}

	return apis, nil
}

func (apiClient *ApiClient) Create(newApi *ApiRequest) (*Api, error) {

	_, body, errs := NewRequest(apiClient.config).Post(apiClient.config.HostAddress + ApisPath).Send(newApi).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new api, error: %v", errs)
	}

	createdApi := &Api{}
	err := json.Unmarshal([]byte(body), createdApi)
	if err != nil {
		return nil, fmt.Errorf("could not parse api creation response, error: %v %s", err, body)
	}

	if createdApi.Id == "" {
		return nil, fmt.Errorf("could not create api, error: %v", body)
	}

	return createdApi, nil
}

func (apiClient *ApiClient) DeleteByName(name string) error {
	return apiClient.DeleteById(name)
}

func (apiClient *ApiClient) DeleteById(id string) error {

	res, _, errs := NewRequest(apiClient.config).Delete(apiClient.config.HostAddress + ApisPath + id).End()
	if errs != nil {
		return fmt.Errorf("could not delete api, result: %v error: %v", res, errs)
	}

	return nil
}

func (apiClient *ApiClient) UpdateByName(name string, apiRequest *ApiRequest) (*Api, error) {
	return apiClient.UpdateById(name, apiRequest)
}

func (apiClient *ApiClient) UpdateById(id string, apiRequest *ApiRequest) (*Api, error) {

	_, body, errs := NewRequest(apiClient.config).Patch(apiClient.config.HostAddress + ApisPath + id).Send(apiRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update api, error: %v", errs)
	}

	updatedApi := &Api{}
	err := json.Unmarshal([]byte(body), updatedApi)
	if err != nil {
		return nil, fmt.Errorf("could not parse api update response, error: %v", err)
	}

	if updatedApi.Id == "" {
		return nil, fmt.Errorf("could not update certificate, error: %v", body)
	}

	return updatedApi, nil
}
