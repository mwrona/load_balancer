package handler

import (
	"log"
	"net/http"
	"scalarm_load_balancer/services"
)

type hTTPError struct {
	message string
	code    int
}

func newHTTPError(message string, code int) *hTTPError {
	return &hTTPError{message, code}
}

func (e *hTTPError) Error() string {
	return e.message
}

func (e *hTTPError) Code() int {
	return e.code
}

//multi thread use - do not modify entries
type appContext struct {
	redirectionsList   services.TypesMap
	servicesTypesList  services.TypesMap
	loadBalancerScheme string
}

func AppContext(redirectionsList, servicesTypesList services.TypesMap, loadBalancerScheme string) *appContext {
	return &appContext{
		redirectionsList:   redirectionsList,
		servicesTypesList:  servicesTypesList,
		loadBalancerScheme: loadBalancerScheme,
	}
}

type contextHandlerFunction func(*appContext, http.ResponseWriter, *http.Request) error

type contextHandler struct {
	context *appContext
	f       contextHandlerFunction
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, _ := ch.f(ch.context, w, r).(*hTTPError)
	if err != nil {
		log.Printf("%s\nResponse: %v; %s\n\n", r.URL.RequestURI(), err.Code(), err.Error())
		http.Error(w, err.Error(), err.Code())
	}
}

func Context(context *appContext, f contextHandlerFunction) http.Handler {
	return contextHandler{context, f}
}
