package model

import (
	"io/ioutil"
	"scalarm_load_balancer/utils"
)

type Config struct {
	Address string
}

const filename = "config.txt"

func LoadConfig() (*Config, error) {
	address, err := ioutil.ReadFile(filename)
	utils.Check(err)

	return &Config{Address: string(address[:])}, nil
}
