package handler

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

func messageWriter(who, query, message string, w http.ResponseWriter) {
	log.Printf("%s: \nQuery: %s\nResponse: %s\n\n", who, query, message)
	fmt.Fprintf(w, message)
}

func Registration(context *model.Context, w http.ResponseWriter, r *http.Request) error {
	address := r.FormValue("address")
	service_name := r.FormValue("name")
	if address == "" {
		return model.NewHTTPError("RegistrationHandler", r.URL.String(), "Missing address", 422)
	}
	if service_name == "" {
		return model.NewHTTPError("RegistrationHandler", r.URL.String(), "Missing service name", 422)
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		return model.NewHTTPError("RegistrationHandler", r.URL.String(), "Service "+service_name+" does not exist", 422)
	}

	if err := sl.AddService(address); err == nil {
		messageWriter("RegistrationHandler", r.URL.String(), "Registered new "+sl.Name()+": "+address, w)
		context.StateChan <- 's'
	} else {
		messageWriter("RegistrationHandler", r.URL.String(), err.Error(), w)
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
		return model.NewHTTPError("ListHandler", r.URL.String(), "Missing service name", 422)
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		return model.NewHTTPError("ListHandler", r.URL.String(), "Service "+service_name+" does not exist", 422)
	}

	log.Printf("ListHandler:\nQuery %s\nMessage: %s list\n\n", sl.Name(), r.URL.String())

	fmt.Fprintln(w, "Available "+sl.Name()+":\n")
	for _, val := range sl.GetServicesList() {
		fmt.Fprintln(w, val)
	}

	return nil
}

func Error(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Service list is empty or all services are not responding. Please try again later.", 404)
}
