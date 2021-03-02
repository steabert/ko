package router_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/steabert/ko/router"
)

type CallRouter struct {
	Called *bool
}

func (s CallRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*s.Called = true
}

func TestNonExistent(t *testing.T) {
	middleware := router.NewFileRouter(".")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/nonexistent")
	if err != nil {
		log.Fatal(err)
	}
	if rsp.StatusCode != 404 {
		log.Fatal("expected file not found")
	}
}

func TestExistent(t *testing.T) {
	middleware := router.NewFileRouter(".")

	ts := httptest.NewServer(middleware(nil))
	rsp, err := http.Get(ts.URL + "/static.go")
	if err != nil {
		log.Fatal(err)
	}
	if rsp.StatusCode != 200 {
		log.Fatal("expected file not found")
	}
}
