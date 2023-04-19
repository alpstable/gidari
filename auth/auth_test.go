// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
)

const (
	krakenTestURL   = "https://api.kraken.com/"
	coinbaseTestURL = "https://api-public.sandbox.exchange.coinbase.com"
)

type roundTrip struct {
	rtripper RoundTrip
}

// RoundTrip is a interfaces the http.Transport interface for auth integration
// tests.
func (r *roundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.rtripper(req)
}

func newTestKrakenRoundTrip(t *testing.T, key, secret string) *roundTrip {
	t.Helper()

	rtripper, err := NewKrakenRoundTrip(key, secret)
	if err != nil {
		t.Fatal(err)
	}

	return &roundTrip{rtripper: rtripper}
}

func newTestCoinbaseRoundTrip(t *testing.T, key, secret, passphrase string) *roundTrip {
	t.Helper()

	rtripper, err := NewCoinbaseRoundTrip(key, secret, passphrase)
	if err != nil {
		t.Fatal(err)
	}

	return &roundTrip{rtripper: rtripper}
}

func TestNewKrakenRoundTripper(t *testing.T) {
	t.Parallel()

	krakenAPIKey := os.Getenv("KRAKEN_API_KEY")
	krakenAPISecret := os.Getenv("KRAKEN_API_SECRET")

	coinbaseAPISecret := os.Getenv("COINBASE_API_SECRET")
	coinbaseAPIKey := os.Getenv("COINBASE_API_KEY")
	coinbaseAPIPassphrase := os.Getenv("COINBASE_API_PASSPHRASE")

	tests := []struct {
		name      string
		method    string
		api       string
		path      string
		roundTrip *roundTrip
	}{
		{
			name:      "test kraken public api",
			method:    http.MethodGet,
			api:       krakenTestURL,
			path:      "/0/public/Time",
			roundTrip: newTestKrakenRoundTrip(t, krakenAPIKey, krakenAPISecret),
		},
		{
			name:      "test kraken private api",
			method:    http.MethodPost,
			api:       krakenTestURL,
			path:      "/0/private/Balance",
			roundTrip: newTestKrakenRoundTrip(t, krakenAPIKey, krakenAPISecret),
		},
		{
			name:      "test coinbase private api",
			method:    http.MethodGet,
			api:       coinbaseTestURL,
			path:      "/accounts",
			roundTrip: newTestCoinbaseRoundTrip(t, coinbaseAPIKey, coinbaseAPISecret, coinbaseAPIPassphrase),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Create a new client with the CoinbaseRoundTripper.
			client := &http.Client{Transport: test.roundTrip}

			// Create a new request for /accounts.
			u, err := url.JoinPath(test.api, test.path)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(context.Background(), test.method, u, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Send the request.
			rsp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			// Check that the response is not nil.
			if rsp == nil {
				t.Fatal("response is nil")
			}

			// Check that the response status code is 200.
			if rsp.StatusCode != http.StatusOK {
				var buf bytes.Buffer
				if _, err := io.Copy(&buf, rsp.Body); err != nil {
					t.Fatalf("failed to copy response body: %v", err)
				}

				t.Fatalf("response status code is not 200: %d: %s", rsp.StatusCode, buf.String())
			}

			// Close the response body.
			if err := rsp.Body.Close(); err != nil {
				t.Fatalf("failed to close response body: %v", err)
			}
		})
	}
}
