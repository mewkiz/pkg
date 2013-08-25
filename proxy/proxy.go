// Package proxy provides proxy server utility functions.
package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// Listen initiates a reverse proxy on a given port. It handles incoming client
// requests and sends them to their intended target server, proxying the
// response back to the client.
//
// director has the ability to modify requests before they are sent. It can also
// be used to simply perform an action based on the requests.
func Listen(port int, director func(req *http.Request)) (err error) {
	// handle handles incoming client requests and sends them to their intended
	// target server, proxying the response back to the client.
	handle := func(w http.ResponseWriter, req *http.Request) {
		proxy := httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, req)
	}
	http.HandleFunc("/", handle)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}
