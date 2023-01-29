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
	"strconv"
	"time"

	"github.com/alpstable/gidari/tools"
)

const (
	// apiKeyTimestmapBase is the base time for calculating the timestamp parameter.
	apiKeyTimestampBase = 10
)

// APIKey is transport for authenticating with an API KEy. API Key authentication should only be used to access your
// own account. If your application requires access to other accounts, do not use API Key. API key authentication
// requires each request to be signed (enhanced security measure). Your API keys should be assigned to access only
// accounts and permission scopes that are necessary for your app to function.
type APIKey struct {
	key        string
	passphrase string
	secret     string
	url        *url.URL
}

// NewAPIKey will return an APIKey authentication transport.
func NewAPIKey() *APIKey {
	return new(APIKey)
}

// SetKey will set the key field on APIKey.
func (auth *APIKey) SetKey(key string) *APIKey {
	auth.key = key

	return auth
}

// SetPassphrase will set the key field on APIKey.
func (auth *APIKey) SetPassphrase(passphrase string) *APIKey {
	auth.passphrase = passphrase

	return auth
}

// SetSecret will set the key field on APIKey.
func (auth *APIKey) SetSecret(secret string) *APIKey {
	auth.secret = secret

	return auth
}

// SetURL will set the key field on APIKey.
func (auth *APIKey) SetURL(u string) *APIKey {
	auth.url, _ = url.Parse(u)

	return auth
}

// RoundTrip authorizes the request with a signed API Key Authorization header.
func (auth *APIKey) RoundTrip(req *http.Request) (*http.Response, error) {
	if auth.url == nil {
		return nil, ErrURLRequired
	}

	var (
		timestamp = strconv.FormatInt(time.Now().Unix(), apiKeyTimestampBase)
		msg       = tools.NewHTTPMessage(req, timestamp)
	)

	sig, err := msg.Sign(auth.secret)
	if err != nil {
		return nil, fmt.Errorf("signature creation error: %w", err)
	}

	req.URL.Scheme = auth.url.Scheme
	req.URL.Host = auth.url.Host

	req.Header.Set("content-type", "application/json")
	req.Header.Add("cb-access-key", auth.key)
	req.Header.Add("cb-access-passphrase", auth.passphrase)
	req.Header.Add("cb-access-sign", sig)
	req.Header.Add("cb-access-timestamp", timestamp)

	rsp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return rsp, nil
}
