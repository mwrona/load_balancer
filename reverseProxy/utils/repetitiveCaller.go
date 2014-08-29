package utils

import (
	"log"
	"strconv"
	"time"
)

func RepititveCaller(f func() (interface{}, error), intervals []int) (out interface{}, err error) {
	if intervals == nil {
		intervals = []int{15, 30, 60, 120, 240}
	}

	intervals = append(intervals, -1)

	for _, duration := range intervals {
		out, err = f()
		if err == nil || duration == -1 {
			return
		}
		log.Printf("RepititveCaller : call failed, err: \n" + err.Error() + "\nReattempt in " + strconv.Itoa(duration) + "s")
		time.Sleep(time.Second * time.Duration(duration))
	}
	return
}
