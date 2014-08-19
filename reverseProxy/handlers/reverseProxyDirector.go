package handlers

import (
	"fmt"
	"load_balancer/reverseProxy/env"
	"load_balancer/reverseProxy/model"
	"log"
	"net/http"
)

func ReverseProxyDirector(context *model.Context) func(*http.Request) {
	return func(req *http.Request) {
		fmt.Println()
		log.Printf("Reverse Proxy : query: %v", req.URL)
		if host, err := context.ServersList.GetNext(); err != nil {
			req.URL.Scheme = env.Protocol
			req.URL.Path = "/error/"
			req.URL.Host = context.ProxyAddress + ":" + context.ProxyPort
			log.Printf("Reverse Proxy : error, redirect to %v\n\n", req.URL)
		} else {
			req.URL.Scheme = "http"
			req.Header.Add("X-Forwarded-Proto", env.Protocol)
			req.URL.Host = host
			log.Printf("Reverse Proxy : redirect to %v\n\n", req.URL)
		}
	}
}
