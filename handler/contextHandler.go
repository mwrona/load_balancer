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

type AppContext struct {
	RedirectionsList   services.TypesMap
	ServicesTypesList  services.TypesMap
	LoadBalancerScheme string
}

type contextHandlerFunction func(*AppContext, http.ResponseWriter, *http.Request) error

type contextHandler struct {
	context *AppContext
	f       contextHandlerFunction
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, _ := ch.f(ch.context, w, r).(*hTTPError)
	if err != nil {
		log.Printf("%s\nResponse: %v; %s\n\n", r.URL.RequestURI(), err.Code(), err.Error())
		http.Error(w, err.Error(), err.Code())
	}
}

func Context(context *AppContext, f contextHandlerFunction) http.Handler {
	return contextHandler{context, f}
}
