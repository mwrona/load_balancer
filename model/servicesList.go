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
	schema                 string
}

func NewServicesList(schema string) *ServicesList {
	//log.Printf("Service List : CreateServicesList")
	return &ServicesList{it: -1, list: make([]*serviceInfo, 0, 0), mutexSL: &sync.Mutex{},
		failedConnectionsLimit: 2, schema: schema}
}

func (sl *ServicesList) AddService(address string) error {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	for _, service := range sl.list {
		if service.address == address {
			//log.Printf("Service List : AddService: host already exists")
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
	sl.list = append(sl.list[:i], sl.list[i+1:]...)
	//log.Printf("Service List : RemoveService: deleted " + sl.list[i].address + ":" + sl.list[i].port)

	if sl.it >= i {
		sl.it--
	}
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

func (sl *ServicesList) CheckState() {
	for i := 0; i < len(sl.list); i++ {
		if sl.list[i].failedConnections <= sl.failedConnectionsLimit {
			resp, err := http.Get(sl.schema + "://" + sl.list[i].address + "/status")
			resp.Body.Close()
			sl.mutexSL.Lock()
			if err != nil { // TODO
				sl.list[i].failedConnections++
			} else {
				sl.list[i].failedConnections = 0
			}
			sl.mutexSL.Unlock()
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
