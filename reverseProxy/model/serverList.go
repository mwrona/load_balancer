package model

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type ServerInfo struct {
	port              string
	address           string
	failedConnections int
}

type ServersList struct {
	list                   []*ServerInfo
	it                     int
	mutexSL                *sync.Mutex
	failedConnectionsLimit int
}

func (serverInfo *ServerInfo) GetFullAddress() string {
	return serverInfo.address + ":" + serverInfo.port
}

func CreateServersList() *ServersList {
	//log.Printf("Server List : CreateServersList")
	return &ServersList{it: -1, list: make([]*ServerInfo, 0, 0), mutexSL: &sync.Mutex{}, failedConnectionsLimit: 2}
}

func (ServersList *ServersList) AddServer(address, port string) error {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()

	for _, server := range ServersList.list {
		if server.port == port && server.address == address {
			//log.Printf("Server List : AddServer: host already exists")
			return errors.New("Server List : AddServer: host already exists")
		}
	}

	serverInfo := &ServerInfo{port: port, address: address}

	ServersList.list = append(ServersList.list, serverInfo)
	//log.Printf("Server List : AddServer: added " + ServersList.list[len(ServersList.list)-1].GetFullAddress())
	return nil
}

func (ServersList *ServersList) UnregisterServer(address, port string) {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()

	for _, v := range ServersList.list {
		if v.address == address && v.port == port {
			v.failedConnections = 1000
			break
		}
	}
	//log.Printf("Server List : UnregisterServer: unregister " + address + ":" + port)
}

func (ServersList *ServersList) removeServer(i int) {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()
	ServersList.list = append(ServersList.list[:i], ServersList.list[i+1:]...)
	//log.Printf("Server List : RemoveServer: deleted " + ServersList.list[i].address + ":" + ServersList.list[i].port)

	if ServersList.it >= i {
		ServersList.it--
	}
}

func (ServersList *ServersList) GetNext() (string, error) {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()

	if len(ServersList.list) == 0 {
		//log.Printf("Server List : GetNext: server list is empty")
		return "", errors.New("Server List : GetNext: server list is empty")
	}

	lenght := len(ServersList.list)

	ServersList.it = (ServersList.it + 1) % len(ServersList.list)
	for ServersList.list[ServersList.it].failedConnections > 0 {
		ServersList.it = (ServersList.it + 1) % len(ServersList.list)
		lenght--
		if lenght == 0 {
			//log.Print("Server List : GetNext: all servers are not responding")
			return "", errors.New("Server List : GetNext: all servers are not responding")
		}
	}

	res := ServersList.list[ServersList.it].GetFullAddress()
	//log.Printf("Server List : GetNext: got " + res)
	return res, nil
}

func (ServersList *ServersList) GetServersList() []string {
	list := make([]string, len(ServersList.list), len(ServersList.list))
	for id, val := range ServersList.list {
		list[id] = val.GetFullAddress()
	}
	return list
}

func (ServersList *ServersList) CheckState() {
	for i := 0; i < len(ServersList.list); i++ {
		if ServersList.list[i].failedConnections <= ServersList.failedConnectionsLimit {
			_, err := http.Get("http://" + ServersList.list[i].address + ":" + ServersList.list[i].port + "/status")
			ServersList.mutexSL.Lock()
			if err != nil { // TODO
				ServersList.list[i].failedConnections++
			} else {
				ServersList.list[i].failedConnections = 0
			}
			ServersList.mutexSL.Unlock()
		}
		if ServersList.list[i].failedConnections > ServersList.failedConnectionsLimit {
			log.Printf("Server List : Check State: removed " + ServersList.list[i].address + ":" + ServersList.list[i].port + "\n\n")
			ServersList.removeServer(i)
			i--
		} else {
			if ServersList.list[i].failedConnections == 0 {
				//log.Printf("Server List : Check State: " + ServersList.list[i].address + ":" + ServersList.list[i].port + " ok")
			} else {
				log.Printf("Server List : Check State: " + ServersList.list[i].address + ":" +
					ServersList.list[i].port + " failed " + strconv.Itoa(ServersList.list[i].failedConnections) + " times\n\n")
			}

		}
	}
}
