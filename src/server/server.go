package main

import (
    "fmt"
	"net"
    "net/http"
    "os"
	"net/url"
	"reverseProxy/utils"
	"strings"
	"bytes"
)

var port string
var address string

var proxyAddress string
var proxyPort string

func getHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s port, %s", port, r.URL.Path)
    utils.Log("Server", "Query received")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
    utils.Log("Server", "Status send")
}

func unregisterHandler(w http.ResponseWriter, r *http.Request) {
    resp, err := http.PostForm("http://" + proxyAddress + ":" + proxyPort + "/unregister", url.Values{"port": {port}, "address": {address}})
    utils.Log("Server", "Unregistering " + address + ":" + port)
	utils.Check(err)
	resp.Body.Close()
	os.Exit(0)
}

func waitForProxyAddress() {
	config, err := utils.LoadConfig()
	utils.Check(err)
	
	mcaddr, err := net.ResolveUDPAddr("udp", config.Address)
	utils.Check(err)
	
	conn, err := net.ListenMulticastUDP("udp", nil, mcaddr)
	utils.Check(err)
	
	utils.Log("Server", "Listen for proxy address on multicast address " + config.Address)
	b := make([]byte, 20)
	_, _, err = conn.ReadFromUDP(b)
	utils.Check(err)
	
	b = bytes.Trim(b, "\x00")
	rawData := strings.Split(string(b), ":")
	proxyAddress = rawData[0]
	proxyPort = rawData[1]
	utils.Log("Server", "Proxy sending from " + proxyAddress + ":" + proxyPort)
	conn.Close()
}

func main() {
	port = os.Args[1]
	address = utils.GetIP()
	
	/*address = "10.22.115.158"
	proxyAddress = "10.22.109.142"
	proxyPort = "8080"*/
	
	utils.Log("Server", "address: " + address + ":" + port)
	waitForProxyAddress()
	
	resp, err := http.PostForm("http://" + proxyAddress + ":" + proxyPort + "/register", url.Values{"port": {port}, "address": {address}})
	utils.Log("Server", "Sending address to reverse proxy")
	utils.Check(err)
	
	resp.Body.Close()
	
    http.HandleFunc("/get/", getHandler)
    http.HandleFunc("/status/", getHandler)
    http.HandleFunc("/unregister", unregisterHandler)
    utils.Log("Server", "start")
    http.ListenAndServe(":" + port, nil)
}
