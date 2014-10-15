package handler

import "net/http"

func Authentication(allowedAddress string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != "localhost" && r.Host != allowedAddress {
			http.Error(w, "Registration from remote client is forbidden", 403)
			return
		}

		h.ServeHTTP(w, r)
	})
}
