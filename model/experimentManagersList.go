package model

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type experimentManagerInfo struct {
	port              string
	address           string
	failedConnections int
}

type ExperimentManagersList struct {
	list                   []*experimentManagerInfo
	it                     int
	mutexSL                *sync.Mutex
	failedConnectionsLimit int
}

func (experimentManagerInfo *experimentManagerInfo) FullAddress() string {
	return experimentManagerInfo.address + ":" + experimentManagerInfo.port
}

func NewExperimentManagersList() *ExperimentManagersList {
	//log.Printf("Server List : CreateExperimentManagersList")
	return &ExperimentManagersList{it: -1, list: make([]*experimentManagerInfo, 0, 0), mutexSL: &sync.Mutex{}, failedConnectionsLimit: 2}
}

func (eml *ExperimentManagersList) AddServer(address, port string) error {
	eml.mutexSL.Lock()
	defer eml.mutexSL.Unlock()

	for _, server := range eml.list {
		if server.port == port && server.address == address {
			//log.Printf("Server List : AddServer: host already exists")
			return errors.New("Server List : AddServer: host already exists")
		}
	}

	experimentManagerInfo := &experimentManagerInfo{port: port, address: address}

	eml.list = append(eml.list, experimentManagerInfo)
	//log.Printf("Server List : AddServer: added " + eml.list[len(eml.list)-1].GetFullAddress())
	return nil
}

func (eml *ExperimentManagersList) UnregisterServer(address, port string) {
	eml.mutexSL.Lock()
	defer eml.mutexSL.Unlock()

	for _, v := range eml.list {
		if v.address == address && v.port == port {
			v.failedConnections = 1000
			break
		}
	}
	//log.Printf("Server List : UnregisterServer: unregister " + address + ":" + port)
}

func (eml *ExperimentManagersList) removeServer(i int) {
	eml.mutexSL.Lock()
	defer eml.mutexSL.Unlock()
	eml.list = append(eml.list[:i], eml.list[i+1:]...)
	//log.Printf("Server List : RemoveServer: deleted " + eml.list[i].address + ":" + eml.list[i].port)

	if eml.it >= i {
		eml.it--
	}
}

func (eml *ExperimentManagersList) GetNext() (string, error) {
	eml.mutexSL.Lock()
	defer eml.mutexSL.Unlock()

	if len(eml.list) == 0 {
		//log.Printf("Server List : GetNext: server list is empty")
		return "", errors.New("Server List : GetNext: server list is empty")
	}

	lenght := len(eml.list)

	eml.it = (eml.it + 1) % len(eml.list)
	for eml.list[eml.it].failedConnections > 0 {
		eml.it = (eml.it + 1) % len(eml.list)
		lenght--
		if lenght == 0 {
			//log.Print("Server List : GetNext: all servers are not responding")
			return "", errors.New("Server List : GetNext: all servers are not responding")
		}
	}

	res := eml.list[eml.it].FullAddress()
	//log.Printf("Server List : GetNext: got " + res)
	return res, nil
}

func (eml *ExperimentManagersList) GetExperimentManagersList() []string {
	list := make([]string, len(eml.list), len(eml.list))
	for id, val := range eml.list {
		list[id] = val.FullAddress()
	}
	return list
}

func (eml *ExperimentManagersList) CheckState() {
	for i := 0; i < len(eml.list); i++ {
		if eml.list[i].failedConnections <= eml.failedConnectionsLimit {
			resp, err := http.Get("http://" + eml.list[i].address + ":" + eml.list[i].port + "/status")
			resp.Body.Close()
			eml.mutexSL.Lock()
			if err != nil { // TODO
				eml.list[i].failedConnections++
			} else {
				eml.list[i].failedConnections = 0
			}
			eml.mutexSL.Unlock()
		}
		if eml.list[i].failedConnections > eml.failedConnectionsLimit {
			log.Printf("Server List : Check State: removed " + eml.list[i].address + ":" + eml.list[i].port + "\n\n")
			eml.removeServer(i)
			i--
		} else {
			if eml.list[i].failedConnections == 0 {
				//log.Printf("Server List : Check State: " + eml.list[i].address + ":" + eml.list[i].port + " ok")
			} else {
				log.Printf("Server List : Check State: " + eml.list[i].address + ":" +
					eml.list[i].port + " failed " + strconv.Itoa(eml.list[i].failedConnections) + " times\n\n")
			}

		}
	}
}
