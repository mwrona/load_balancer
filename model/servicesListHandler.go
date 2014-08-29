package model

import "net/http"

type servicesListHandlerFunction func(*ServicesList, http.ResponseWriter, *http.Request)

type servicesListHandler struct {
	servicesList *ServicesList
	f            servicesListHandlerFunction
}

func (slh servicesListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slh.f(slh.servicesList, w, r)
}

func ServicesListHandler(servicesList *ServicesList, f servicesListHandlerFunction) servicesListHandler {
	return servicesListHandler{servicesList, f}
}
