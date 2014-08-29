package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"scalarm_load_balancer/handlers"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/services"
	"scalarm_load_balancer/utils"
)

func main() {
	var configFile string
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		configFile = "config.json"
	}
	config, err := model.LoadConfig(configFile)
	if err != nil {
		fmt.Println("An error occurred while loading configuration: " + configFile)
		fmt.Println(err.Error())
		return
	}

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
		ExperimentManagersList: model.NewServicesList("http"),
		StorageManagersList:    model.NewServicesList("http"),
		LoadBalancerAddress:    config.LoadBalancerAddress,
		LoadBalancerScheme:     config.LoadBalancerScheme,
	}

	if _, err := utils.RepetitiveCaller(
		func() (interface{}, error) {
			return nil, utils.InformationServiseRegistration(config.LoadBalancerAddress,
				config.InformationServiceAddress,
				config.InformationServiceScheme)
		}, nil, "InformationServiseRegistration"); err != nil {
		log.Printf("Registration to Information Service failed")
		return
	}

	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context),
		Transport: TransportCert}
	http.Handle("/", reverseProxy)

	http.Handle("experiment_managers/register", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.RegistrationHandler))
	http.Handle("experiment_managers/unregister", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.UnregistrationHandler))
	http.Handle("experiment_managers/list", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.ListHandler))

	http.Handle("storage_managers/register", model.ServicesListHandler(context.StorageManagersList,
		handlers.RegistrationHandler))
	http.Handle("storage_managers/unregister", model.ServicesListHandler(context.StorageManagersList,
		handlers.UnregistrationHandler))
	http.Handle("storage_managers/list", model.ServicesListHandler(context.StorageManagersList,
		handlers.ListHandler))

	http.HandleFunc("/error/", handlers.ErrorHandler)

	go services.StartMulticastAddressSender(config.LoadBalancerAddress, config.MulticastAddress)
	go services.ServicesStatusChecker(context.ExperimentManagersList)
	go services.ServicesStatusChecker(context.StorageManagersList)

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
	if err != nil {
		fmt.Println("An error occurred while running service on port " + config.Port)
		fmt.Println(err.Error())
	}
}
