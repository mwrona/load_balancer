// +build prod

package env

import "net/http"

const Protocol = "https"

func StartServer(server *http.Server, certFile, keyFile string) error {
	return server.ListenAndServeTLS(certFile, keyFile)
}
