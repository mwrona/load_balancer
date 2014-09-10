package main

import (
	"crypto/tls"
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
		log.Println("An error occurred while loading configuration: " + configFile)
		log.Println(err.Error())
		return
	}

	if t := os.Getenv("INFORMATION_SERVICE_URL"); t != "" {
		config.InformationServiceAddress = t
	}
	if t := os.Getenv("INFORMATION_SERVICE_LOGIN"); t != "" {
		config.InformationServiceUser = t
	}
	if t := os.Getenv("INFORMATION_SERVICE_PASSWORD"); t != "" {
		config.InformationServicePass = t
	}

	var TLSClientConfigCert *tls.Config
	var TransportCert *http.Transport

	if config.CertificateCheckDisable {
		log.Printf("Certificates checking disabled")
		TLSClientConfigCert = &tls.Config{InsecureSkipVerify: true}
		TransportCert = &http.Transport{
			TLSClientConfig: TLSClientConfigCert,
		}
	} else {
		TLSClientConfigCert = &tls.Config{}
		TransportCert = &http.Transport{}
	}

	context := &model.Context{
		ExperimentManagersList:    model.NewServicesList("http", "ExperimentManager"),
		StorageManagersList:       model.NewServicesList("http", "StorageManager"),
		InformationServiceAddress: config.InformationServiceAddress,
		InformationServiceScheme:  config.InformationServiceScheme,
		LoadBalancerAddress:       config.LocalLoadBalancerAddress,
		LoadBalancerScheme:        config.LoadBalancerScheme,
	}

	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context),
		Transport: TransportCert}
	http.Handle("/", reverseProxy)

	http.Handle("/experiment_managers/register", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.RegistrationHandler))
	http.Handle("/experiment_managers/unregister", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.UnregistrationHandler))
	http.Handle("/experiment_managers", model.ServicesListHandler(context.ExperimentManagersList,
		handlers.ListHandler))

	http.Handle("/storage_managers/register", model.ServicesListHandler(context.StorageManagersList,
		handlers.RegistrationHandler))
	http.Handle("/storage_managers/unregister", model.ServicesListHandler(context.StorageManagersList,
		handlers.UnregistrationHandler))
	http.Handle("/storage_managers", model.ServicesListHandler(context.StorageManagersList,
		handlers.ListHandler))

	http.HandleFunc("/error/", handlers.ErrorHandler)

	if _, err := utils.RepetitiveCaller(
		func() (interface{}, error) {
			return nil, utils.InformationServiceRegistration(config.RemoteLoadBalancerAddress,
				config.InformationServiceAddress,
				config.InformationServiceScheme,
				config.InformationServiceUser,
				config.InformationServicePass)
		},
		nil,
		"InformationServiseRegistration",
	); err != nil {
		log.Printf("Registration to Information Service failed")
		return
	}

	go services.StartMulticastAddressSender(config.LocalLoadBalancerAddress, config.MulticastAddress)
	go services.ServicesStatusChecker(context.ExperimentManagersList)
	go services.ServicesStatusChecker(context.StorageManagersList)

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
		log.Println("An error occurred while running service on port " + config.Port)
		log.Println(err.Error())
	}
}
