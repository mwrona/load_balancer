package model

type SerivesesListMap map[string]*ServicesList

type Context struct {
	RedirectionsList    SerivesesListMap
	ServicesTypesList   SerivesesListMap
	LoadBalancerAddress string
	LoadBalancerScheme  string
	StateChan           chan byte
}
