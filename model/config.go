package model

import (
	"encoding/json"
	"io/ioutil"
	"scalarm_load_balancer/utils"
)

type Config struct {
	LoadBalancerAddress       string
	Port                      string
	MulticastAddress          string
	LoadBalancerScheme        string
	CertificateCheckDisable   bool
	InformationServiceAddress string
	InformationServiceScheme  string
	CertFilePath              string
	KeyFilePath               string
}

func LoadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	utils.Check(err)
	config := &Config{}
	err = json.Unmarshal(file, config)
	utils.Check(err)
	return config, nil
}
