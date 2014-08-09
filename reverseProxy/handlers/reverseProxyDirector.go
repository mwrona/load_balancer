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
			req.URL.Host = context.ProxyAddress + ":" + context.ProxyPort
			log.Printf("Reverse Proxy : error, empty server list, redirect to %v", req.URL)
		} else {
			req.URL.Scheme = "http"
			req.URL.Host = context.ServersList.GetNext()
			log.Printf("Reverse Proxy : redirect to %v", req.URL)
		}
	}
}
