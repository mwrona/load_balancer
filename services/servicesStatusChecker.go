package services

import (
	"scalarm_load_balancer/model"
	"time"
)

func ServicesStatusChecker(servicesList *model.ServicesList) {
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			servicesList.CheckState()
		}
	}
}
