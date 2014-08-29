package model

type Context struct {
	ExperimentManagersList *ServicesList
	StorageManagersList    *ServicesList
	LoadBalancerAddress    string
	LoadBalancerScheme     string
}
