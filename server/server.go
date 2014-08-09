package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"load_balancer/reverseProxy/model"
	"load_balancer/reverseProxy/utils"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var port string
var address string

var proxyAddress string
var proxyPort string

var protocol string = "https"

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s port, %s", port, r.URL.Path)
	log.Printf("Server : Query received")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
	log.Printf("Server : Status send")
}

func unregisterHandler(w http.ResponseWriter, r *http.Request) {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var Client = &http.Client{Transport: tr}

	resp, err := Client.PostForm(protocol+"://"+proxyAddress+":"+proxyPort+"/unregister", url.Values{"port": {port}, "address": {address}})
	log.Printf("Server : Unregistering " + address + ":" + port)
	utils.Check(err)
	resp.Body.Close()
	os.Exit(0)
}

func waitForProxyAddress() {
	config, err := model.LoadConfig()
	utils.Check(err)

	mcaddr, err := net.ResolveUDPAddr("udp", config.Address)
	utils.Check(err)

	conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	utils.Check(err)

	log.Printf("Server : Listen for proxy address on multicast address " + config.Address)
	b := make([]byte, 20)
	_, _, err = conn.ReadFromUDP(b)
	utils.Check(err)

	b = bytes.Trim(b, "\x00")
	rawData := strings.Split(string(b), ":")
	proxyAddress = rawData[0]
	proxyPort = rawData[1]
	log.Printf("Server : Proxy sending from " + proxyAddress + ":" + proxyPort)
	conn.Close()
}

func main() {
	port = os.Args[1]
	address = utils.GetIP()
	address = "localhost"
	/*address = "10.22.115.158"
	proxyAddress = "10.22.109.142"
	proxyPort = "8080"*/

	log.Printf("Server : address: " + address + ":" + port)
	waitForProxyAddress()

	//certOff
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var Client = &http.Client{Transport: tr}

	//

	resp, err := Client.PostForm(protocol+"://"+proxyAddress+":"+proxyPort+"/register", url.Values{"port": {port}, "address": {address}})
	log.Printf("Server : Sending address to reverse proxy")
	utils.Check(err)

	resp.Body.Close()

	http.HandleFunc("/", getHandler)
	http.HandleFunc("/status/", getHandler)
	http.HandleFunc("/unregister", unregisterHandler)
	log.Printf("Server : start")
	http.ListenAndServe(":"+port, nil)
}
