package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"scalarm_load_balancer/handler"
	"scalarm_load_balancer/model"
	"scalarm_load_balancer/services"

	"github.com/natefinch/lumberjack"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetOutput(&lumberjack.Logger{
		Dir:        "log",
		MaxSize:    100 * lumberjack.Megabyte,
		MaxBackups: 3,
		MaxAge:     28, //days
	})

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

	//setting context
	context := &model.Context{
		RedirectionsList:    redirectionsList,
		ServicesTypesList:   servicesTypesList,
		LoadBalancerAddress: config.PrivateLoadBalancerAddress,
		LoadBalancerScheme:  config.LoadBalancerScheme,
		StateChan:           stateChan,
	}
	//disabling certificate checking
	TLSClientConfigCert := &tls.Config{InsecureSkipVerify: true}
	TransportCert := &http.Transport{
		TLSClientConfig: TLSClientConfigCert,
	}

	//setting reverse proxy
	director := handler.ReverseProxyDirector(context)
	reverseProxy := &httputil.ReverseProxy{Director: director, Transport: TransportCert}
	http.Handle("/", handler.Context(nil, handler.Websocket(director, reverseProxy)))

	http.Handle("/register", handler.Authentication(config.PrivateLoadBalancerAddress, handler.Context(context, handler.Registration)))
	//http.Handle("/unregister", model.ServicesListHandler(context.RedirectionsList,
	//	handlers.UnregistrationHandler)))
	http.Handle("/list", handler.Context(context, handler.List))

	http.HandleFunc("/error/", handler.RedirectionError)

	//starting services
	go services.StartMulticastAddressSender(config.PrivateLoadBalancerAddress, config.MulticastAddress)
	go services.StateDeamon(context.StateChan, servicesTypesList)

	//setting up server
	server := &http.Server{
		Addr:      ":" + config.Port,
		TLSConfig: TLSClientConfigCert,
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
						func(w http.ResponseWriter, r *http.Request) {
							http.Redirect(w, r, "https://"+r.Host+r.RequestURI,
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
