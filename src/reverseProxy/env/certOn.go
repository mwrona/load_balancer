// +build !certOff

package env

import (
	"crypto/tls"
	"net/http"
)

var TLSClientConfigCert = &tls.Config{}
var TransportCert = &http.Transport{}

const CertOff = false
