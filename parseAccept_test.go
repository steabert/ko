package ko

import (
	"fmt"
	"net/http"
	"testing"
)

func TestParseQValues(t *testing.T) {
	header := make(http.Header)
	header.Add("Accept-Encoding", "gzip;q=0.8, br;q=0.5 , indentity")

	accepts := parseAccept(header, "Accept-Encoding")
	fmt.Println(accepts)
}
