// +build !prod

package env

import "net/http"

const Protocol = "http://"

func StartServer(server *http.Server, certFile, keyFile string) error {
	return server.ListenAndServe()
}
