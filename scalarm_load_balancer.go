package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"scalarm_load_balancer/handlers"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/services"
	"scalarm_load_balancer/utils"
)

func main() {
	config, err := model.LoadConfig("config.json")

	var TLSClientConfigCert *tls.Config
	var TransportCert *http.Transport

	if config.CertificateCheckDisable {
		log.Printf("Reverse Proxy : Certificates checking disabled")
		TLSClientConfigCert = &tls.Config{InsecureSkipVerify: true}
		TransportCert = &http.Transport{
			TLSClientConfig: TLSClientConfigCert,
		}
	} else {
		TLSClientConfigCert = &tls.Config{}
		TransportCert = &http.Transport{}
	}

	context := &model.Context{
		ExperimentManagersList: model.NewExperimentManagersList(),
		LoadBalancerAddress:    config.LoadBalancerAddress,
		LoadBalancerScheme:     config.LoadBalancerScheme,
	}

	if _, err := utils.RepititveCaller(
		func() (interface{}, error) {
			return nil, utils.InformationServiseRegistration(config.LoadBalancerAddress,
				config.InformationServiceAddress,
				config.InformationServiceScheme)
		}, nil); err != nil {
		log.Printf("Registration to Information Service failed")
		return
	}

	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context), Transport: TransportCert}
	http.Handle("/", reverseProxy)

	http.Handle("/register", model.ContextHandler(context, handlers.RegisterHandler))
	http.Handle("/unregister", model.ContextHandler(context, handlers.UnregisterHandler))
	http.Handle("/list", model.ContextHandler(context, handlers.ListHandler))
	http.HandleFunc("/error/", handlers.ErrorHandler)

	go services.MulticastAddressSender(config.LoadBalancerAddress, config.MulticastAddress)
	go services.ExperimentManagersStatusChecker(context.ExperimentManagersList)

	log.Printf("Reverse Proxy : Start")

	server := &http.Server{
		Addr:      ":" + config.Port,
		TLSConfig: TLSClientConfigCert,
	}

	if config.LoadBalancerScheme == "http" {
		err = server.ListenAndServe()
	} else { // "https"
		err = server.ListenAndServeTLS(config.CertFilePath, config.KeyFilePath)
	}
	utils.Check(err)
}
