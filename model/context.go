package model

type SerivesesListMap map[string]*ServicesList

type Context struct {
	RedirectionsList   SerivesesListMap
	ServicesTypesList  SerivesesListMap
	LoadBalancerScheme string
}
