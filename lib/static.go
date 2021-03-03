package lib

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// NewStaticMiddleware creates a router that:
//  - serves static content from file if available
//  - passes request to next handler if not
// Note that this middleware terminates if the file was found!
func NewStaticMiddleware(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := path.Join(root, r.URL.Path)

			// Normalize directory routes to index.html
			info, err := os.Stat(filepath.FromSlash(p))
			if err == nil && info.IsDir() {
				p = path.Join(p, "index.html")
			}

			encodings := strings.Split(r.Header.Get("Accept-Encoding"), ",")
			for _, enc := range encodings {
				switch strings.Trim(enc, "") {
				case "gzip":
					fp := filepath.FromSlash(p + ".gz")
					info, err := os.Stat(fp)
					if err != nil {
						break
					}
					f, err := os.Open(fp)
					if err != nil {
						http.Error(w, "", http.StatusForbidden)
						return
					}
					w.Header().Add("Content-Encoding", "gzip")
					http.ServeContent(w, r, p, info.ModTime(), f)
					return
				}
			}

			fp := filepath.FromSlash(p)
			if info, err = os.Stat(fp); err == nil {
				f, err := os.Open(fp)
				if err != nil {
					http.Error(w, "", http.StatusForbidden)
					return
				}
				http.ServeContent(w, r, p, info.ModTime(), f)
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
