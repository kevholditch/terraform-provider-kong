package gokong

import (
	"encoding/json"
	"fmt"
)

type SnisClient struct {
	config *Config
}

type SnisRequest struct {
	Name             string `json:"name,omitempty"`
	SslCertificateId string `json:"ssl_certificate_id,omitempty"`
}

type Sni struct {
	Name             string `json:"name,omitempty"`
	SslCertificateId string `json:"ssl_certificate_id,omitempty"`
}

type Snis struct {
	Results []*Sni `json:"data,omitempty"`
	Total   int    `json:"total,omitempty"`
}

const SnisPath = "/snis/"

func (snisClient *SnisClient) Create(snisRequest *SnisRequest) (*Sni, error) {

	_, body, errs := newPost(snisClient.config, snisClient.config.HostAddress+SnisPath).Send(snisRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new sni, error: %v", errs)
	}

	sni := &Sni{}
	err := json.Unmarshal([]byte(body), sni)
	if err != nil {
		return nil, fmt.Errorf("could not parse sni creation response, error: %v", err)
	}

	if sni.SslCertificateId == "" {
		return nil, fmt.Errorf("could not create sni, error: %v", body)
	}

	return sni, nil
}

func (snisClient *SnisClient) GetByName(name string) (*Sni, error) {

	_, body, errs := newGet(snisClient.config, snisClient.config.HostAddress+SnisPath+name).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get sni, error: %v", errs)
	}

	sni := &Sni{}
	err := json.Unmarshal([]byte(body), sni)
	if err != nil {
		return nil, fmt.Errorf("could not parse sni get response, error: %v", err)
	}

	if sni.Name == "" {
		return nil, nil
	}

	return sni, nil
}

func (snisClient *SnisClient) List() (*Snis, error) {

	_, body, errs := newGet(snisClient.config, snisClient.config.HostAddress+SnisPath).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get snis, error: %v", errs)
	}

	snis := &Snis{}
	err := json.Unmarshal([]byte(body), snis)
	if err != nil {
		return nil, fmt.Errorf("could not parse snis list response, error: %v", err)
	}

	return snis, nil
}

func (snisClient *SnisClient) DeleteByName(name string) error {

	res, _, errs := newDelete(snisClient.config, snisClient.config.HostAddress+SnisPath+name).End()
	if errs != nil {
		return fmt.Errorf("could not delete sni, result: %v error: %v", res, errs)
	}

	return nil
}

func (snisClient *SnisClient) UpdateByName(name string, snisRequest *SnisRequest) (*Sni, error) {

	_, body, errs := newPatch(snisClient.config, snisClient.config.HostAddress+SnisPath+name).Send(snisRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update sni, error: %v", errs)
	}

	updatedSni := &Sni{}
	err := json.Unmarshal([]byte(body), updatedSni)
	if err != nil {
		return nil, fmt.Errorf("could not parse sni update response, error: %v", err)
	}

	if updatedSni.SslCertificateId == "" {
		return nil, fmt.Errorf("could not update sni, error: %v", body)
	}

	return updatedSni, nil
}
