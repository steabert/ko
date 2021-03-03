package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/steabert/ko/lib"
)

func main() {
	var public string
	var backend string
	flag.StringVar(&public, "public", "", "directory or archive to serve files from")
	flag.StringVar(&backend, "backend", "", "fallback backend")
	flag.Parse()

	var handler http.Handler = nil

	if backend != "" {
		backendURL, err := url.Parse(backend)
		if err != nil {
			fmt.Println("invalid backend URL: ", err.Error())
			return
		}
		fmt.Println("...enabling reverse proxy: ", backendURL)
		handler = lib.NewProxyMiddleware(*backendURL)(handler)
	}

	if public != "" {
		fmt.Println("...serving from folder: ", public)
		handler = lib.NewStaticMiddleware(public)(handler)
	}

	host := "127.0.0.1:4080"
	fmt.Printf("🐮 listening on %s, what would you like me to serve? ...\n", host)
	http.ListenAndServe(host, handler)
}
