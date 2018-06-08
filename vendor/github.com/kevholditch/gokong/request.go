package gokong

import (
	"crypto/tls"

	"github.com/parnurzeal/gorequest"
)

func NewRequest(adminConfig *Config) *gorequest.SuperAgent {
	request := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: adminConfig.InsecureSkipVerify})
	if adminConfig.Username != "" || adminConfig.Password != "" {
		request.SetBasicAuth(adminConfig.Username, adminConfig.Password)
	}
	return request
}
