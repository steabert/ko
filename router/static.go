package router

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
)

// StaticRouter serves files relative to a root directory
type StaticRouter struct {
	root      string
	filepaths map[string]bool
}

// NewStaticRouter creates a router that can serve static content
func NewStaticRouter(root string) (*StaticRouter, error) {
	router := StaticRouter{root: root, filepaths: make(map[string]bool)}
	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			router.filepaths[p] = true
		}
		return nil
	})
	return &router, nil
}

// CanRoute confirms if a path can be served by the router
func (router StaticRouter) CanRoute(route string) bool {
	fp := path.Join(router.root, route)
	return router.filepaths[fp]
}

// ServeHTTP
func (router StaticRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fp := path.Join(router.root, r.URL.Path)
	// If file path exists, serve from file system
	if !router.filepaths[fp] {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	fmt.Println("serving: ", fp)
	http.ServeFile(w, r, fp)
	return
}
