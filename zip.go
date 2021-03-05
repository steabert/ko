package ko

import (
	"archive/zip"
	"io/ioutil"
	"net/http"
)

func NewZIPMiddleware(zipPath string) func(http.Handler) http.Handler {
	files := map[string]*zip.File{}
	z, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	for _, file := range z.File {
		files["/"+file.Name] = file
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			f, ok := files[r.URL.Path]
			if !ok {
				if next != nil {
					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "File not found", http.StatusNotFound)
					return
				}
			}
			b, err := readAll(f)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
			}
			w.Write(b)
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
