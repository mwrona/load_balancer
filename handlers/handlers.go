package handlers

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

func RegistrationHandler(servicesList *model.ServicesList, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("RegistrationHandler: error, missing addres\n\n")
	} else {
		if err := servicesList.AddService(address); err == nil {
			fmt.Fprintf(w, "Registered %s:  %s", servicesList.Name(), address)
			log.Printf("RegistrationHandler: registered " + servicesList.Name() + ": " + address + "\n\n")
		} else {
			fmt.Fprintf(w, "Host already exists")
			log.Printf("RegistrationHandler %s: %v \n\n", servicesList.Name(), err)
		}
	}
}

func UnregistrationHandler(servicesList *model.ServicesList, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("UnregistrationHandler: error, missing address\n\n")
	} else {
		servicesList.UnregisterService(address)
		fmt.Fprintf(w, "Unregistered %s:  %s", servicesList.Name(), address)
		log.Printf("UnregistrationHandler: unregistered " + servicesList.Name() + ": " + address + "\n\n")
	}
}

func ListHandler(servicesList *model.ServicesList, w http.ResponseWriter, r *http.Request) {
	log.Printf("ListHandler: printing services list\n\n")
	fmt.Fprintln(w, "Available "+servicesList.Name()+":\n")
	for _, val := range servicesList.GetServicesList() {
		fmt.Fprintln(w, val)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service list is empty or all services are not responding. Please try again later.")
}
