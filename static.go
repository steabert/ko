package ko

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// TODO: instead of using a root dir directly, give the main
// middleware the possibility to look up and read file content,
// provided by either a dir or zip archive reader.

// NewStaticMiddleware creates a router that:
//  - serves static content from file if available
//  - passes request to next handler if not
// Note that this middleware terminates if the file was found!
func NewStaticMiddleware(root string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := path.Clean(r.URL.Path)
			p := path.Join(root, name)

			// Normalize directory routes to index.html
			info, err := os.Stat(filepath.FromSlash(p))
			if err == nil && info.IsDir() {
				p = path.Join(p, "index.html")
			}

			encodings := ParseAccept(r.Header, "Accept-Encoding")
			fmt.Println("encodings", encodings)
			for _, enc := range encodings {
				suffix, exists := knownEncSuffix[enc]
				if !exists {
					continue
				}

				fp := filepath.FromSlash(p + suffix)
				if info, err := os.Stat(fp); err == nil {
					f, err := os.Open(fp)
					if err != nil {
						http.Error(w, "", http.StatusForbidden)
						return
					}
					w.Header().Add("Content-Encoding", enc)
					http.ServeContent(w, r, name, info.ModTime(), f)
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
				http.ServeContent(w, r, name, info.ModTime(), f)
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
