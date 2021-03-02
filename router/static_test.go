package router_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/steabert/ko/router"
)

type Spy struct {
	Called bool
}

func (s Spy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Called = true
}

func TestCanRoute(t *testing.T) {
	fallback := Spy{Called: false}
	router := router.NewFileRouter(".")
	handler := router(fallback)
	handler.ServeHTTP(http.ResponseWriter{}, &http.Request{URL: &url.URL{Path: "static.go"}})
	if !r.CanRoute("static.go") {
		t.Fatal("expected files but not found")
	}
}

func TestCantRoute(t *testing.T) {
	r := router.NewFileRouter(".")
	if r.CanRoute("missing.txt") {
		t.Fatal("expected file to be missing")
	}
}
