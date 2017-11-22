package gokong

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"reflect"
)

type ApiClient struct {
	config *Config
	client *gorequest.SuperAgent
}

type ApiRequest struct {
	Name                   string   `json:"name"`
	Hosts                  []string `json:"hosts,omitempty"`
	Uris                   []string `json:"uris,omitempty"`
	Methods                []string `json:"methods,omitempty"`
	UpstreamUrl            string   `json:"upstream_url"`
	StripUri               bool     `json:"strip_uri"`
	PreserveHost           bool     `json:"preserve_host,omitempty"`
	Retries                int      `json:"retries,omitempty"`
	UpstreamConnectTimeout int      `json:"upstream_connect_timeout,omitempty"`
	UpstreamSendTimeout    int      `json:"upstream_send_timeout,omitempty"`
	UpstreamReadTimeout    int      `json:"upstream_read_timeout,omitempty"`
	HttpsOnly              bool     `json:"https_only,omitempty"`
	HttpIfTerminated       bool     `json:"http_if_terminated,omitempty"`
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

type GetAllFilter struct {
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

	_, body, errs := apiClient.client.Get(apiClient.config.HostAddress + ApisPath + id).End()
	if errs != nil {
		return nil, errors.New(fmt.Sprintf("Could not get api, error: %v", errs))
	}

	api := &Api{}
	err := json.Unmarshal([]byte(body), api)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse api get response, error: %v", err))
	}

	if api.Id == "" {
		return nil, nil
	}

	return api, nil
}

func (apiClient *ApiClient) GetAll() (*Apis, error) {
	return apiClient.GetAllFiltered(nil)
}

func addQueryString(currentUrl string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	u, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return currentUrl, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (apiClient *ApiClient) GetAllFiltered(filter *GetAllFilter) (*Apis, error) {

	address, err := addQueryString(apiClient.config.HostAddress+ApisPath, filter)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not build query string for apis filter, error: %v", err))
	}

	_, body, errs := apiClient.client.Get(address).End()
	if errs != nil {
		return nil, errors.New(fmt.Sprintf("Could not get apis, error: %v", errs))
	}

	apis := &Apis{}
	err = json.Unmarshal([]byte(body), apis)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse apis get response, error: %v", err))
	}

	return apis, nil
}

func (apiClient *ApiClient) Create(newApi *ApiRequest) (*Api, error) {

	_, body, errs := apiClient.client.Post(apiClient.config.HostAddress + ApisPath).Send(newApi).End()
	if errs != nil {
		return nil, errors.New(fmt.Sprintf("Could not create new api, error: %v", errs))
	}

	createdApi := &Api{}
	err := json.Unmarshal([]byte(body), createdApi)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse api creation response, error: %v", err))
	}

	return createdApi, nil
}

func (apiClient *ApiClient) DeleteByName(name string) error {
	return apiClient.DeleteById(name)
}

func (apiClient *ApiClient) DeleteById(id string) error {

	res, _, errs := apiClient.client.Delete(apiClient.config.HostAddress + ApisPath + id).End()
	if errs != nil {
		return errors.New(fmt.Sprintf("Could not delete api, result: %v error: %v", res, errs))
	}

	return nil
}

func (apiClient *ApiClient) UpdateByName(name string, apiRequest *ApiRequest) (*Api, error) {
	return apiClient.UpdateById(name, apiRequest)
}

func (apiClient *ApiClient) UpdateById(id string, apiRequest *ApiRequest) (*Api, error) {

	_, body, errs := apiClient.client.Patch(apiClient.config.HostAddress + ApisPath + id).Send(apiRequest).End()
	if errs != nil {
		return nil, errors.New(fmt.Sprintf("Could not update api, error: %v", errs))
	}

	updatedApi := &Api{}
	err := json.Unmarshal([]byte(body), updatedApi)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse api update response, error: %v", err))
	}

	return updatedApi, nil
}
