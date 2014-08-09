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
		log.Printf("Reverse Proxy : error, missing port")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : error, missing address")
	} else {
		context.ServersList.AddServer(address, port)
		fmt.Fprintf(w, "Registered server:  %s", address+":"+port)
		log.Printf("Reverse Proxy : registered server: " + address + ":" + port)
	}
}

func UnregisterHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	address := r.FormValue("address")
	if port == "" {
		fmt.Fprintf(w, "Error: missing port")
		log.Printf("Reverse Proxy : error, missing port")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("Reverse Proxy : error, missing address")
	} else {
		context.ServersList.UnregisterServer(address, port)
		fmt.Fprintf(w, "Unregistered server:  %s", address+":"+port)
		log.Printf("Reverse Proxy : unregistered server: " + address + ":" + port)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server list is empty!")
}

func ListHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	log.Printf("Reverse Proxy : printing servers list")
	fmt.Fprintf(w, "Servers available:\n")
	for _, val := range context.ServersList.GetServersList() {
		fmt.Fprintf(w, val+"\n")
	}
}
