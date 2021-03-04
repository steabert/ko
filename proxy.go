package ko

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// ProxyRouter routes requests to another server
type ProxyRouter struct {
	backend   url.URL
	transport http.RoundTripper
}

// handleHTTP proxies regular HTTP requests
func (router ProxyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Host = router.backend.Host
	r.URL.Scheme = router.backend.Scheme
	r.URL.Host = router.backend.Host
	rsp, err := router.transport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// If we are dealing with a websocket handshake,
	// hijack the connection and we're done
	if rsp.StatusCode == http.StatusSwitchingProtocols {
		fmt.Println("...connection upgrade on: ", r.URL.Path)
		router.handleUpgrade(w, r, rsp)
		return
	}

	fmt.Println("...relay on: ", r.URL.Path)

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

// NewProxyMiddleware creates a proxy router
func NewProxyMiddleware(backend url.URL) func(http.Handler) http.Handler {
	// Keep TLS config.
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
		// Set this value so that the underlying transport round-tripper
		// doesn't try to auto decode the body of objects with
		// content-encoding set to `gzip`.
		//
		// Refer:
		//    https://golang.org/src/net/http/transport.go?h=roundTrip#L1843
		DisableCompression: true,
	}

	// router := ProxyRouter{backend: backend, transport: transport}
	proxy := httputil.NewSingleHostReverseProxy(&backend)
	proxy.Transport = transport

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
			if next != nil {
				next.ServeHTTP(w, r)
			}
		})
	}
}
