package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"scalarm_load_balancer/env"
	"scalarm_load_balancer/handlers"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/services"
	"scalarm_load_balancer/utils"
)

func main() {
	if env.CertOff {
		log.Printf("Reverse Proxy : Certificates checking disabled")
	}

	context := &model.Context{
		ExperimentManagersList: model.NewExperimentManagersList(),
		ProxyPort:              "9000",
		ProxyAddress:           "localhost"}

	if _, err := utils.RepititveCaller(
		func() (interface{}, error) {
			return nil, utils.InformationSeriviseRegistration(context.ProxyAddress, context.ProxyPort)
		}, nil); err != nil {
		log.Printf("Registration to Information Service failed")
		return
	}

	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context), Transport: env.TransportCert}
	http.Handle("/", reverseProxy)

	http.Handle("/register", model.ContextHandler(context, handlers.RegisterHandler))
	http.Handle("/unregister", model.ContextHandler(context, handlers.UnregisterHandler))
	http.Handle("/list", model.ContextHandler(context, handlers.ListHandler))
	http.HandleFunc("/error/", handlers.ErrorHandler)

	go services.MulticastAddressSender(context.ProxyAddress, context.ProxyPort)
	go services.ExperimentManagersStatusChecker(context.ExperimentManagersList)

	log.Printf("Reverse Proxy : Start")

	server := &http.Server{
		Addr:      ":" + context.ProxyPort,
		TLSConfig: env.TLSClientConfigCert,
	}

	//err := server.ListenAndServe()
	//err := server.ListenAndServeTLS("cert.pem", "key.pem")
	err := env.StartServer(server, "cert.pem", "key.pem")
	utils.Check(err)
}
