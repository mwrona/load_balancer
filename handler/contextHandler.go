package handler

import (
	"log"
	"net/http"
	"scalarm_load_balancer/model"
)

type contextHandlerFunction func(*model.Context, http.ResponseWriter, *http.Request) error

type contextHandler struct {
	context *model.Context
	f       contextHandlerFunction
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, _ := ch.f(ch.context, w, r).(*model.HTTPError)
	if err != nil {
		log.Printf("%s\nResponse: %v; %s\n\n", r.URL.RequestURI(), err.Code(), err.Error())
		http.Error(w, err.Error(), err.Code())
	}
}

func Context(context *model.Context, f contextHandlerFunction) http.Handler {
	return contextHandler{context, f}
}
