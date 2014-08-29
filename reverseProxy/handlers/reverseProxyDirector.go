package handlers

import (
	"fmt"
	"load_balancer/reverseProxy/env"
	"load_balancer/reverseProxy/model"
	"log"
	"net/http"
	"strings"
)

func ReverseProxyDirector(context *model.Context) func(*http.Request) {
	return func(req *http.Request) {
		fmt.Println()
		log.Printf("Reverse Proxy : query: %v", req.URL.RequestURI())

		splited := strings.SplitN(req.URL.RequestURI(), "/", 3)
		fmt.Println(splited)
		if len(splited) < 3 {
			splited = append(splited, "")
		}
		if splited[1] == "storage" {
			req.URL.Scheme = "http"
			req.Header.Add("X-Forwarded-Proto", env.Protocol)
			req.URL.Host = "localhost:20000"
			req.URL.Path = "/" + splited[2]
			log.Printf("Reverse Proxy : redirect to %v\n\n", req.URL)
		} else if splited[1] == "information" {
			req.URL.Scheme = "http"
			req.Header.Add("X-Forwarded-Proto", env.Protocol)
			req.URL.Host = "localhost:11300"
			req.URL.Path = "/" + splited[2]
			log.Printf("Reverse Proxy : redirect to %v\n\n", req.URL)

		} else if host, err := context.ExperimentManagersList.GetNext(); err != nil {
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
