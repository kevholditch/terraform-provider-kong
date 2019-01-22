package gokong

import (
	"encoding/json"
	"fmt"
)

type RouteClient struct {
	config *Config
}

type RouteRequest struct {
	Name          *string   `json:"name"`
	Protocols     []*string `json:"protocols"`
	Methods       []*string `json:"methods"`
	Hosts         []*string `json:"hosts"`
	Paths         []*string `json:"paths"`
	RegexPriority *int      `json:"regex_priority"`
	StripPath     *bool     `json:"strip_path"`
	PreserveHost  *bool     `json:"preserve_host"`
	Snis          []*string `json:"snis"`
	Sources       []*string `json:"sources"`
	Destinations  []*string `json:"destinations"`
	Service       *Id       `json:"service"`
}

type Route struct {
	Id            *string   `json:"id"`
	Name          *string   `json:"name"`
	CreatedAt     *int      `json:"created_at"`
	UpdatedAt     *int      `json:"updated_at"`
	Protocols     []*string `json:"protocols"`
	Methods       []*string `json:"methods"`
	Hosts         []*string `json:"hosts"`
	Paths         []*string `json:"paths"`
	RegexPriority *int      `json:"regex_priority"`
	StripPath     *bool     `json:"strip_path"`
	PreserveHost  *bool     `json:"preserve_host"`
	Snis          []*string `json:"snis"`
	Sources       []*string `json:"sources"`
	Destinations  []*string `json:"destinations"`
	Service       *Id       `json:"service"`
}

type Routes struct {
	Data  []*Route `json:"data"`
	Total int      `json:"total"`
	Next  string   `json:"next"`
}

type RouteQueryString struct {
	Offset int
	Size   int
}

const RoutesPath = "/routes/"

func (routeClient *RouteClient) GetByName(name string) (*Route, error) {
	return routeClient.GetById(name)
}

func (routeClient *RouteClient) GetById(id string) (*Route, error) {
	r, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get the route, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	route := &Route{}
	err := json.Unmarshal([]byte(body), route)
	if err != nil {
		return nil, fmt.Errorf("could not parse route get response, error: %v", err)
	}

	if route.Id == nil {
		return nil, nil
	}

	return route, nil
}

func (routeClient *RouteClient) Create(routeRequest *RouteRequest) (*Route, error) {
	r, body, errs := newPost(routeClient.config, routeClient.config.HostAddress+RoutesPath).Send(routeRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not register the route, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	createdRoute := &Route{}
	err := json.Unmarshal([]byte(body), createdRoute)
	if err != nil {
		return nil, fmt.Errorf("could not parse route get response, error: %v", err)
	}

	if createdRoute.Id == nil {
		return nil, fmt.Errorf("could not register the route, error: %v", body)
	}

	return createdRoute, nil
}

func (routeClient *RouteClient) List(query *RouteQueryString) ([]*Route, error) {
	routes := []*Route{}
	data := &Routes{}

	if query.Size < 100 {
		query.Size = 100
	}

	if query.Size > 1000 {
		query.Size = 1000
	}

	for {
		r, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+RoutesPath).Query(query).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the route, error: %v", errs)
		}

		if r.StatusCode == 401 || r.StatusCode == 403 {
			return nil, fmt.Errorf("not authorised, message from kong: %s", body)
		}

		err := json.Unmarshal([]byte(body), data)
		if err != nil {
			return nil, fmt.Errorf("could not parse route get response, error: %v", err)
		}

		routes = append(routes, data.Data...)

		if data.Next == "" {
			break
		}

		query.Offset += query.Size
	}

	return routes, nil
}

func (routeClient *RouteClient) GetRoutesFromServiceName(name string) ([]*Route, error) {
	return routeClient.GetRoutesFromServiceId(name)
}

func (routeClient *RouteClient) GetRoutesFromServiceId(id string) ([]*Route, error) {
	routes := []*Route{}
	data := &Routes{}

	for {
		r, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+fmt.Sprintf("/services/%s/routes", id)).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the route, error: %v", errs)
		}

		if r.StatusCode == 401 || r.StatusCode == 403 {
			return nil, fmt.Errorf("not authorised, message from kong: %s", body)
		}

		err := json.Unmarshal([]byte(body), data)
		if err != nil {
			return nil, fmt.Errorf("could not parse route get response, error: %v", err)
		}

		routes = append(routes, data.Data...)

		if data.Next == "" {
			break
		}

	}
	return routes, nil
}

func (routeClient *RouteClient) UpdateByName(name string, routeRequest *RouteRequest) (*Route, error) {
	return routeClient.UpdateById(name, routeRequest)
}

func (routeClient *RouteClient) UpdateById(id string, routeRequest *RouteRequest) (*Route, error) {
	r, body, errs := newPatch(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).Send(routeRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update route, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	updatedRoute := &Route{}
	err := json.Unmarshal([]byte(body), updatedRoute)
	if err != nil {
		return nil, fmt.Errorf("could not parse route update response, error: %v", err)
	}

	if updatedRoute.Id == nil {
		return nil, fmt.Errorf("could not update route, error: %v", body)
	}

	return updatedRoute, nil
}

func (routeClient *RouteClient) DeleteByName(name string) error {
	return routeClient.DeleteById(name)
}

func (routeClient *RouteClient) DeleteById(id string) error {
	r, body, errs := newDelete(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete the route, result: %v error: %v", r, errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return fmt.Errorf("not authorised, message from kong: %s", body)
	}

	return nil
}
