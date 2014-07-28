package utils

import (
	"io/ioutil"
)

type Config struct {
    Address string
}

var filename = "config.txt"

func LoadConfig() (*Config, error) {
    address, err := ioutil.ReadFile(filename)
	Check(err)
	
    return &Config{Address: string(address[:])}, nil
}
