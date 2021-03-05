package ko_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steabert/ko"
)

type CallRouter struct {
	Called *bool
}

func (s CallRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*s.Called = true
}

func TestNonExistent(t *testing.T) {
	middleware := ko.NewStaticMiddleware(".")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL + "/nonexistent")
	if err != nil {
		log.Fatal(err)
	}

	if rsp.StatusCode != 404 {
		log.Fatal("expected file not found")
	}
}

func TestExistent(t *testing.T) {
	middleware := ko.NewStaticMiddleware(".")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL + "/static.go")
	if err != nil {
		log.Fatal(err)
	}

	if rsp.StatusCode != 200 {
		log.Fatal("expected file not found")
	}
}

func TestIndexExist(t *testing.T) {
	middleware := ko.NewStaticMiddleware("testdir")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	if rsp.StatusCode != 200 {
		log.Fatal("expected index.html to be found")
	}
}

func TestIndexNoneExist(t *testing.T) {
	middleware := ko.NewStaticMiddleware(".")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	if rsp.StatusCode != 404 {
		log.Fatal("expected index.html not to be found")
	}
}

func TestContentType(t *testing.T) {
	middleware := ko.NewStaticMiddleware("testdir")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL + "/index.html")
	if err != nil {
		log.Fatal(err)
	}

	ce := rsp.Header.Get("Content-Type")
	if !strings.Contains(ce, "text/html") {
		t.Fatalf("expected text/html got %s", ce)
	}

	if rsp.StatusCode != 200 {
		log.Fatal("expected file to be found")
	}
}

func TestContentType2(t *testing.T) {
	middleware := ko.NewStaticMiddleware("testdir")

	ts := httptest.NewServer(middleware(nil))
	defer ts.Close()

	rsp, err := http.Get(ts.URL + "/test.js")
	if err != nil {
		log.Fatal(err)
	}

	ct := rsp.Header.Get("Content-Type")
	if !strings.Contains(ct, "javascript") {
		t.Fatalf("expected javascript got %s", ct)
	}

	// The "Content-Encoding" header will be gone,
	// which is explained at the `Uncompressed` field:
	// https://golang.org/pkg/net/http/#Response
	if !rsp.Uncompressed {
		t.Fatalf("expected uncompressed response")
	}

	if rsp.StatusCode != 200 {
		log.Fatal("expected file to be found")
	}
}
