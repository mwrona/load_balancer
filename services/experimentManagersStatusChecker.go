package services

import (
	"scalarm_load_balancer/model"
	"time"
)

func ExperimentManagersStatusChecker(experimentManagersList *model.ExperimentManagersList) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			experimentManagersList.CheckState()
		}
	}
}
