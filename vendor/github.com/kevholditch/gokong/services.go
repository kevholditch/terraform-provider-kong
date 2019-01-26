package gokong

import (
	"encoding/json"
	"fmt"
)

type ServiceClient struct {
	config *Config
}

type ServiceRequest struct {
	Name           *string `json:"name"`
	Protocol       *string `json:"protocol"`
	Host           *string `json:"host"`
	Port           *int    `json:"port,omitempty"`
	Path           *string `json:"path,omitempty"`
	Retries        *int    `json:"retries,omitempty"`
	ConnectTimeout *int    `json:"connect_timeout,omitempty"`
	WriteTimeout   *int    `json:"write_timeout,omitempty"`
	ReadTimeout    *int    `json:"read_timeout,omitempty"`
	Url            *string `json:"url,omitempty"`
}

type Service struct {
	Id             *string `json:"id"`
	CreatedAt      *int    `json:"created_at"`
	UpdatedAt      *int    `json:"updated_at"`
	Protocol       *string `json:"protocol"`
	Host           *string `json:"host"`
	Port           *int    `json:"port"`
	Path           *string `json:"path"`
	Name           *string `json:"name"`
	Retries        *int    `json:"retries"`
	ConnectTimeout *int    `json:"connect_timeout"`
	WriteTimeout   *int    `json:"write_timeout"`
	ReadTimeout    *int    `json:"read_timeout"`
	Url            *string `json:"url"`
}

type Services struct {
	Data []*Service `json:"data"`
	Next *string    `json:"next"`
}

type ServiceQueryString struct {
	Offset int
	Size   int
}

const ServicesPath = "/services/"

func (serviceClient *ServiceClient) Create(serviceRequest *ServiceRequest) (*Service, error) {

	if serviceRequest.Port == nil {
		serviceRequest.Port = Int(80)
	}

	if serviceRequest.Retries == nil {
		serviceRequest.Retries = Int(5)
	}

	if serviceRequest.ConnectTimeout == nil {
		serviceRequest.ConnectTimeout = Int(60000)
	}

	if serviceRequest.ReadTimeout == nil {
		serviceRequest.ReadTimeout = Int(60000)
	}

	if serviceRequest.Retries == nil {
		serviceRequest.Retries = Int(60000)
	}

	r, body, errs := newPost(serviceClient.config, serviceClient.config.HostAddress+ServicesPath).Send(serviceRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not register the service, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	createdService := &Service{}
	err := json.Unmarshal([]byte(body), createdService)
	if err != nil {
		return nil, fmt.Errorf("could not parse service get response, error: %v", err)
	}

	if createdService.Id == nil {
		return nil, fmt.Errorf("could not register the service, error: %v", body)
	}

	return createdService, nil
}

func (serviceClient *ServiceClient) GetServiceByName(name string) (*Service, error) {
	return serviceClient.GetServiceById(name)
}

func (serviceClient *ServiceClient) GetServiceById(id string) (*Service, error) {
	return serviceClient.getService(serviceClient.config.HostAddress + ServicesPath + id)
}

func (serviceClient *ServiceClient) GetServiceFromRouteId(id string) (*Service, error) {
	return serviceClient.getService(serviceClient.config.HostAddress + "/routes/" + id + "/service")
}

func (serviceClient *ServiceClient) getService(endpoint string) (*Service, error) {
	r, body, errs := newGet(serviceClient.config, endpoint).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get the service, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	service := &Service{}
	err := json.Unmarshal([]byte(body), service)
	if err != nil {
		return nil, fmt.Errorf("could not parse service get response, error: %v", err)
	}

	if service.Id == nil {
		return nil, nil
	}

	return service, nil
}

func (serviceClient *ServiceClient) GetServices(query *ServiceQueryString) ([]*Service, error) {
	services := []*Service{}
	data := &Services{}

	if query.Size == 0 || query.Size < 100 {
		query.Size = 100
	}

	if query.Size > 1000 {
		query.Size = 1000
	}

	for {
		r, body, errs := newGet(serviceClient.config, serviceClient.config.HostAddress+ServicesPath).Query(query).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the service, error: %v", errs)
		}

		if r.StatusCode == 401 || r.StatusCode == 403 {
			return nil, fmt.Errorf("not authorised, message from kong: %s", body)
		}

		err := json.Unmarshal([]byte(body), data)
		if err != nil {
			return nil, fmt.Errorf("could not parse service get response, error: %v", err)
		}

		services = append(services, data.Data...)

		if data.Next == nil || *data.Next == "" {
			break
		}

		query.Offset += query.Size
	}

	return services, nil
}

func (serviceClient *ServiceClient) UpdateServiceByName(name string, serviceRequest *ServiceRequest) (*Service, error) {
	return serviceClient.UpdateServiceById(name, serviceRequest)
}

func (serviceClient *ServiceClient) UpdateServiceById(id string, serviceRequest *ServiceRequest) (*Service, error) {
	return serviceClient.updateService(serviceClient.config.HostAddress+ServicesPath+id, serviceRequest)
}

func (serviceClient *ServiceClient) UpdateServicebyRouteId(id string, serviceRequest *ServiceRequest) (*Service, error) {
	return serviceClient.updateService(serviceClient.config.HostAddress+"/routes/"+id+"/service", serviceRequest)
}

func (serviceClient *ServiceClient) updateService(endpoint string, serviceRequest *ServiceRequest) (*Service, error) {
	r, body, errs := newPatch(serviceClient.config, endpoint).Send(serviceRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update service, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	updatedService := &Service{}
	err := json.Unmarshal([]byte(body), updatedService)
	if err != nil {
		return nil, fmt.Errorf("could not parse service update response, error: %v", err)
	}

	if updatedService.Id == nil {
		return nil, fmt.Errorf("could not update service, error: %v", body)
	}

	return updatedService, nil
}

func (serviceClient *ServiceClient) DeleteServiceByName(name string) error {
	return serviceClient.DeleteServiceById(name)
}

func (serviceClient *ServiceClient) DeleteServiceById(id string) error {
	r, body, errs := newDelete(serviceClient.config, serviceClient.config.HostAddress+ServicesPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete the service, result: %v error: %v", r, errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return fmt.Errorf("not authorised, message from kong: %s", body)
	}

	return nil
}
