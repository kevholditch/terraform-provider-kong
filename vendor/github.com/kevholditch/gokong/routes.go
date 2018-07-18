package gokong

import (
	"encoding/json"
	"fmt"
)

type RouteClient struct {
	config *Config
}

type RouteRequest struct {
	Protocols    []*string           `json:"protocols"`
	Methods      []*string           `json:"methods,omitempty"`
	Hosts        []*string           `json:"hosts,omitempty"`
	Paths        []*string           `json:"paths,omitempty"`
	StripPath    *bool               `json:"strip_path,omitempty"`
	PreserveHost *bool               `json:"preserve_host,omitempty"`
	Service      *RouteServiceObject `json:"service"`
}

type RouteServiceObject struct {
	Id string `json:"id"`
}

type Route struct {
	Id            *string             `json:"id"`
	CreatedAt     *int                `json:"created_at"`
	UpdatedAt     *int                `json:"updated_at"`
	Protocols     []*string           `json:"protocols"`
	Methods       []*string           `json:"methods"`
	Hosts         []*string           `json:"hosts"`
	Paths         []*string           `json:"paths"`
	RegexPriority *int                `json:"regex_priority"`
	StripPath     *bool               `json:"strip_path"`
	PreserveHost  *bool               `json:"preserve_host"`
	Service       *RouteServiceObject `json:"service"`
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

func (routeClient *RouteClient) AddRoute(routeRequest *RouteRequest) (*Route, error) {
	_, body, errs := newPost(routeClient.config, routeClient.config.HostAddress+RoutesPath).Send(routeRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not register the route, error: %v", errs)
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

func (routeClient *RouteClient) GetRoute(id string) (*Route, error) {
	_, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get the route, error: %v", errs)
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

func (routeClient *RouteClient) GetRoutes(query *RouteQueryString) ([]*Route, error) {
	routes := []*Route{}
	data := &Routes{}

	if query.Size < 100 {
		query.Size = 100
	}

	if query.Size > 1000 {
		query.Size = 1000
	}

	for {
		_, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+RoutesPath).Query(query).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the route, error: %v", errs)
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
		_, body, errs := newGet(routeClient.config, routeClient.config.HostAddress+fmt.Sprintf("/services/%s/routes", id)).End()
		if errs != nil {
			return nil, fmt.Errorf("could not get the route, error: %v", errs)
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

func (routeClient *RouteClient) UpdateRoute(id string, routeRequest *RouteRequest) (*Route, error) {
	_, body, errs := newPatch(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).Send(routeRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update route, error: %v", errs)
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

func (routeClient *RouteClient) DeleteRoute(id string) error {
	res, _, errs := newDelete(routeClient.config, routeClient.config.HostAddress+RoutesPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete the route, result: %v error: %v", res, errs)
	}

	return nil
}
