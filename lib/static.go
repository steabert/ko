package lib

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

// NewStaticMiddleware creates a router that:
//  - serves static content from file if available
//  - passes request to next handler if not
// Note that this middleware terminates if the file was found!
func NewStaticMiddleware(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fp := path.Join(root, r.URL.Path)

			s, err := os.Stat(fp)
			if err == nil && !s.IsDir() {
				fmt.Println("serving: ", fp)
				http.ServeFile(w, r, fp)
				return
			}

			if next != nil {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "File not found", http.StatusNotFound)
			}
		})
	}
}
