package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/steabert/ko"
)

func main() {
	var public string
	var backend string
	var secure bool
	var port int
	flag.StringVar(&public, "public", "", "directory or archive to serve files from")
	flag.StringVar(&backend, "backend", "", "fallback URL")
	flag.BoolVar(&secure, "secure", false, "use HTTPS")
	flag.IntVar(&port, "port", 8080, "port number")
	flag.Parse()

	var handler http.Handler = nil

	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Println("invalid backend URL: ", err.Error())
			return
		}
		fmt.Println("...enabling reverse proxy: ", backendURL)
		handler = ko.NewProxyMiddleware(*backendURL)(handler)
	}

	if public != "" {
		fmt.Println("...serving from folder: ", public)
		handler = ko.NewStaticMiddleware(public)(handler)
	}

	var scheme string
	if secure {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("🐮 listening on %s://%s, what would you like me to serve? ...\n", scheme, host)
	if secure {
		http.ListenAndServeTLS(host, "server.cert", "server.key", handler)
	} else {
		http.ListenAndServe(host, handler)
	}
}
