package ko_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/steabert/ko"
)

func TestReadArchiveExistingFile(t *testing.T) {
	middleware := ko.NewZIPMiddleware("testdir.zip", "testdir")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != 200 {
		t.Fatal("expected index.html to be found")
	}
}

func TestReadArchiveNoneExistingFile(t *testing.T) {
	middleware := ko.NewZIPMiddleware("testdir.zip", "testdir")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/missing.html")
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != 404 {
		t.Fatal("expected missing.html to not be found")
	}
}
