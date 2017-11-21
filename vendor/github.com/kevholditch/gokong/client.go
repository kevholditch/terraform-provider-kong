package gokong

import (
	"github.com/parnurzeal/gorequest"
	"os"
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
