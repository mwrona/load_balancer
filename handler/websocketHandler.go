package handler

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"scalarm_load_balancer/model"
)

func isWebsocket(req *http.Request) bool {
	return req.Header.Get("Connection") == "Upgrade" && req.Header.Get("Upgrade") == "websocket"
}

func Websocket(director func(*http.Request), h http.Handler) contextHandlerFunction {
	errorMessage := "Unable to establish websocket connection"
	return func(context *model.Context, w http.ResponseWriter, req *http.Request) error {
		if !isWebsocket(req) {
			h.ServeHTTP(w, req)
			return nil
		}
		director(req)
		if req.URL.Host == "/error/" {
			RedirectionError(w, req)
			return nil
		}
		//weboscket
		w.Header().Set("Upgrade", "websocket")

		//connection to target
		var d net.Conn
		var err error
		if req.URL.Scheme == "http" {
			d, err = net.Dial("tcp", req.URL.Host)
		} else {
			d, err = tls.Dial("tcp", req.URL.Host, &tls.Config{InsecureSkipVerify: true})
		}

		if err != nil {
			log.Printf("Websocket connection failed: %v", err.Error())
			return model.NewHTTPError(errorMessage, 404)
		}
		defer d.Close()

		//hijacking connection
		hj, ok := w.(http.Hijacker)
		if !ok {
			log.Printf("Websocket connection failed: Casting to http.Hijacker failed")
			return model.NewHTTPError(errorMessage, 500)
		}
		nc, _, err := hj.Hijack()
		if err != nil {
			log.Printf("Websocket connection failed: %v", err.Error())
			return model.NewHTTPError(errorMessage, 500)
		}
		defer nc.Close()

		//copying requst to target and responce to client
		err = req.Write(d)
		if err != nil {
			log.Printf("Websocket connection failed: %v", err.Error())
			return model.NewHTTPError(errorMessage, 404)
		}

		errc := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errc <- err
		}
		go cp(d, nc)
		go cp(nc, d)
		<-errc
		return nil
	}
}
