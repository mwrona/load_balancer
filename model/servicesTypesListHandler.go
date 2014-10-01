package model

import "net/http"

type servicesTypesListHandlerFunction func(map[string]*ServicesList, http.ResponseWriter, *http.Request)

type servicesTypesListHandler struct {
	servicesTypesList map[string]*ServicesList
	f                 servicesTypesListHandlerFunction
}

func (stlh servicesTypesListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stlh.f(stlh.servicesTypesList, w, r)
}

func ServicesTypesListHandler(servicesList map[string]*ServicesList,
	f servicesTypesListHandlerFunction) servicesTypesListHandler {
	return servicesTypesListHandler{servicesList, f}
}
