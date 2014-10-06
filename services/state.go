package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"scalarm_load_balancer/model"
)

type State struct {
	Name          string
	Scheme        string
	AddressesList []string
}

func StateDeamon(stateChan chan byte, services map[string]*model.ServicesList) {
	for {
		select {
		case <-stateChan:
			saveState(services)
		}
	}
}

func saveState(services map[string]*model.ServicesList) {
	statesList := make([]State, 0, 0)
	for _, s := range services {
		statesList = append(statesList, State{s.Name(), s.Scheme(), s.GetServicesList()})
	}
	data, err := json.Marshal(statesList)
	if err != nil {
		log.Printf("An error occurred while saving state: " + err.Error())
		return
	}
	err = ioutil.WriteFile("state.json", data, 0644)
	if err != nil {
		log.Printf("An error occurred while saving state: " + err.Error())
		return
	}
	log.Println("State saved succesfully")
}

func LoadState(services map[string]*model.ServicesList) {
	data, err := ioutil.ReadFile("state.json")
	if err != nil {
		log.Printf("An error occurred while loading state: " + err.Error())
		return
	}
	var statesList []State
	err = json.Unmarshal(data, &statesList)
	if err != nil {
		log.Printf("An error occurred while loading state: " + err.Error())
		return
	}
	for _, state := range statesList {
		if sl, ok := services[state.Name]; ok && sl.Scheme() == state.Scheme {
			for _, address := range state.AddressesList {
				sl.AddService(address)
			}
		} else {
			log.Printf("LoadState: No such service as: %s with scheme %s", state.Name, state.Scheme)
		}
	}
	log.Println("Previous state loaded succesfully")
}
