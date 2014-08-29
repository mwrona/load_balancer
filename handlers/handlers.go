package handlers

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

func RegisterHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : registered server: error, missing addres\n\n")
	} else {
		if err := context.ExperimentManagersList.AddServer(address); err == nil {
			fmt.Fprintf(w, "Registered server:  %s", address)
			log.Printf("Reverse Proxy : registered server: " + address + "\n\n")
		} else {
			fmt.Fprintf(w, "Host already exists")
			log.Printf("Reverse Proxy : %v \n\n", err)
		}
	}
}

func UnregisterHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : error, missing address\n\n")
	} else {
		context.ExperimentManagersList.UnregisterServer(address)
		fmt.Fprintf(w, "Unregistered server:  %s", address)
		log.Printf("Reverse Proxy : unregistered server: " + address + "\n\n")
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server list is empty or all servers are not responding.")
}

func ListHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	log.Printf("Reverse Proxy : printing servers list\n\n")
	fmt.Fprintln(w, "Servers available:\n")
	for _, val := range context.ExperimentManagersList.GetExperimentManagersList() {
		fmt.Fprintln(w, val)
	}
}
