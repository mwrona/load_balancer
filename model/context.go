package model

type Context struct {
	ExperimentManagersList    *ServicesList
	StorageManagersList       *ServicesList
	InformationServiceAddress string
	InformationServiceScheme  string
	LoadBalancerAddress       string
	LoadBalancerScheme        string
}
