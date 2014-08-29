package handlers

import (
	"fmt"
	"log"
	"net/http"
	"scalarm_load_balancer/model"
	"strings"
)

func ReverseProxyDirector(context *model.Context) func(*http.Request) {
	return func(req *http.Request) {
		fmt.Println()
		log.Printf("ReverseProxyDirector : query: %v", req.URL.RequestURI())

		splited := strings.SplitN(req.URL.Path, "/", 3)
		if len(splited) < 3 {
			splited = append(splited, "")
		}

		req.Header.Add("X-Forwarded-Proto", context.LoadBalancerScheme)

		if splited[1] == "information" {
			req.URL.Scheme = context.InformationServiceScheme
			req.URL.Host = context.InformationServiceAddress
			req.URL.Path = "/" + splited[2]
		} else {
			var servicesList *model.ServicesList
			var path string

			if splited[1] == "storage" {
				servicesList = context.StorageManagersList
				path = "/" + splited[2]
			} else {
				servicesList = context.ExperimentManagersList
				path = req.URL.Path
			}

			if host, err := servicesList.GetNext(); err != nil {
				req.URL.Scheme = context.LoadBalancerScheme
				req.URL.Host = context.LoadBalancerAddress
				req.URL.Path = "/error/"
			} else {
				req.URL.Scheme = "http"
				req.URL.Host = host
				req.URL.Path = path
			}
		}

		log.Printf("ReverseProxyDirector : redirect to %v\n\n", req.URL)
	}
}
