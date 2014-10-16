package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var stateChan chan byte
var stateDaemonInit bool

func init() {
	stateChan = make(chan byte, 100)
	stateDaemonInit = false
}

type State struct {
	Name          string
	Scheme        string
	AddressesList []string
}

func StartStateDaemon(services TypesMap) {
	if stateDaemonInit == false {
		go stateDaemon(services)
		stateDaemonInit = true
	}
}

func stateDaemon(services TypesMap) {
	var s byte
	for {
		select {
		case s = <-stateChan:
			if s == 'l' {
				for s != 'e' {
					select {
					case s = <-stateChan:
					}
				}
			} else {
				SaveState(services)
			}
		}
	}
}

func SaveState(services TypesMap) {
	statesList := make([]State, 0, 0)
	for _, s := range services {
		statesList = append(statesList, State{s.Name(), s.Scheme(), s.AddressesList()})
	}
	data, err := json.Marshal(statesList)
	if err != nil {
		log.Printf("An error occurred while saving state: %s", err.Error())
		return
	}
	err = ioutil.WriteFile("state.json", data, 0644)
	if err != nil {
		log.Printf("An error occurred while saving state: %s", err.Error())
		return
	}
	log.Println("State saved succesfully")
}

func LoadState(services TypesMap) {
	stateChan <- 'l'
	defer func() {
		stateChan <- 'e'
	}()
	data, err := ioutil.ReadFile("state.json")
	if err != nil {
		log.Printf("An error occurred while loading state: %s", err.Error())
		return
	}
	var statesList []State
	err = json.Unmarshal(data, &statesList)
	if err != nil {
		log.Printf("An error occurred while loading state: %s", err.Error())
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
