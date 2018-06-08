package gokong

import (
	"encoding/json"
	"fmt"
)

type ServiceClient struct {
	config *Config
}

type ServiceRequest struct {
	Name           string `json:"name"`
	Protocol       string `json:"protocol"`
	Host           string `json:"host"`
	Port           int    `json:"port,omitempty"`
	Path           string `json:"path,omitempty"`
	Retries        int    `json:"retries,omitempty"`
	ConnectTimeout int    `json:"connect_timeout,omitempty"`
	WriteTimeout   int    `json:"write_timeout,omitempty"`
	ReadTimeout    int    `json:"read_timeout,omitempty"`
}

type Service struct {
	Id             string `json:"id"`
	CreatedAt      int    `json:"created_at"`
	UpdatedAt      int    `json:"updated_at"`
	Protocol       string `json:"protocol"`
	Host           string `json:"host"`
	Port           int    `json:"int"`
	Path           string `json:"path"`
	Name           string `json:"name"`
	Retries        int    `json:"retries"`
	ConnectTimeout int    `json:"connect_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	ReadTimeout    int    `json:"read_timeout"`
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

func (serviceClient *ServiceClient) AddService(serviceRequest *ServiceRequest) (*Service, error) {

	if serviceRequest.Port == 0 {
		serviceRequest.Port = 80
	}

	if serviceRequest.Retries == 0 {
		serviceRequest.Retries = 5
	}

	if serviceRequest.ConnectTimeout == 0 {
		serviceRequest.ConnectTimeout = 60000
	}

	if serviceRequest.ReadTimeout == 0 {
		serviceRequest.ReadTimeout = 60000
	}

	if serviceRequest.Retries == 0 {
		serviceRequest.Retries = 60000
	}

	_, body, errs := NewRequest(serviceClient.config).Post(serviceClient.config.HostAddress + ServicesPath).Send(serviceRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not register the service, error: %v", errs)
	}

	createdService := &Service{}
	err := json.Unmarshal([]byte(body), createdService)
	if err != nil {
		return nil, fmt.Errorf("could not parse service get response, error: %v", err)
	}

	if createdService.Id == "" {
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
	_, body, errs := NewRequest(serviceClient.config).Get(endpoint).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get the service, error: %v", errs)
	}

	service := &Service{}
	err := json.Unmarshal([]byte(body), service)
	if err != nil {
		return nil, fmt.Errorf("could not parse service get response, error: %v", err)
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
		_, body, errs := NewRequest(serviceClient.config).Get(serviceClient.config.HostAddress + ServicesPath).Query(query).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the service, error: %v", errs)
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
	_, body, errs := NewRequest(serviceClient.config).Patch(endpoint).Send(serviceRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update service, error: %v", errs)
	}

	updatedService := &Service{}
	err := json.Unmarshal([]byte(body), updatedService)
	if err != nil {
		return nil, fmt.Errorf("could not parse service update response, error: %v", err)
	}

	if updatedService.Id == "" {
		return nil, fmt.Errorf("could not update service, error: %v", body)
	}

	return updatedService, nil
}

func (serviceClient *ServiceClient) DeleteServiceByName(name string) error {
	return serviceClient.DeleteServiceById(name)
}

func (serviceClient *ServiceClient) DeleteServiceById(id string) error {
	res, _, errs := NewRequest(serviceClient.config).Delete(serviceClient.config.HostAddress + ServicesPath + id).End()
	if errs != nil {
		return fmt.Errorf("could not delete the service, result: %v error: %v", res, errs)
	}

	return nil
}
