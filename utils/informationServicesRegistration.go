package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func InformationSeriviseRegistration(address, port string) error {
	//TODO
	log.Printf("Registration to Load Balancer on: " + address + ":" + port)
	data := url.Values{"address": {address + ":" + port}}
	request, err := http.NewRequest("POST", "http://localhost:11300/experiment_managers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("scalarm", "scalarm")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf(string(body))
	return nil
}
