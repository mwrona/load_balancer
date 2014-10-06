package model

import (
	"errors"
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
	stateChan              chan byte
}

func NewServicesList(scheme, name string, stateChan chan byte) *ServicesList {
	//log.Printf("Service List : CreateServicesList")
	if scheme == "" {
		scheme = "http"
	}
	return &ServicesList{it: -1, list: make([]*serviceInfo, 0, 0), mutexSL: &sync.Mutex{},
		failedConnectionsLimit: 6, scheme: scheme, name: name, stateChan: stateChan}
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
			//log.Printf("Service List : AddService: host already exists")
			service.failedConnections = 0
			return errors.New("ServicesList : AddService: host already exists")
		}
	}

	serviceInfo := &serviceInfo{address: address}

	sl.list = append(sl.list, serviceInfo)
	//log.Printf("Service List : AddService: added " + sl.list[len(sl.list)-1].GetFullAddress())
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
	//log.Printf("Service List : UnregisterService: unregister " + address + ":" + port)
}

func (sl *ServicesList) removeService(i int) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()
	copy(sl.list[i:], sl.list[i+1:])
	sl.list[len(sl.list)-1] = nil
	sl.list = sl.list[:len(sl.list)-1]

	//sl.list = append(sl.list[:i], sl.list[i+1:]...)
	//log.Printf("Service List : RemoveService: deleted " + sl.list[i].address + ":" + sl.list[i].port)

	if sl.it >= i {
		sl.it--
	}

	sl.stateChan <- 's'
}

func (sl *ServicesList) GetNext() (string, error) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	if len(sl.list) == 0 {
		//log.Printf("Service List : GetNext: service list is empty")
		return "", errors.New("ServicesList : GetNext: services list is empty")
	}

	lenght := len(sl.list)

	sl.it = (sl.it + 1) % len(sl.list)
	for sl.list[sl.it].failedConnections > 0 {
		sl.it = (sl.it + 1) % len(sl.list)
		lenght--
		if lenght == 0 {
			//log.Print("Service List : GetNext: all services are not responding")
			return "", errors.New("ServicesList : GetNext: all services are not responding")
		}
	}

	res := sl.list[sl.it].address
	//log.Printf("Service List : GetNext: got " + res)
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
			resp, err := http.Get(sl.scheme + "://" + sl.list[i].address + "/status")
			if err != nil { // TODO
				sl.updateFailedConnections(i, sl.list[i].failedConnections+1)
			} else {
				resp.Body.Close()
				sl.updateFailedConnections(i, 0)
			}
		}
		if sl.list[i].failedConnections > sl.failedConnectionsLimit {
			log.Printf("ServicesList : CheckState: removed " + sl.list[i].address + "\n\n")
			sl.removeService(i)
			i--
		} else {
			if sl.list[i].failedConnections == 0 {
				//log.Printf("Service List : Check State: " + sl.list[i].address + ":" + sl.list[i].port + " ok")
			} else {
				log.Printf("ServicesList : CheckState: " + sl.list[i].address + " failed " +
					strconv.Itoa(sl.list[i].failedConnections) + " times\n\n")
			}

		}
	}
}
