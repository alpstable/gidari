// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

// Package auth contains a non-exhaustive list of custom authentication round
// trippers to be used as authentication middleware with a gidari HTTP Service.
package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var errInvalidRoundTripArgs = fmt.Errorf("invalid auth arguments")

// RoundTrip is an HTTP round tripper that acts as a middleware to add
// auth requirements to HTTP requests.
type RoundTrip func(*http.Request) (*http.Response, error)

// NewBasicAuthRoundTrip will return a "RoundTrip" function that can be used as
// a "RoundTrip" function in an "http.RoundTripper" interface to authenticate
// requests that require basic authentication.
func NewBasicAuthRoundTrip(username, password string) (RoundTrip, error) {
	if username == "" || password == "" {
		return nil, errInvalidRoundTripArgs
	}

	return func(req *http.Request) (*http.Response, error) {
		req.SetBasicAuth(username, password)

		rsp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		return rsp, nil
	}, nil
}

// NewCoinbaseRoundTrip will return a "RoundTrip" function that can be used as
// a "RoundTrip" function in an "http.RoundTripper" interface to authenticate
// requests to the Coinbase Cloud API.
func NewCoinbaseRoundTrip(key, secret, passphrase string) (RoundTrip, error) {
	if key == "" || secret == "" || passphrase == "" {
		return nil, errInvalidRoundTripArgs
	}

	return func(req *http.Request) (*http.Response, error) {
		var body []byte
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		requestPath := req.URL.Path
		if req.URL.RawQuery != "" {
			requestPath = fmt.Sprintf("%s?%s", req.URL.Path, req.URL.RawQuery)
		}

		formatBase := 10
		timestamp := strconv.FormatInt(time.Now().Unix(), formatBase)
		msg := fmt.Sprintf("%s%s%s%s", timestamp, req.Method, requestPath, string(body))

		skey, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, fmt.Errorf("error decoding secret: %w", err)
		}

		signature := hmac.New(sha256.New, skey)

		// Don't handle error because hash.Write method never returns an error.
		signature.Write([]byte(msg))
		sig := base64.StdEncoding.EncodeToString(signature.Sum(nil))

		req.Header.Set("content-type", "application/json")
		req.Header.Add("cb-access-key", key)
		req.Header.Add("cb-access-passphrase", passphrase)
		req.Header.Add("cb-access-sign", sig)
		req.Header.Add("cb-access-timestamp", timestamp)

		rsp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		return rsp, nil
	}, nil
}

// NewKrakenRoundTrip will return a "RoundTrip" function that can be used as
// a "RoundTrip" function in an "http.RoundTripper" interface to authenticate
// requests to the Kraken Cloud API.
//
// Kraken uses a custom authentication algorithm that is based on a combination
// of API key, nonce, and message signature. The signature is generated using a
// hash-based message authentication code (HMAC) with SHA-512 as the hash
// function.
//
// When making an API request to Kraken, you will need to include the API key
// and a message signature in the request headers. The message signature is a
// hash of the request data and the nonce using your API secret as the key. The
// nonce is a unique integer value that must be incremented for each request to
// prevent replay attacks.
//
// You can find more information about Kraken's API authentication in their API
// documentation: https://docs.kraken.com/rest/#section/Authentication
func NewKrakenRoundTrip(key, secret string) (RoundTrip, error) {
	if key == "" || secret == "" {
		return nil, errInvalidRoundTripArgs
	}

	return func(req *http.Request) (*http.Response, error) {
		// If the path includes "public" then we don't need to sign the
		// request.
		if strings.Contains(req.URL.Path, "0/public") {
			rsp, err := http.DefaultTransport.RoundTrip(req)
			if err != nil {
				return nil, fmt.Errorf("error making request: %w", err)
			}

			return rsp, nil
		}

		var body []byte
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// Parse the body bytes into url.Values
		values, _ := url.ParseQuery(string(body))
		values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

		sha := sha256.New()
		sha.Write([]byte(values.Get("nonce") + values.Encode()))
		shasum := sha.Sum(nil)

		b64DecodedSecret, _ := base64.StdEncoding.DecodeString(secret)

		mac := hmac.New(sha512.New, b64DecodedSecret)
		mac.Write(append([]byte(req.URL.Path), shasum...))
		macsum := mac.Sum(nil)

		signature := base64.StdEncoding.EncodeToString(macsum)

		req.Header.Add("API-Sign", signature)
		req.Header.Add("API-Key", key)

		// Set the body to the url.Values encoded body.
		req.Body = io.NopCloser(bytes.NewBufferString(values.Encode()))

		rsp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		return rsp, nil
	}, nil
}
