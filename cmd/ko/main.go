package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/steabert/ko/router"
)

func main() {
	var public string
	var backend string
	flag.StringVar(&public, "public", "", "directory or archive to serve files from")
	flag.StringVar(&backend, "backend", "", "fallback backend")
	flag.Parse()

	staticRouter, err := router.NewStaticRouter(public)
	if err != nil {
		fmt.Println("unable to initialize static router: ", err.Error())
		return
	}
	proxyRouter, err := router.NewProxyRouter(backend)
	if err != nil {
		fmt.Println("unable to initialize proxy router: ", err.Error())
		return
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if staticRouter.CanRoute(r.URL.Path) {
			fmt.Printf("... [%s] => %s\n", public, r.URL.Path)
			staticRouter.ServeHTTP(w, r)
		} else {
			fmt.Printf("... [%s] => %s\n", backend, r.URL.Path)
			proxyRouter.ServeHTTP(w, r)
		}
	})

	host := "127.0.0.1:4080"
	fmt.Printf("🐮 listening on %s, what would you like me to serve? ...\n", host)
	http.ListenAndServe(host, handler)
}
