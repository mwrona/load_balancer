package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"scalarm_load_balancer/model"
	"strings"
)

func redirectToError(context *model.Context, req *http.Request, err error) {
	log.Printf("%v\nUnable to redirect: %v", req.URL.RequestURI(), err.Error())

	values := url.Values{}
	values.Add("message", err.Error())

	req.URL.RawQuery = values.Encode()
	req.URL.Scheme = context.LoadBalancerScheme
	req.URL.Host = req.Host
	req.URL.Path = "/error/"
}

func parseURL(context *model.Context, req *http.Request) (string, *model.ServicesList) {
	splitted := strings.SplitN(req.URL.Path, "/", 3)
	if len(splitted) < 3 {
		splitted = append(splitted, "")
	}

	prefix := "/" + splitted[1]
	sl := context.RedirectionsList[prefix]
	path := "/" + splitted[2]

	if sl == nil {
		sl = context.RedirectionsList["/"]
		path = req.URL.Path
	}

	return path, sl
}

func ReverseProxyDirector(context *model.Context) func(*http.Request) {
	return func(req *http.Request) {
		oldURL := req.URL.RequestURI()

		req.Header.Add("X-Forwarded-Proto", context.LoadBalancerScheme)

		path, servicesList := parseURL(context, req)
		if servicesList == nil {
			redirectToError(context, req, fmt.Errorf("Requested redirection does not exists"))
			return
		}

		host, err := servicesList.GetNext()
		if err != nil {
			redirectToError(context, req, err)
			return
		}

		req.URL.Scheme = servicesList.Scheme()
		req.URL.Host = host
		req.URL.Path = path

		fmt.Println()
		log.Printf("%v \nredirect to %v\n\n", oldURL, req.URL)
	}
}
