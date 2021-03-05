package ko

import (
	"net/http"
	"strings"
)

const (
	comma      = ","
	semicolon  = ";"
	whitespace = " "
)

func parseAccept(header http.Header, key string) []string {
	accepts := make([]string, 0, 8)

	for _, headerValue := range header.Values(key) {
		if headerValue == "" {
			continue
		}

		acceptValues := strings.Split(headerValue, comma)
		for _, value := range acceptValues {
			parts := strings.Split(value, semicolon)
			if len(parts) >= 1 {
				accepts = append(accepts, strings.Trim(parts[0], whitespace))
			}
		}
	}

	return accepts
}