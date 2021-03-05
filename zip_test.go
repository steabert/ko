package ko_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/steabert/ko"
)

func TestReadArchiveExistingFile(t *testing.T) {
	middleware := ko.NewZIPMiddleware("../testdir/tmp.zip")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/tmp/file.html")
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != 200 {
		t.Fatal("expected file.html to be found")
	}
}

func TestReadArchiveNoneExistingFile(t *testing.T) {
	middleware := ko.NewZIPMiddleware("../testdir/tmp.zip")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/tmp/missing.html")
	if err != nil {
		t.Fatal(err)
	}
	if rsp.StatusCode != 404 {
		t.Fatal("expected missing.html to not be found")
	}
}
