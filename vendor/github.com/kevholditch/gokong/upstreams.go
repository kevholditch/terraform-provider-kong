package gokong

import (
	"encoding/json"
	"fmt"
)

type UpstreamClient struct {
	config *Config
}

type UpstreamRequest struct {
	Name               string               `json:"name"`
	Slots              int                  `json:"slots,omitempty"`
	HashOn             string               `json:"hash_on,omitempty"`
	HashFallback       string               `json:"hash_fallback,omitempty"`
	HashOnHeader       string               `json:"hash_on_header,omitempty"`
	HashFallbackHeader string               `json:"hash_fallback_header,omitempty"`
	HealthChecks       *UpstreamHealthCheck `json:"healthchecks,omitempty"`
}

type UpstreamHealthCheck struct {
	Active  *UpstreamHealthCheckActive  `json:"active,omitempty"`
	Passive *UpstreamHealthCheckPassive `json:"passive,omitempty"`
}

type UpstreamHealthCheckActive struct {
	Concurrency int              `json:"concurrency,omitempty"`
	Healthy     *ActiveHealthy   `json:"healthy,omitempty"`
	HttpPath    string           `json:"http_path,omitempty"`
	Timeout     int              `json:"timeout,omitempty"`
	Unhealthy   *ActiveUnhealthy `json:"unhealthy,omitempty"`
}

type ActiveHealthy struct {
	HttpStatuses []int `json:"http_statuses,omitempty"`
	Interval     int   `json:"interval,omitempty"`
	Successes    int   `json:"successes,omitempty"`
}

type ActiveUnhealthy struct {
	HttpFailures int   `json:"http_failures,omitempty"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
	Interval     int   `json:"interval,omitempty"`
	TcpFailures  int   `json:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty"`
}

type UpstreamHealthCheckPassive struct {
	Healthy   *PassiveHealthy   `json:"healthy,omitempty"`
	Unhealthy *PassiveUnhealthy `json:"unhealthy,omitempty"`
}

type PassiveHealthy struct {
	HttpStatuses []int `json:"http_statuses,omitempty"`
	Successes    int   `json:"successes,omitempty"`
}

type PassiveUnhealthy struct {
	HttpFailures int   `json:"http_failures,omitempty"`
	HttpStatuses []int `json:"http_statuses,omitempty"`
	TcpFailures  int   `json:"tcp_failures,omitempty"`
	Timeouts     int   `json:"timeouts,omitempty"`
}

type Upstream struct {
	Id string `json:"id,omitempty"`
	UpstreamRequest
}

type Upstreams struct {
	Results []*Upstream `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
	Next    string      `json:"next,omitempty"`
	Offset  string      `json:"offset,omitempty"`
}

type UpstreamFilter struct {
	Id     string `url:"id,omitempty"`
	Name   string `url:"name,omitempty"`
	Slots  int    `url:"slots,omitempty"`
	Size   int    `url:"size,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

const UpstreamsPath = "/upstreams/"

func (upstreamClient *UpstreamClient) GetByName(name string) (*Upstream, error) {
	return upstreamClient.GetById(name)
}

func (upstreamClient *UpstreamClient) GetById(id string) (*Upstream, error) {

	_, body, errs := newGet(upstreamClient.config, upstreamClient.config.HostAddress+UpstreamsPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get upstream, error: %v", errs)
	}

	upstream := &Upstream{}
	err := json.Unmarshal([]byte(body), upstream)
	if err != nil {
		return nil, fmt.Errorf("could not parse upstream get response, error: %v", err)
	}

	if upstream.Id == "" {
		return nil, nil
	}

	return upstream, nil
}

func (upstreamClient *UpstreamClient) Create(upstreamRequest *UpstreamRequest) (*Upstream, error) {

	_, body, errs := newPost(upstreamClient.config, upstreamClient.config.HostAddress+UpstreamsPath).Send(upstreamRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new upstream, error: %v", errs)
	}

	createdUpstream := &Upstream{}
	err := json.Unmarshal([]byte(body), createdUpstream)
	if err != nil {
		return nil, fmt.Errorf("could not parse upstream creation response, error: %v", err)
	}

	if createdUpstream.Id == "" {
		return nil, fmt.Errorf("could not create update, error: %v", body)
	}

	return createdUpstream, nil
}

func (upstreamClient *UpstreamClient) DeleteByName(name string) error {
	return upstreamClient.DeleteById(name)
}

func (upstreamClient *UpstreamClient) DeleteById(id string) error {

	res, _, errs := newDelete(upstreamClient.config, upstreamClient.config.HostAddress+UpstreamsPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete upstream, result: %v error: %v", res, errs)
	}

	return nil
}

func (upstreamClient *UpstreamClient) List() (*Upstreams, error) {
	return upstreamClient.ListFiltered(nil)
}

func (upstreamClient *UpstreamClient) ListFiltered(filter *UpstreamFilter) (*Upstreams, error) {

	address, err := addQueryString(upstreamClient.config.HostAddress+UpstreamsPath, filter)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for upstreams filter, error: %v", err)
	}

	_, body, errs := newGet(upstreamClient.config, address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get upstreams, error: %v", errs)
	}

	upstreams := &Upstreams{}
	err = json.Unmarshal([]byte(body), upstreams)
	if err != nil {
		return nil, fmt.Errorf("could not parse upstreams list response, error: %v", err)
	}

	return upstreams, nil
}

func (upstreamClient *UpstreamClient) UpdateByName(name string, upstreamRequest *UpstreamRequest) (*Upstream, error) {
	return upstreamClient.UpdateById(name, upstreamRequest)
}

func (upstreamClient *UpstreamClient) UpdateById(id string, upstreamRequest *UpstreamRequest) (*Upstream, error) {

	_, body, errs := newPatch(upstreamClient.config, upstreamClient.config.HostAddress+UpstreamsPath+id).Send(upstreamRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update upstream, error: %v", errs)
	}

	updatedUpstream := &Upstream{}
	err := json.Unmarshal([]byte(body), updatedUpstream)
	if err != nil {
		return nil, fmt.Errorf("could not parse upstream update response, error: %v", err)
	}

	if updatedUpstream.Id == "" {
		return nil, fmt.Errorf("could not update upstream, error: %v", body)
	}

	return updatedUpstream, nil
}
