package model

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	RemoteLoadBalancerAddress string
	Port                      string
	MulticastAddress          string
	LocalLoadBalancerAddress  string
	LoadBalancerScheme        string
	CertificateCheckDisable   bool
	InformationServiceAddress string
	InformationServiceScheme  string
	InformationServiceUser    string
	InformationServicePass    string
	CertFilePath              string
	KeyFilePath               string
}

func LoadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
