// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package web

import (
	"context"
	"github.com/alpine-hodler/gidari/internal/web/auth"
	"golang.org/x/time/rate"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFetchWithBasicAuth(t *testing.T) {
	t.Parallel()

	t.Run("authorization success", func(t *testing.T) {
		t.Parallel()

		const username = "test@email.com"
		const password = "test"

		ts := createTestServerWithBasicAuth(username, password)
		defer ts.Close()

		ctx := context.Background()

		basicAuth := auth.NewBasic()
		basicAuth.SetEmail(username)
		basicAuth.SetPassword(password)
		basicAuth.SetURL(ts.URL)

		client, err := NewClient(ctx, basicAuth)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		_, err = Fetch(ctx, &FetchConfig{
			C:           client,
			Method:      http.MethodGet,
			URL:         uri,
			RateLimiter: rate.NewLimiter(1, 1),
		})
		if err != nil {
			t.Fatalf("fetch error: %v", err)
		}
	})

	t.Run("authorization failed", func(t *testing.T) {
		t.Parallel()

		const username = "test@email.com"
		const password = "test"

		ts := createTestServerWithBasicAuth(username, password)
		defer ts.Close()

		for _, tcase := range []struct {
			username, password string
		}{
			{username: "wrong", password: "wrong"},
			{username: username},
			{password: password},
			{username: "", password: ""},
		} {
			ctx := context.Background()

			basicAuth := auth.NewBasic()
			basicAuth.SetEmail(tcase.username)
			basicAuth.SetPassword(tcase.password)
			basicAuth.SetURL(ts.URL)

			client, err := NewClient(ctx, basicAuth)
			if err != nil {
				t.Fatalf("error creating client: %v", err)
			}

			uri, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatalf("error parsing url: %v", err)
			}

			_, err = Fetch(ctx, &FetchConfig{
				C:           client,
				Method:      http.MethodGet,
				URL:         uri,
				RateLimiter: rate.NewLimiter(1, 1),
			})
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
		}
	})

	t.Run("empty url only in basic auth tripper", func(t *testing.T) {
		t.Parallel()

		const username = "test@email.com"
		const password = "test"

		ts := createTestServerWithBasicAuth(username, password)
		defer ts.Close()

		ctx := context.Background()

		// Don't set url for tripper
		basicAuth := auth.NewBasic()
		basicAuth.SetEmail(username)
		basicAuth.SetPassword(password)

		client, err := NewClient(ctx, basicAuth)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		_, err = Fetch(ctx, &FetchConfig{
			C:           client,
			Method:      http.MethodGet,
			URL:         uri,
			RateLimiter: rate.NewLimiter(1, 1),
		})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}

// createTestServerWithBasicAuth is a helper that creates a httptest.Server with a handler that has basic authentication.
func createTestServerWithBasicAuth(username, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUsername, reqPassword, ok := r.BasicAuth()
		if !ok || reqUsername != username || reqPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}
