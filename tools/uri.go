// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package tools

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ErrParsingURL is returned when there is an error parsing the url.
var ErrParsingURL = fmt.Errorf("error parsing url")

// SplitURL will return the endpoint parts from the url.
func SplitURL(url *url.URL) []string {
	parts := strings.Split(strings.TrimPrefix(url.EscapedPath(), "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return []string{}
	}

	return parts
}

// ParseDBTableFromURL will return the table name from the url.
func ParseDBTableFromURL(url *url.URL) (string, error) {
	parts := SplitURL(url)
	if len(parts) == 0 {
		return "", ErrParsingURL
	}

	return parts[len(parts)-1], nil
}

// SplitURLFromRequest will return the endpoint parts from the request.
func SplitURLFromRequest(req http.Request) []string {
	return SplitURL(req.URL)
}

// ParseDBTableFromRequest will return the table name from the request.
func ParseDBTableFromRequest(req http.Request) (string, error) {
	return ParseDBTableFromURL(req.URL)
}
