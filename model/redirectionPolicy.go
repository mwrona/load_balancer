package model

type RedirectionPolicy struct {
	Path                  string
	Name                  string
	DisableStatusChecking bool
	Scheme                string
	StatusPath            string
}
