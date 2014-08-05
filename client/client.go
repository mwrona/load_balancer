package main

import (
	"fmt"
	"load_balancer/reverseProxy/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

func request(ch chan int) {
	start := time.Now()
	resp, err := http.Get("http://" + proxyAddress + ":" + proxyPort + "/query/get")
	respTime += time.Since(start)

	utils.Check(err)
	resp.Body.Close()
	ch <- 0
}

var count = 0
var respTime = time.Duration(0)
var proxyAddress = "192.168.202.80"
var proxyPort = "8080"

func main() {
	maxCount, err := strconv.Atoi(os.Args[1])
	utils.Check(err)

	ch := make(chan int, maxCount)
	for count = 0; count < maxCount; count++ {
		go request(ch)
	}

	for i := 0; i < maxCount; i++ {
		<-ch
	}

	fmt.Println("Average response time:", respTime/time.Duration(count))
}
