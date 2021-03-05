package ko_test

import (
	"net/http"
	"testing"

	"github.com/steabert/ko"
)

func TestAcceptedEncodings(t *testing.T) {
	header := make(http.Header)
	header.Add("Accept-Encoding", "gzip;q=0.8, br;q=0.5 , identity")

	accepts := ko.AcceptedEncodings(header, "Accept-Encoding") // [gzip br identity]
	if accepts[0] != "gzip" {
		t.Fatalf("expected: %s, actual: %s", "gzip", accepts[0])
	}
	if accepts[1] != "br" {
		t.Fatalf("expected: %s, actual: %s", "br", accepts[1])
	}
	if accepts[2] != "identity" {
		t.Fatalf("expected: %s, actual: %s", "identity", accepts[2])
	}
}
