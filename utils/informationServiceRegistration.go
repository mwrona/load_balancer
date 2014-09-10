package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func InformationServiceRegistration(loadBalancerAddress, informationServiceAddress, informationServiceScheme,
	informationServiceUser, informationServicePass string) error {

	log.Printf("Registration to Load Balancer on: " + informationServiceAddress)
	data := url.Values{"address": {loadBalancerAddress}}
	request, err := http.NewRequest("POST", informationServiceScheme+"://"+informationServiceAddress+"/experiment_managers", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(informationServiceUser, informationServicePass)

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
