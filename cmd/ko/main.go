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
	var archive string
	var secure bool
	var port int
	flag.StringVar(&backend, "backend", "", "fallback URL")
	flag.StringVar(&root, "root", ".", "location of files to host (dir or path into zip archive)")
	flag.StringVar(&archive, "archive", "", "archive to serve files from")
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

	if archive != "" {
		stack = append(stack, fmt.Sprintf("%s@%s", archive, root))
		handler = ko.NewZIPMiddleware(archive, root)(handler)
	} else if root != "" {
		stack = append(stack, fmt.Sprintf("%s", root))
		handler = ko.NewStaticMiddleware(root)(handler)
	}

	if handler == nil {
		panic("🐮 I'm useless without a handler")
	}

	fmt.Println("🐮 Serving from (in order of priority) ...")
	for i := len(stack); i > 0; i-- {
		fmt.Println(" > ", stack[i-1])
	}

	var scheme string
	if secure {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := fmt.Sprintf("localhost:%d", port)

	fmt.Printf("listening on %s://%s ...\n", scheme, host)

	var err error
	if secure {
		err = ko.ListenAndServeTLS(host, handler)
	} else {
		err = http.ListenAndServe(host, handler)
	}
	if err != nil {
		fmt.Println("Server failed with: ", err.Error())
	}
}
