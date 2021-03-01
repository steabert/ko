package router_test

import (
	"github.com/steabert/ko/router"
	"testing"
)

func TestCanRoute(t *testing.T) {
	r := router.NewStaticRouter(".")
	if !r.CanRoute("static.go") {
		t.Fatal("expected files but not found")
	}
}

func TestCantRoute(t *testing.T) {
	r := router.NewStaticRouter(".")
	if r.CanRoute("missing.txt") {
		t.Fatal("expected file to be missing")
	}
}