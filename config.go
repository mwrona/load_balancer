package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mwrona/scalarm_load_balancer/services"
)

type Config struct {
	Port                       string
	MulticastAddress           string
	PrivateLoadBalancerAddress string
	LoadBalancerScheme         string
	CertFilePath               string
	KeyFilePath                string
	LogDirectory               string
	StateDirectory             string
	RedirectionConfig          []services.RedirectionPolicy
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

	if config.MulticastAddress == "" {
		return nil, fmt.Errorf("Multicast address is missing")
	}
	if config.PrivateLoadBalancerAddress == "" {
		config.PrivateLoadBalancerAddress = "localhost"
	}
	if config.LoadBalancerScheme == "" {
		config.LoadBalancerScheme = "https"
	}
	if config.Port == "" {
		if config.LoadBalancerScheme == "https" {
			config.Port = "443"
		} else if config.LoadBalancerScheme == "http" {
			config.Port = "80"
		}
	}
	if config.LoadBalancerScheme != "https" && config.LoadBalancerScheme != "http" {
		return nil, fmt.Errorf("Unsuported protocol in LoadBalancerScheme")
	}
	if config.CertFilePath == "" {
		config.CertFilePath = "cert.pem"
	}
	if config.KeyFilePath == "" {
		config.KeyFilePath = "key.pem"
	}
	if config.RedirectionConfig == nil {
		return nil, fmt.Errorf("RedirectionConfig is missing")
	}
	if config.LogDirectory == "" {
		config.LogDirectory = "log"
	}
	if !strings.HasSuffix(config.StateDirectory, "/") && config.StateDirectory != "" {
		config.StateDirectory += "/"
	}

	return config, nil
}
