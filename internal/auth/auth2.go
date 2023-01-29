// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package auth

import (
	"fmt"
	"net/http"
	"net/url"
)

// Auth2 is an OAuth2 http transport.
type Auth2 struct {
	bearer string
	url    *url.URL
}

// NewAuth2 will return an OAuth2 http transport.
func NewAuth2() *Auth2 {
	return new(Auth2)
}

// SetBearer will set the bearer field on Auth2.
func (auth *Auth2) SetBearer(val string) *Auth2 {
	auth.bearer = val

	return auth
}

// SetURL will set the key field on APIKey.
func (auth *Auth2) SetURL(u string) *Auth2 {
	auth.url, _ = url.Parse(u)

	return auth
}

// RoundTrip authorizes the request with a signed OAuth1 Authorization header using the author and TokenSource.
func (auth *Auth2) RoundTrip(req *http.Request) (*http.Response, error) {
	if auth.url == nil {
		return nil, ErrURLRequired
	}

	req.URL.Scheme = auth.url.Scheme
	req.URL.Host = auth.url.Host
	req.Header.Set(authorizationHeaderParam, fmt.Sprintf("%s %s", bearerHeaderPrefix, auth.bearer))

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	return rsp, nil
}
