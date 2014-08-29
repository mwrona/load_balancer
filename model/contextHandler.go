package model

import "net/http"

type contextHandlerFunction func(*Context, http.ResponseWriter, *http.Request)

type contextHandler struct {
	context *Context
	f       contextHandlerFunction
}

func (ch contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.f(ch.context, w, r)
}

func ContextHandler(context *Context, f contextHandlerFunction) contextHandler {
	return contextHandler{context, f}
}
