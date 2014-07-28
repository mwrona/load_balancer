package utils
import (
	"sync"
	"net/http"
	"strconv"
)

type ServerInfo struct{
	port string
	address string
	failedConnections int
}

type ServerList struct{
	list []*ServerInfo 
	it int
	mutexSL *sync.Mutex
}

func (serverInfo *ServerInfo) GetFullAddress() string {
	return serverInfo.address + ":" + serverInfo.port
}

func (serverInfo *ServerInfo) EqualTo(serverInfo2 ServerInfo) bool {
	if serverInfo.port != serverInfo2.port {
		return false
	}
	if serverInfo.address != serverInfo2.address {
		return false
	}
	if serverInfo.failedConnections != serverInfo2.failedConnections {
		return false
	}
	return true
}  

func (serverInfo *ServerInfo) Copy() *ServerInfo {
	tmp := *serverInfo
	return &tmp
}

func CreateServerList() *ServerList{
	Log("Server List", "CreateServerList")
	return &ServerList{it: -1, list: make([]*ServerInfo, 0, 0), mutexSL: &sync.Mutex{}}
}

func (serverList *ServerList) CopyListElement(i int) *ServerInfo {
	serverList.mutexSL.Lock()
	defer serverList.mutexSL.Unlock()
	return serverList.list[i].Copy()
}


func (serverList *ServerList) AddServer(address, port string) {
	serverList.mutexSL.Lock()
	defer serverList.mutexSL.Unlock()
	
	serverInfo := &ServerInfo{port: port, address: address}
	
	serverList.list = append(serverList.list, serverInfo)
	Log("Server List", "AddServer: added " + serverList.list[len(serverList.list) - 1].GetFullAddress())
}

func (serverList *ServerList) UnregisterServer(address, port string) {
	serverList.mutexSL.Lock()
	defer serverList.mutexSL.Unlock()
	
	for _, v := range(serverList.list) {
		if v.address == address && v.port == port {
			v.failedConnections = 1000 
			break
		}
	}
	Log("Server List", "UnregisterServer: unregister " + address + ":" + port)
}

func (serverList *ServerList) removeServer(i int) {
	serverList.mutexSL.Lock()
	defer serverList.mutexSL.Unlock()
	Log("Server List", "RemoveServer: deleted " + serverList.list[i].address + ":" + serverList.list[i].port)
	serverList.list = append(serverList.list[:i], serverList.list[i+1:]...)
	
	if serverList.it >= i {
		serverList.it--
	}
}

func (serverList *ServerList) GetNext() string {
	serverList.mutexSL.Lock()
	defer serverList.mutexSL.Unlock()	
	
	serverList.it = (serverList.it + 1) % len(serverList.list)
	for serverList.list[serverList.it].failedConnections > 0 {
		serverList.it = (serverList.it + 1) % len(serverList.list)
	}
	
	res := serverList.list[serverList.it].GetFullAddress()
	Log("Server List", "GetNext: got " + res)
	return res
}

func (serverList *ServerList) GetServerList() []string {
	list := make([]string, len(serverList.list), len(serverList.list))
	for id, val := range(serverList.list) {
		list[id] = val.GetFullAddress()
	}
	return list
}

func (serverList *ServerList) IsEmpty() bool {
	if(len(serverList.list) == 0) { 
		return true
	}
	return false
}

func (serverList *ServerList) CheckState() {
	for i := 0; i < len(serverList.list); i++  {
		if serverList.list[i].failedConnections <= 5 {			
			_, err := http.Get("http://" + serverList.list[i].address + ":" + serverList.list[i].port + "/status")
			if err != nil { // TODO
				serverList.list[i].failedConnections++			
			}
		}
		if serverList.list[i].failedConnections > 0 {
			Log("Server List", "Check State: removed" + serverList.list[i].address + ":" + serverList.list[i].port)
			serverList.removeServer(i)
			i--
		} else {
			if  serverList.list[i].failedConnections == 0 {
				Log("Server List", "Check State: " + serverList.list[i].address + ":" + serverList.list[i].port + " ok")
			} else {
				Log("Server List", "Check State: " + serverList.list[i].address + ":" +
							serverList.list[i].port + " failed " + strconv.Itoa(serverList.list[i].failedConnections) + " times");
			}
			
		}
	}
}
