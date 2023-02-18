// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"net/http"
	"os"
	"testing"
)

type roundTrip struct {
	rtripper RoundTrip
}

// RoundTrip is a interfaces the http.Transport interface for auth integration
// tests.
func (r *roundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.rtripper(req)
}

func TestNewCoinbaseRoundTripper(t *testing.T) {
	t.Parallel()

	const api = "https://api-public.sandbox.exchange.coinbase.com"

	secret := os.Getenv("COINBASE_API_SECRET")
	key := os.Getenv("COINBASE_API_KEY")
	passphrase := os.Getenv("COINBASE_API_PASSPHRASE")

	// If the environment variabls are not set, we should skip this test.
	if secret == "" || key == "" || passphrase == "" {
		t.Skip("COINBASE_API_SECRET, COINBASE_API_KEY, and " +
			"COINBASE_API_PASSPHRASE must be set to run this test")
	}

	// Create a new CoinbaseRoundTripper.
	rtripper := NewCoinbaseRoundTrip(key, secret, passphrase)

	// Check that the CoinbaseRoundTripper is not nil.
	if rtripper == nil {
		t.Fatal("CoinbaseRoundTripper is nil")
	}

	roundTrip := &roundTrip{rtripper: rtripper}

	// Create a new client with the CoinbaseRoundTripper.
	client := &http.Client{Transport: roundTrip}

	// Create a new request for /accounts.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, api+"/accounts", nil)
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
		t.Fatal("response status code is not 200")
	}

	// Close the response body.
	if err := rsp.Body.Close(); err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}
