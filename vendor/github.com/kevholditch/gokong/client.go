package gokong

import (
	"github.com/google/go-querystring/query"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const EnvKongAdminHostAddress = "KONG_ADMIN_ADDR"

type KongAdminClient struct {
	config *Config
	client *gorequest.SuperAgent
}

type Config struct {
	HostAddress string
}

func addQueryString(currentUrl string, filter interface{}) (string, error) {
	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return currentUrl, nil
	}

	u, err := url.Parse(currentUrl)
	if err != nil {
		return currentUrl, err
	}

	qs, err := query.Values(filter)
	if err != nil {
		return currentUrl, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func NewDefaultConfig() *Config {
	config := &Config{
		HostAddress: "http://localhost:8001",
	}

	if os.Getenv(EnvKongAdminHostAddress) != "" {
		config.HostAddress = strings.TrimRight(os.Getenv(EnvKongAdminHostAddress), "/")
	}

	return config
}

func NewClient(config *Config) *KongAdminClient {
	return &KongAdminClient{
		config: config,
		client: gorequest.New(),
	}
}

func (kongAdminClient *KongAdminClient) Status() *StatusClient {
	return &StatusClient{
		config: kongAdminClient.config,
		client: kongAdminClient.client,
	}

}

func (kongAdminClient *KongAdminClient) Apis() *ApiClient {
	return &ApiClient{
		config: kongAdminClient.config,
		client: kongAdminClient.client,
	}
}

func (kongAdminClient *KongAdminClient) Consumers() *ConsumerClient {
	return &ConsumerClient{
		config: kongAdminClient.config,
		client: kongAdminClient.client,
	}
}
