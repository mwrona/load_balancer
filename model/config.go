package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Port                       string
	MulticastAddress           string
	PrivateLoadBalancerAddress string
	LoadBalancerScheme         string
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

	if config.Port == "" {
		config.Port = "443"
	}
	if config.MulticastAddress == "" {
		return nil, fmt.Errorf("Multicast address is missing")
	}
	if config.PrivateLoadBalancerAddress == "" {
		config.PrivateLoadBalancerAddress = "localhost"
	}
	if config.LoadBalancerScheme == "" {
		config.LoadBalancerScheme = "https"
	}
	if config.CertFilePath == "" {
		config.CertFilePath = "cert.pem"
	}
	if config.KeyFilePath == "" {
		config.KeyFilePath = "key.pem"
	}
	if config.RedirectionConfig == nil {
		return nil, fmt.Errorf("Redirection Config is missing")
	}

	return config, nil
}
