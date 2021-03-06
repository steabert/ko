package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/steabert/ko"
)

func main() {
	var backend string
	var root string
	var zip string
	var https bool
	var port int
	flag.StringVar(&backend, "backend", "", "fallback URL")
	flag.StringVar(&root, "root", ".", "directory or zip archive prefix")
	flag.StringVar(&zip, "zip", "", "zip archive to serve files from")
	flag.BoolVar(&https, "https", false, "use HTTPS (self-signed certificate)")
	flag.IntVar(&port, "port", 8090, "port number")
	flag.Parse()

	var scheme string
	if https {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("🐮 listening on %s://%s\n", scheme, host)

	var handler http.Handler = nil

	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Println("invalid backend URL: ", err.Error())
			return
		}
		fmt.Printf("> serving from backend: %s\n", backendURL)
		handler = ko.NewProxyMiddleware(*backendURL)(handler)
	}

	if zip != "" {
		fmt.Printf("> serving from archive: %s@%s\n", zip, root)
		handler = ko.NewZIPMiddleware(zip, root)(handler)
	} else if root != "" {
		fmt.Printf("> serving from directory: %s\n", root)
		handler = ko.NewStaticMiddleware(root)(handler)
	}

	if handler == nil {
		panic("🐮 I'm useless without a handler")
	}

	var err error
	if https {
		err = ko.ListenAndServeTLS(host, handler)
	} else {
		err = http.ListenAndServe(host, handler)
	}
	if err != nil {
		fmt.Println("Server failed with: ", err.Error())
	}
}
