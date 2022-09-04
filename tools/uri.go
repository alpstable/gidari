package tools

import (
	"fmt"
	"net/http"
	"strings"
)

// MongoURI will return the full URI for accessing a mongo database.
func MongoURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("mongodb://%s%s:%s/%s", auth, host, port, database), nil
}

// PostgressURI will return the full URI for accessing a postgres database.
func PostgresURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" || username != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("postgresql://%s%s:%s/%s?sslmode=disable", auth, host, port, database), nil
}

// RedisURI will return the full URI for accessing a redis cache.
func RedisURI(host, username, password, port, database string) (string, error) {
	auth := ""
	if password != "" {
		auth = fmt.Sprintf("%s:%s@", username, password)
	}
	return fmt.Sprintf("redis://%s%s:%s/%s", auth, host, port, database), nil

}

// EndpointPartsFromHTTPRequest will return the endpoint parts from the request.
func EndpointPartsFromHTTPRequest(req http.Request) []string {
	parts := strings.Split(strings.TrimPrefix(req.URL.EscapedPath(), "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return []string{}
	}
	return parts
}

// TableFromHTTPRequest will return the table name from the request.
func TableFromHTTPRequest(req http.Request) (string, error) {
	endpointParts := EndpointPartsFromHTTPRequest(req)
	if len(endpointParts) == 0 {
		return "", fmt.Errorf("no endpoint parts found in url: %s", req.URL)
	}

	table := endpointParts[len(endpointParts)-1]
	return table, nil
}
