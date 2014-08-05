package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"reverseProxy/env"
	"reverseProxy/utils"
	"time"

	"code.google.com/p/go.net/ipv4"
)

var serverList *utils.ServerList

var proxyAddress string
var proxyPort string

func queryHandler(req *http.Request) {

	if serverList.IsEmpty() {
		req.URL.Scheme = env.Protocol
		req.URL.Path = "/error/"
		req.URL.RawQuery = ""
		req.URL.Host = proxyAddress + ":" + proxyPort
		utils.Log("Reverse Proxy", "error, empty server list")
	} else {
		req.URL.Scheme = "http"
		req.URL.Host = serverList.GetNext()
		req.URL.Path = "/get/"
		req.URL.RawQuery = ""
		utils.Log("Reverse Proxy", "redirect to "+req.URL.Host+req.URL.Path)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	address := r.FormValue("address")
	if port == "" {
		fmt.Fprintf(w, "Error: missing port")
		utils.Log("Reverse Proxy", "error, missing port")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		utils.Log("Reverse Proxy", "error, missing address")
	} else {
		serverList.AddServer(address, port)
		fmt.Fprintf(w, "Registered server:  %s", address+":"+port)
		utils.Log("Reverse Proxy", "registered server: "+address+":"+port)
	}
}

func unregisterHandler(w http.ResponseWriter, r *http.Request) {
	port := r.FormValue("port")
	address := r.FormValue("address")
	if port == "" {
		fmt.Fprintf(w, "Error: missing port")
		utils.Log("Reverse Proxy", "error, missing port")
	} else if address == "" {
		fmt.Fprintf(w, "Error: missing address")
		utils.Log("Reverse Proxy", "error, missing address")
	} else {
		serverList.UnregisterServer(address, port)
		fmt.Fprintf(w, "Unregistered server:  %s", address+":"+port)
		utils.Log("Reverse Proxy", "unregistered server: "+address+":"+port)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server list is empty!")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log("Reverse Proxy", "printing servers list")
	fmt.Fprintf(w, "Servers available:\n")
	for _, val := range serverList.GetServerList() {
		fmt.Fprintf(w, val+"\n")
	}
}

func multicastAddressSender() {
	config, err := utils.LoadConfig()
	utils.Check(err)

	mcaddr, err := net.ResolveUDPAddr("udp", config.Address)
	utils.Check(err)

	// conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	// utils.Check(err)

	c, err := net.ListenPacket("udp4", "")
	utils.Check(err)
	defer c.Close()

	conn := ipv4.NewPacketConn(c)

	conn.JoinGroup(nil, mcaddr)
	conn.SetMulticastLoopback(true)

	b := make([]byte, 20)
	copy(b, proxyAddress+":"+proxyPort)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			// _, err = conn.WriteToUDP(b, mcaddr)
			_, err = conn.WriteTo(b, nil, mcaddr)
			utils.Check(err)
		}
	}
}

func serversStatusChecker() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			serverList.CheckState()
		}
	}
}

func main() {
	if env.CertOff {
		utils.Log("Reverse Proxy", "Certificates checking disabled")
	}
	serverList = utils.CreateServerList()
	proxyPort = "8080"
	proxyAddress = utils.GetIP()
	proxyAddress = "localhost"

	h := &httputil.ReverseProxy{Director: queryHandler, Transport: env.TransportCert}

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/unregister", unregisterHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/error/", errorHandler)

	http.Handle("/query/", h)

	go multicastAddressSender()
	go serversStatusChecker()

	utils.Log("Reverse Proxy", "Start")

	server := &http.Server{Addr: env.Protocol + proxyAddress + ":" + proxyPort, TLSConfig: env.TLSClientConfigCert}

	//err := server.ListenAndServe()
	//err := server.ListenAndServeTLS("cert.pem", "key.pem")
	err := env.StartServer(server, "cert.pem", "key.pem")
	utils.Check(err)
}
