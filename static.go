package ko

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// NewStaticMiddleware serves files from a local folder
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

			encodings := AcceptedEncodings(r.Header, "Accept-Encoding")
			for _, enc := range encodings {
				suffix, exists := knownEncSuffix[enc]
				if !exists {
					continue
				}

				fp := filepath.FromSlash(p + suffix)
				if info, err := os.Stat(fp); err == nil {
					f, err := os.Open(fp)
					if err != nil {
						http.Error(w, "Internal error", http.StatusInternalServerError)
						return
					}
					w.Header().Add("Content-Encoding", enc)
					fmt.Println("serving: ", w)
					http.ServeContent(w, r, name, info.ModTime(), f)
					return
				}
			}

			if next != nil {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "File not found", http.StatusNotFound)
			}
		})
	}
}
