package model

import (
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
	list    []*ServerInfo
	it      int
	mutexSL *sync.Mutex
}

func (serverInfo *ServerInfo) GetFullAddress() string {
	return serverInfo.address + ":" + serverInfo.port
}

func CreateServersList() *ServersList {
	log.Printf("Server List : CreateServersList")
	return &ServersList{it: -1, list: make([]*ServerInfo, 0, 0), mutexSL: &sync.Mutex{}}
}

func (ServersList *ServersList) AddServer(address, port string) {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()

	serverInfo := &ServerInfo{port: port, address: address}

	ServersList.list = append(ServersList.list, serverInfo)
	log.Printf("Server List : AddServer: added " + ServersList.list[len(ServersList.list)-1].GetFullAddress())
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
	log.Printf("Server List : UnregisterServer: unregister " + address + ":" + port)
}

func (ServersList *ServersList) removeServer(i int) {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()
	log.Printf("Server List : RemoveServer: deleted " + ServersList.list[i].address + ":" + ServersList.list[i].port)
	ServersList.list = append(ServersList.list[:i], ServersList.list[i+1:]...)

	if ServersList.it >= i {
		ServersList.it--
	}
}

func (ServersList *ServersList) GetNext() string {
	ServersList.mutexSL.Lock()
	defer ServersList.mutexSL.Unlock()

	ServersList.it = (ServersList.it + 1) % len(ServersList.list)
	for ServersList.list[ServersList.it].failedConnections > 0 {
		ServersList.it = (ServersList.it + 1) % len(ServersList.list)
	}

	res := ServersList.list[ServersList.it].GetFullAddress()
	log.Printf("Server List : GetNext: got " + res)
	return res //TODO nil if empty
}

func (ServersList *ServersList) GetServersList() []string {
	list := make([]string, len(ServersList.list), len(ServersList.list))
	for id, val := range ServersList.list {
		list[id] = val.GetFullAddress()
	}
	return list
}

func (ServersList *ServersList) IsEmpty() bool {
	if len(ServersList.list) == 0 {
		return true
	}
	return false
}

func (ServersList *ServersList) CheckState() {
	for i := 0; i < len(ServersList.list); i++ {
		if ServersList.list[i].failedConnections <= 5 {
			_, err := http.Get("http://" + ServersList.list[i].address + ":" + ServersList.list[i].port + "/status")
			if err != nil { // TODO
				ServersList.list[i].failedConnections++
			}
		}
		if ServersList.list[i].failedConnections > 0 {
			log.Printf("Server List : Check State: removed" + ServersList.list[i].address + ":" + ServersList.list[i].port)
			ServersList.removeServer(i)
			i--
		} else {
			if ServersList.list[i].failedConnections == 0 {
				log.Printf("Server List : Check State: " + ServersList.list[i].address + ":" + ServersList.list[i].port + " ok")
			} else {
				log.Printf("Server List : Check State: " + ServersList.list[i].address + ":" +
					ServersList.list[i].port + " failed " + strconv.Itoa(ServersList.list[i].failedConnections) + " times")
			}

		}
	}
}
