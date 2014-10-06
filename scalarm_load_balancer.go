package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"scalarm_load_balancer/handlers"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/services"
	"scalarm_load_balancer/utils"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//loading config
	var configFile string
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		configFile = "config.json"
	}
	config, err := model.LoadConfig(configFile)
	if err != nil {
		log.Fatal("An error occurred while loading configuration: " + configFile + "\n" + err.Error())
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
	//creating channel for requests about saving state
	stateChan := make(chan byte, 15)

	//creating services lists and names maps
	redirectionsList := make(map[string]*model.ServicesList)
	servicesTypesList := make(map[string]*model.ServicesList)
	for _, rc := range config.RedirectionConfig {
		nsl := model.NewServicesList(rc.Scheme, rc.Name, stateChan)
		redirectionsList[rc.Path] = nsl
		servicesTypesList[rc.Name] = nsl

		// starting status checking deamons
		if !rc.DisableStatusChecking {
			go services.ServicesStatusChecker(redirectionsList[rc.Path])
		}
	}
	//loading state if exeists
	services.LoadState(servicesTypesList)

	//adding information service data from config
	var informationServiceScheme string
	if is, ok := redirectionsList["/information"]; ok {
		is.AddService(config.InformationServiceAddress)
		informationServiceScheme = is.Scheme()
	} else {
		log.Fatal("Information Service not specified")
	}

	//setting context
	context := &model.Context{
		RedirectionsList:    redirectionsList,
		ServicesTypesList:   servicesTypesList,
		LoadBalancerAddress: config.PrivateLoadBalancerAddress,
		LoadBalancerScheme:  config.LoadBalancerScheme,
		StateChan:           stateChan,
	}

	//setting reverse proxy
	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context)}
	http.Handle("/", reverseProxy)

	http.Handle("/register", model.ContextHandler(context, handlers.RegistrationHandler))
	//http.Handle("/unregister", model.ServicesListHandler(context.RedirectionsList,
	//	handlers.UnregistrationHandler)))
	http.Handle("/list", model.ContextHandler(context, handlers.ListHandler))

	http.HandleFunc("/error/", handlers.ErrorHandler)

	//information service registration
	if _, err := utils.RepetitiveCaller(
		func() (interface{}, error) {
			return nil, utils.InformationServiceRegistration(config.PublicLoadBalancerAddress,
				config.InformationServiceAddress,
				informationServiceScheme,
				config.InformationServiceUser,
				config.InformationServicePass)
		},
		nil,
		"InformationServiseRegistration",
	); err != nil {
		log.Fatal("Registration to Information Service failed")
	}

	//starting services
	go services.StartMulticastAddressSender(config.PrivateLoadBalancerAddress, config.MulticastAddress)
	go services.StateDeamon(context.StateChan, servicesTypesList)

	//setting up server
	server := &http.Server{
		Addr: ":" + config.Port,
	}

	if config.LoadBalancerScheme == "http" {
		err = server.ListenAndServe()
	} else { // "https"
		//redirect http to https
		if config.Port == "443" {
			go func() {
				serverHTTP := &http.Server{
					Addr: ":80",
					Handler: http.HandlerFunc(
						func(w http.ResponseWriter, req *http.Request) {
							http.Redirect(w, req, "https://"+config.PublicLoadBalancerAddress+req.RequestURI,
								http.StatusMovedPermanently)
						}),
				}
				err = serverHTTP.ListenAndServe()
				if err != nil {
					log.Fatal("An error occurred while running service on port 80\n" + err.Error())
				}
			}()
		}

		err = server.ListenAndServeTLS(config.CertFilePath, config.KeyFilePath)
	}
	if err != nil {
		log.Fatal("An error occurred while running service on port " + config.Port + "\n" + err.Error())
	}
}
