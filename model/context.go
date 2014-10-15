package model

type Context struct {
	RedirectionsList    map[string]*ServicesList
	ServicesTypesList   map[string]*ServicesList
	LoadBalancerAddress string
	LoadBalancerScheme  string
	StateChan           chan byte
}
