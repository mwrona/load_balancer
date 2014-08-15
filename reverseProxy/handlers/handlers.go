package handlers

import (
	"fmt"
	"load_balancer/reverseProxy/model"
	"log"
	"net/http"
)

func RegisterHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	address := r.FormValue("address")
	if port == "" {
		fmt.Fprintf(w, "Error: missing port")
		log.Printf("Reverse Proxy : registered server: error, missing port\n\n")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : registered server: error, missing addres\n\n")
	} else {
		if err := context.ServersList.AddServer(address, port); err == nil {
			fmt.Fprintf(w, "Registered server:  %s", address+":"+port)
			log.Printf("Reverse Proxy : registered server: " + address + ":" + port + "\n\n")
		} else {
			fmt.Fprintf(w, "Host already exists")
			log.Printf("Reverse Proxy : %v \n\n", err)
		}
	}
}

func UnregisterHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	address := r.FormValue("address")
	if port == "" {
		fmt.Fprintf(w, "Error: missing port")
		log.Printf("Reverse Proxy : error, missing port\n\n")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : error, missing address\n\n")
	} else {
		context.ServersList.UnregisterServer(address, port)
		fmt.Fprintf(w, "Unregistered server:  %s", address+":"+port)
		log.Printf("Reverse Proxy : unregistered server: " + address + ":" + port + "\n\n")
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server list is empty or all servers are not responding.")
}

func ListHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	log.Printf("Reverse Proxy : printing servers list\n\n")
	fmt.Fprintln(w, "Servers available:\n")
	for _, val := range context.ServersList.GetServersList() {
		fmt.Fprintln(w, val)
	}
}
