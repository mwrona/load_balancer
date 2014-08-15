package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func InformationSeriviseRegistration(address, port string) {
	//TODO
	data := url.Values{"address": {address + ":" + port}}
	request, err := http.NewRequest("POST", "http://localhost:11300/experiment_managers", strings.NewReader(data.Encode()))
	Check(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("scalarm", "scalarm")

	client := http.Client{}
	resp, err := client.Do(request)
	Check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	Check(err)
	log.Printf(string(body))

}
