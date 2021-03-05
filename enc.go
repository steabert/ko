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

var knownEncSuffix = map[string]string{
	"gzip":     ".gz",
	"br":       ".br",
	"identity": "",
}

// AcceptedEncodings extracts accepted values from an Accept- header
func AcceptedEncodings(header http.Header, key string) []string {
	accepts := make([]string, 0, 8)

	for _, headerValue := range header.Values(key) {
		if headerValue == "" {
			continue
		}

		acceptValues := strings.Split(headerValue, comma)
		for _, value := range acceptValues {
			parts := strings.Split(value, semicolon)
			if len(parts) >= 1 {
				enc := strings.Trim(parts[0], whitespace)

				// skip identity since it will be added explicitly
				if enc == "identity" {
					continue
				}

				accepts = append(accepts, enc)
			}
		}
	}

	// always add identity to allow easy generic fallback
	accepts = append(accepts, "identity")

	return accepts
}
