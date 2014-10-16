package services

import (
	"fmt"
	"log"

	"sync"
)

type TypesMap map[string]*List

type serviceInfo struct {
	address           string
	failedConnections int
}

type List struct {
	list                   []*serviceInfo
	it                     int
	mutexSL                *sync.Mutex
	failedConnectionsLimit int
	scheme                 string
	name                   string
	statusPath             string
}

type RedirectionPolicy struct {
	Path                  string
	Name                  string
	DisableStatusChecking bool
	Scheme                string
	StatusPath            string
}

func NewList(rc RedirectionPolicy) *List {
	if rc.Scheme == "" {
		rc.Scheme = "http"
	}
	if rc.StatusPath == "" {
		rc.StatusPath = "/status"
	}
	if rc.Path == "" {
		log.Fatalf("Missing path in RedirectionPolicy")
	}
	if rc.Name == "" {
		log.Fatalf("Missing name in RedirectionPolicy")
	}

	l := &List{
		it:                     -1,
		list:                   make([]*serviceInfo, 0, 0),
		mutexSL:                &sync.Mutex{},
		failedConnectionsLimit: 5,
		scheme:                 rc.Scheme,
		name:                   rc.Name,
		statusPath:             rc.StatusPath}
	// starting status checking daemon
	if !rc.DisableStatusChecking {
		go statusChecker(l)
	}
	return l
}

func (sl *List) Scheme() string {
	return sl.scheme
}

func (sl *List) Name() string {
	return sl.name
}

func (sl *List) updateFailedConnections(i, newValue int) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()
	sl.list[i].failedConnections = newValue
}

func (sl *List) AddService(address string) error {
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

	stateChan <- 's'
	return nil
}

func (sl *List) UnregisterService(address string) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()

	for _, v := range sl.list {
		if v.address == address {
			v.failedConnections = 1000
			break
		}
	}
}

func (sl *List) removeService(i int) {
	sl.mutexSL.Lock()
	defer sl.mutexSL.Unlock()
	copy(sl.list[i:], sl.list[i+1:])
	sl.list[len(sl.list)-1] = nil
	sl.list = sl.list[:len(sl.list)-1]

	if sl.it >= i {
		sl.it--
	}

	stateChan <- 's'
}

func (sl *List) GetNext() (string, error) {
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

func (sl *List) AddressesList() []string {
	list := make([]string, len(sl.list), len(sl.list))
	for id, val := range sl.list {
		list[id] = val.address
	}
	return list
}
