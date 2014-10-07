package model

import (
	"log"
	"net/http"
)

type contextHandlerFunction func(*Context, http.ResponseWriter, *http.Request) error

type contextHandler struct {
	context *Context
	f       contextHandlerFunction
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err, _ := ch.f(ch.context, w, r).(*HTTPError)
	if err != nil {
		log.Printf("%s: \nQuery: %s\nResponse: %v; %s\n\n", err.Who(), err.Query(), err.Code(), err.Error())
		http.Error(w, err.Error(), err.Code())
	}
}

func ContextHandler(context *Context, f contextHandlerFunction) contextHandler {
	return contextHandler{context, f}
}
