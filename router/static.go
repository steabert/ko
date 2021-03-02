package router

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

// NewFileRouter creates a router that can serve static content
func NewFileRouter(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fp := path.Join(root, r.URL.Path)
			s, err := os.Stat(fp)
			if err == nil && !s.IsDir() {
				fmt.Println("serving: ", fp)
				http.ServeFile(w, r, fp)
			}
			next.ServeHTTP(w, r)
		})
	}
}
