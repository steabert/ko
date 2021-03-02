package router

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

// StaticRouter serves files relative to a root directory
type StaticRouter struct {
	root      string
}

// NewStaticRouter creates a router that can serve static content
func NewStaticRouter(root string) *StaticRouter {
	return &StaticRouter{root: root}

}

// CanRoute confirms if a path can be served by the router
func (router StaticRouter) CanRoute(route string) bool {
	fp := path.Join(router.root, route)
	s, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// ServeHTTP
func (router StaticRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !router.CanRoute(r.URL.Path) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	fp := path.Join(router.root, r.URL.Path)
	// If file path exists, serve from file system
	fmt.Println("serving: ", fp)
	http.ServeFile(w, r, fp)
	return
}
