package model

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	PublicLoadBalancerAddress  string
	Port                       string
	MulticastAddress           string
	PrivateLoadBalancerAddress string
	LoadBalancerScheme         string
	InformationServiceAddress  string
	InformationServiceUser     string
	InformationServicePass     string
	CertFilePath               string
	KeyFilePath                string
	RedirectionConfig          []RedirectionPolicy
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

	if config.LoadBalancerScheme == "" {
		config.LoadBalancerScheme = "https"
	}
	if config.PrivateLoadBalancerAddress == "" {
		config.PrivateLoadBalancerAddress = "localhost"
	}

	return config, nil
}
