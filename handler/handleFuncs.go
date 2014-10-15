package handler

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

func messageWriter(query, message string, w http.ResponseWriter) {
	log.Printf("%s\nResponse: %s\n\n", query, message)
	fmt.Fprintf(w, message)
}

func Registration(context *model.Context, w http.ResponseWriter, r *http.Request) error {
	address := r.FormValue("address")
	service_name := r.FormValue("name")
	if address == "" {
		return model.NewHTTPError("Missing address", 412)
	}
	if service_name == "" {
		return model.NewHTTPError("Missing service name", 412)
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		return model.NewHTTPError("Service "+service_name+" does not exist", 412)
	}

	if err := sl.AddService(address); err == nil {
		messageWriter(r.URL.String(), "Registered new "+sl.Name()+": "+address, w)
		context.StateChan <- 's'
	} else {
		messageWriter(r.URL.String(), err.Error(), w)
	}
	return nil
}

/*
func UnregistrationHandler(servicesTypesList map[string]*model.ServicesList, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("UnregistrationHandler: error, missing address\n\n")
	} else if service_name == "" {
		fmt.Fprintf(w, "Error: missing service name")
		log.Printf("RegistrationHandler: error, missing missing service name\n\n")
	} else {
		servicesTypesList.UnregisterService(address)
		fmt.Fprintf(w, "Unregistered %s:  %s", servicesTypesList.Name(), address)
		log.Printf("UnregistrationHandler: unregistered " + servicesTypesList.Name() + ": " + address + "\n\n")
	}
}
*/
func List(context *model.Context, w http.ResponseWriter, r *http.Request) error {
	service_name := r.FormValue("name")
	if service_name == "" {
		return model.NewHTTPError("Missing service name", 422)
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		return model.NewHTTPError("Service "+service_name+" does not exist", 422)
	}

	log.Printf("%s\nMessage: %s list\n\n", r.URL.String(), sl.Name())

	fmt.Fprintln(w, "Available "+sl.Name()+":\n")
	for _, val := range sl.GetServicesList() {
		fmt.Fprintln(w, val)
	}

	return nil
}

func RedirectionError(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Service list is empty or all services are not responding. Please try again later.", 404)
}
