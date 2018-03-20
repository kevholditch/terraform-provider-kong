package gokong

import (
	"github.com/parnurzeal/gorequest"
)

func NewRequest(adminConfig *Config) *gorequest.SuperAgent {
	request := gorequest.New()
	if adminConfig.Username != "" || adminConfig.Password != "" {
		request.SetBasicAuth(adminConfig.Username, adminConfig.Password)
	}
	return request
}
