package handlers

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

func RegistrationHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	service_name := r.FormValue("name")
	if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		log.Printf("RegistrationHandler: error, missing addres\n\n")
		return
	}
	if service_name == "" {
		fmt.Fprintf(w, "Error: missing service name")
		log.Printf("RegistrationHandler: error, missing missing service name\n\n")
		return
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		fmt.Fprintf(w, "Service "+service_name+" does not exist")
		log.Printf("RegistrationHandler: Service " + service_name + " does not exist")
		return
	}

	if err := sl.AddService(address); err == nil {
		fmt.Fprintf(w, "Registered %s:  %s", sl.Name(), address)
		log.Printf("RegistrationHandler: registered " + sl.Name() + ": " + address + "\n\n")
		context.StateChan <- 's'
	} else {
		fmt.Fprintf(w, "Host already exists")
		log.Printf("RegistrationHandler %s: %v \n\n", sl.Name(), err)
	}

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
func ListHandler(context *model.Context, w http.ResponseWriter, r *http.Request) {
	service_name := r.FormValue("name")
	if service_name == "" {
		fmt.Fprintf(w, "Error: missing service name")
		log.Printf("RegistrationHandler: error, missing missing service name\n\n")
		return
	}

	sl, ok := context.ServicesTypesList[service_name]
	if ok == false {
		fmt.Fprintf(w, "Service "+service_name+" does not exist")
		log.Printf("ListHandler: Service " + service_name + " does not exist")
		return
	}

	log.Printf("ListHandler: printing services list\n\n")
	fmt.Fprintln(w, "Available "+sl.Name()+":\n")
	for _, val := range sl.GetServicesList() {
		fmt.Fprintln(w, val)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service list is empty or all services are not responding. Please try again later.")
}
