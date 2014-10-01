package model

type Context struct {
	RedirectionsList          map[string]*ServicesList
	InformationServiceAddress string
	InformationServiceScheme  string
	LoadBalancerAddress       string
	LoadBalancerScheme        string
}
