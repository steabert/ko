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

	stack := make([]string, 0, 2)
	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Println("invalid backend URL: ", err.Error())
			return
		}
		stack = append(stack, fmt.Sprintf("%s", backendURL))
		handler = ko.NewProxyMiddleware(*backendURL)(handler)
	}

	if public != "" {
		stack = append(stack, fmt.Sprintf("%s", public))
		handler = ko.NewStaticMiddleware(public)(handler)
	}

	// TODO: this should be checked by the args
	if handler == nil {
		panic("🐮 I'm useless without a handler")
	}

	fmt.Println("🐮 Serving from (in order of priority) ...")
	for i := len(stack); i > 0; i-- {
		fmt.Println("->", stack[i-1])
	}

	var scheme string
	if secure {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("listening on %s://%s ...\n", scheme, host)
	if secure {
		//TODO: include these somehow in the code as constants (then we can run this from anywhere)
		http.ListenAndServeTLS(host, "server.cert", "server.key", handler)
	} else {
		http.ListenAndServe(host, handler)
	}
}
