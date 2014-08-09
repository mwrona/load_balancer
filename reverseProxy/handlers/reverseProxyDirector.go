package handlers

import (
	"load_balancer/reverseProxy/env"
	"load_balancer/reverseProxy/model"
	"log"
	"net/http"
)

func ReverseProxyDirector(context *model.Context) func(*http.Request) {
	return func(req *http.Request) {
		if context.ServersList.IsEmpty() {
			req.URL.Scheme = env.Protocol
			req.URL.Path = "/error/"
			req.URL.RawQuery = ""
			req.URL.Host = context.ProxyAddress + ":" + context.ProxyPort
			log.Printf("Reverse Proxy : error, empty server list")
		} else {
			req.URL.Scheme = "http"
			req.URL.Host = context.ServersList.GetNext()
			req.URL.Path = "/get/"
			req.URL.RawQuery = ""
			log.Printf("Reverse Proxy : redirect to " + req.URL.Host + req.URL.Path)
		}

	}
}
