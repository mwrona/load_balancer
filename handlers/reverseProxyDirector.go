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

		splitted := strings.SplitN(req.URL.Path, "/", 3)
		if len(splitted) < 3 {
			splitted = append(splitted, "")
		}

		req.Header.Add("X-Forwarded-Proto", context.LoadBalancerScheme)

		if splitted[1] == "information" {
			req.URL.Scheme = context.InformationServiceScheme
			req.URL.Host = context.InformationServiceAddress
			req.URL.Path = "/" + splitted[2]
		} else {
			var servicesList *model.ServicesList
			var path string

			if splitted[1] == "storage" {
				servicesList = context.StorageManagersList
				path = "/" + splitted[2]
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
