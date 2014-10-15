package model

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type serviceInfo struct {
	address           string
	failedConnections int
}

type ServicesList struct {
	list                   []*serviceInfo
	it                     int
	mutexSL                *sync.Mutex
	failedConnectionsLimit int
	scheme                 string
	name                   string
	statusPath             string
	stateChan              chan byte
}

func NewServicesList(rc RedirectionPolicy, stateChan chan byte) *ServicesList {
	if rc.Scheme == "" {
		rc.Scheme = "http"
	}
	if rc.StatusPath == "" {
		rc.StatusPath = "/status"
	}
	return &ServicesList{
		it:                     -1,
		list:                   make([]*serviceInfo, 0, 0),
		mutexSL:                &sync.Mutex{},
		failedConnectionsLimit: 5,
		scheme:                 rc.Scheme,
		name:                   rc.Name,
		statusPath:             rc.StatusPath,
		stateChan:              stateChan}
}

func (sl *ServicesList) Scheme() string {
	return sl.scheme
}

func (sl *ServicesList) Name() string {
	return sl.name
}

func (sl *ServicesList) AddService(address string) error {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	for _, service := range sl.list {
		if service.address == address {
			service.failedConnections = 0
			return fmt.Errorf("Host %s already exists", address)
		}
	}

	serviceInfo := &serviceInfo{address: address}

	sl.list = append(sl.list, serviceInfo)
	//sl.stateChan <- 's'
	return nil
}

func (sl *ServicesList) UnregisterService(address string) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	for _, v := range sl.list {
		if v.address == address {
			v.failedConnections = 1000
			break
		}
	}
}

func (sl *ServicesList) removeService(i int) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()
	copy(sl.list[i:], sl.list[i+1:])
	sl.list[len(sl.list)-1] = nil
	sl.list = sl.list[:len(sl.list)-1]

	if sl.it >= i {
		sl.it--
	}

	sl.stateChan <- 's'
}

func (sl *ServicesList) GetNext() (string, error) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	if len(sl.list) == 0 {
		return "", fmt.Errorf("Services list is empty")
	}

	lenght := len(sl.list)

	sl.it = (sl.it + 1) % len(sl.list)
	for sl.list[sl.it].failedConnections > 0 {
		sl.it = (sl.it + 1) % len(sl.list)
		lenght--
		if lenght == 0 {
			return "", fmt.Errorf("All services are not responding")
		}
	}

	res := sl.list[sl.it].address
	return res, nil
}

func (sl *ServicesList) GetServicesList() []string {
	list := make([]string, len(sl.list), len(sl.list))
	for id, val := range sl.list {
		list[id] = val.address
	}
	return list
}

func (sl *ServicesList) updateFailedConnections(i, newValue int) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()
	sl.list[i].failedConnections = newValue
}

func (sl *ServicesList) CheckState() {
	for i := 0; i < len(sl.list); i++ {
		if sl.list[i].failedConnections <= sl.failedConnectionsLimit {
			resp, err := http.Get(sl.scheme + "://" + sl.list[i].address + sl.statusPath)

			if err != nil || resp.StatusCode != 200 { // TODO
				sl.updateFailedConnections(i, sl.list[i].failedConnections+1)
			} else {
				resp.Body.Close()
				sl.updateFailedConnections(i, 0)
			}
		}

		if sl.list[i].failedConnections > sl.failedConnectionsLimit {
			log.Printf(sl.name + " status check: removed " + sl.list[i].address + "\n\n")
			sl.removeService(i)
			i--
			continue
		}

		if sl.list[i].failedConnections != 0 {
			log.Printf(sl.name + " status check: " + sl.list[i].address + " failed " +
				strconv.Itoa(sl.list[i].failedConnections) + " times\n\n")
		}
	}
}
