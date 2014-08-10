package main

import (
	"load_balancer/reverseProxy/env"
	"load_balancer/reverseProxy/handlers"
	"load_balancer/reverseProxy/model"
	"load_balancer/reverseProxy/services"
	"load_balancer/reverseProxy/utils"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	if env.CertOff {
		log.Printf("Reverse Proxy : Certificates checking disabled")
	}

	context := &model.Context{
		ServersList:  model.CreateServersList(),
		ProxyPort:    "8080",
		ProxyAddress: "localhost"}

	reverseProxy := &httputil.ReverseProxy{Director: handlers.ReverseProxyDirector(context), Transport: env.TransportCert}
	http.Handle("/", reverseProxy)

	http.Handle("/register", model.ContextHandler(context, handlers.RegisterHandler))
	http.Handle("/unregister", model.ContextHandler(context, handlers.UnregisterHandler))
	http.Handle("/list", model.ContextHandler(context, handlers.ListHandler))
	http.HandleFunc("/error/", handlers.ErrorHandler)

	go services.MulticastAddressSender(context.ProxyAddress, context.ProxyPort)
	go services.ServersStatusChecker(context.ServersList)

	log.Printf("Reverse Proxy : Start")

	server := &http.Server{
		Addr:      context.ProxyAddress + ":" + context.ProxyPort,
		TLSConfig: env.TLSClientConfigCert,
	}

	//err := server.ListenAndServe()
	//err := server.ListenAndServeTLS("cert.pem", "key.pem")
	err := env.StartServer(server, "cert.pem", "key.pem")
	utils.Check(err)
}
