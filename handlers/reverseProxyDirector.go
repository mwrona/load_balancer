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

		prefix := "/" + splitted[1]
		sl, ok := context.RedirectionsList[prefix]
		path := "/" + splitted[2]

		if ok == false {
			sl, ok = context.RedirectionsList["/"]
			path = req.URL.Path
		}

		if ok {
			if host, err := sl.GetNext(); err == nil {
				req.URL.Scheme = sl.Scheme()
				req.URL.Host = host
				req.URL.Path = path
			} else {
				ok = false
			}
		}

		if !ok {
			req.URL.Scheme = context.LoadBalancerScheme
			req.URL.Host = context.LoadBalancerAddress
			req.URL.Path = "/error/"
		}

		log.Printf("ReverseProxyDirector : redirect to %v\n\n", req.URL)
	}
}
