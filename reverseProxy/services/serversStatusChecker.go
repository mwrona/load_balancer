package services

import (
	"load_balancer/reverseProxy/model"
	"time"
)

func ServersStatusChecker(serversList *model.ServersList) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			serversList.CheckState()
		}
	}
}
