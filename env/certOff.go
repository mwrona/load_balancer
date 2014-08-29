// +build certOff

package env

import (
	"crypto/tls"
	"net/http"
)

var TLSClientConfigCert = &tls.Config{InsecureSkipVerify: true}

var TransportCert = &http.Transport{
	TLSClientConfig: TLSClientConfigCert,
}

const CertOff = true
