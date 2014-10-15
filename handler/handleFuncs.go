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

func ServicesManagment(f func(string, *model.ServicesList, http.ResponseWriter, *http.Request)) contextHandlerFunction {
	return func(context *model.Context, w http.ResponseWriter, r *http.Request) error {
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
		f(address, sl, w, r)
		return nil
	}
}

func Registration(address string, sl *model.ServicesList, w http.ResponseWriter, r *http.Request) {
	if err := sl.AddService(address); err == nil {
		messageWriter(r.URL.String(), "Registered new "+sl.Name()+": "+address, w)
	} else {
		messageWriter(r.URL.String(), err.Error(), w)
	}
}

func Unregistration(address string, sl *model.ServicesList, w http.ResponseWriter, r *http.Request) {
	sl.UnregisterService(address)
	messageWriter(r.URL.String(), "Unegistered "+sl.Name()+": "+address, w)
}

func printServicesList(sl *model.ServicesList, w http.ResponseWriter) {
	fmt.Fprintln(w, sl.Name()+":\n")
	for _, val := range sl.GetServicesList() {
		fmt.Fprintln(w, "\t", val)
	}
}

func printAllServicesList(slt model.SerivesesListMap, w http.ResponseWriter) {
	for _, sl := range slt {
		printServicesList(sl, w)
		fmt.Fprintln(w)
	}
}

func List(context *model.Context, w http.ResponseWriter, r *http.Request) error {
	service_name := r.FormValue("name")
	if service_name == "" {
		printAllServicesList(context.ServicesTypesList, w)
		log.Printf("%s\nMessage: all services list\n\n", r.URL.String())
		return nil
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		return model.NewHTTPError("Service "+service_name+" does not exist", 412)
	}
	log.Printf("%s\nMessage: %s list\n\n", r.URL.String(), sl.Name())

	printServicesList(sl, w)

	return nil
}

func RedirectionError(w http.ResponseWriter, req *http.Request) {
	message := req.FormValue("message")
	if message != "" {
		http.Error(w, message, 404)
	} else {
		http.Error(w, "Service list is empty or all services are not responding.", 404)
	}
}
