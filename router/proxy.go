package router

import (
	"fmt"
	"io"
	"net/http"
)

// ProxyRouter routes requests to another server
type ProxyRouter struct {
	host string
}

// NewProxyRouter creates a proxy router
func NewProxyRouter(host string) (*ProxyRouter, error) {
	router := ProxyRouter{host: host}
	return &router, nil
}

// handleHTTP proxies regular HTTP requests
func (router ProxyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	r.Host = router.host
	r.URL.Scheme = "https"
	r.URL.Host = router.host
	rsp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// If we are dealing with a websocket handshake,
	// hijack the connection and we're done
	if rsp.StatusCode == http.StatusSwitchingProtocols {
		fmt.Println("upgrading connection...")
		router.handleUpgrade(w, r, rsp)
		return
	}

	// Transfer the roundtrip response to the response writer
	defer rsp.Body.Close()
	header := w.Header()
	for key, values := range rsp.Header {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	w.WriteHeader(rsp.StatusCode)
	io.Copy(w, rsp.Body)
}

// handleUpgrade takes over connection for WebSocket
func (router ProxyRouter) handleUpgrade(w http.ResponseWriter, r *http.Request, rsp *http.Response) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, clientBuf, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	serverConn, ok := rsp.Body.(io.ReadWriteCloser)
	if !ok {
		http.Error(w, "cannot write to response body", http.StatusInternalServerError)
		return
	}
	defer serverConn.Close()

	// Flush the headers down the buffer, body will be copied in pipe.
	rsp.Body = nil
	err = rsp.Write(clientBuf)
	if err != nil {
		http.Error(w, "Failed to write upgrade response headers", http.StatusInternalServerError)
		return
	}
	err = clientBuf.Flush()
	if err != nil {
		http.Error(w, "Failed to flush upgrade connection", http.StatusInternalServerError)
		return
	}

	errc := make(chan error, 1)
	go pipeStreams(clientConn, serverConn, errc)
	go pipeStreams(serverConn, clientConn, errc)
	<-errc

	fmt.Println("closing connection upgrade")
	return
}

func pipeStreams(dst io.Writer, src io.Reader, errc chan<- error) {
	_, err := io.Copy(dst, src)
	errc <- err
}
