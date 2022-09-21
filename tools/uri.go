package tools

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrParsingURL is returned when there is an error parsing the url.
	ErrParsingURL = fmt.Errorf("error parsing url")
)

// SplitURLPath will return the endpoint parts from the request.
func SplitURLPath(req http.Request) []string {
	parts := strings.Split(strings.TrimPrefix(req.URL.EscapedPath(), "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return []string{}
	}
	return parts
}

// ParseDBTableFromURL will return the table name from the request.
func ParseDBTableFromURL(req http.Request) (string, error) {
	endpointParts := SplitURLPath(req)
	if len(endpointParts) == 0 {
		return "", ErrParsingURL
	}

	table := endpointParts[len(endpointParts)-1]
	return table, nil
}
