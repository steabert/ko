package ko

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"net/http"
	"path"
)

// NewZIPMiddleware serves files from within a ZIP archive
func NewZIPMiddleware(zipPath, prefix string) func(http.Handler) http.Handler {
	files := map[string]*zip.File{}
	z, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	for _, file := range z.File {
		files[file.Name] = file
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name := path.Clean(r.URL.Path)
			p := path.Join(prefix, name)

			// Normalize directory routes to index.html
			if _, ok := files[p+"/"]; ok {
				p = path.Join(p, "index.html")
			}

			encodings := AcceptedEncodings(r.Header, "Accept-Encoding")
			for _, enc := range encodings {
				suffix, exists := knownEncSuffix[enc]
				if !exists {
					continue
				}

				f, ok := files[p+suffix]
				if !ok {
					continue
				}

				w.Header().Add("Content-Encoding", enc)
				b, err := readAll(f)
				if err != nil {
					http.Error(w, "Internal error", http.StatusInternalServerError)
					return
				}
				http.ServeContent(w, r, name, f.ModTime(), bytes.NewReader(b))
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

// readAll is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func readAll(file *zip.File) ([]byte, error) {
	fc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		return nil, err
	}

	return content, nil
}
