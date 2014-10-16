package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"scalarm_load_balancer/handler"
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
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatal("An error occurred while loading configuration: " + configFile + "\n" + err.Error())
	}

	//creating services lists and names maps
	redirectionsList := make(services.TypesMap)
	servicesTypesList := make(services.TypesMap)
	for _, rc := range config.RedirectionConfig {
		nsl := services.NewList(rc)
		redirectionsList[rc.Path] = nsl
		servicesTypesList[rc.Name] = nsl
	}
	//do not modify redirectionsList and servicesTypesList after this point
	//StateDeamon - on signal saves state
	go services.StartStateDaemon(servicesTypesList)
	//loading state if exists
	go services.LoadState(servicesTypesList)

	//setting context
	context := handler.AppContext(redirectionsList, servicesTypesList, config.LoadBalancerScheme)

	//disabling certificate checking
	TLSClientConfigCert := &tls.Config{InsecureSkipVerify: true}
	TransportCert := &http.Transport{
		TLSClientConfig: TLSClientConfigCert,
	}

	//setting routing
	director := handler.ReverseProxyDirector(context)
	reverseProxy := &httputil.ReverseProxy{Director: director, Transport: TransportCert}
	http.Handle("/", handler.Context(nil, handler.Websocket(director, reverseProxy)))

	http.Handle("/register", handler.Authentication(
		config.PrivateLoadBalancerAddress,
		handler.Context(
			context,
			handler.ServicesManagment(handler.Registration))))

	http.Handle("/unregister", handler.Authentication(
		config.PrivateLoadBalancerAddress,
		handler.Context(
			context,
			handler.ServicesManagment(handler.Unregistration))))

	http.Handle("/list", handler.Context(context, handler.List))

	http.HandleFunc("/error", handler.RedirectionError)

	//starting periodical multicast
	go StartMulticastAddressSender(config.PrivateLoadBalancerAddress, config.MulticastAddress)

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
