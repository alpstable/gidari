//go:build utest

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
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alpstable/gidari/internal/web/auth"
	"github.com/alpstable/gidari/tools"
	"golang.org/x/time/rate"
)

func TestFetchWithBasicAuth(t *testing.T) {
	t.Parallel()

	t.Run("authorization success", func(t *testing.T) {
		t.Parallel()

		const username = "test@email.com"
		const password = "test"

		testServer := createTestServerWithBasicAuth(username, password)
		defer testServer.Close()

		ctx := context.Background()

		tripper := auth.NewBasic()
		tripper.SetEmail(username)
		tripper.SetPassword(password)
		tripper.SetURL(testServer.URL)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
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

		testServer := createTestServerWithBasicAuth(username, password)
		defer testServer.Close()

		for _, tcase := range []struct {
			username, password string
		}{
			{username: "wrong", password: "wrong"},
			{username: username},
			{password: password},
			{username: "", password: ""},
		} {
			ctx := context.Background()

			tripper := auth.NewBasic()
			tripper.SetEmail(tcase.username)
			tripper.SetPassword(tcase.password)
			tripper.SetURL(testServer.URL)

			client, err := NewClient(ctx, tripper)
			if err != nil {
				t.Fatalf("error creating client: %v", err)
			}

			uri, err := url.Parse(testServer.URL)
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

		testServer := createTestServerWithBasicAuth(username, password)
		defer testServer.Close()

		ctx := context.Background()

		// Don't set url for tripper.
		tripper := auth.NewBasic()
		tripper.SetEmail(username)
		tripper.SetPassword(password)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
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

func TestFetchWithAuth2(t *testing.T) {
	t.Parallel()

	t.Run("authorization success", func(t *testing.T) {
		t.Parallel()

		const bearer = "AbCd1234"

		testServer := createTestServerWithOAuth2(bearer)
		defer testServer.Close()

		ctx := context.Background()

		tripper := auth.NewAuth2()
		tripper.SetBearer(bearer)
		tripper.SetURL(testServer.URL)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
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

		const bearer = "AbCd1234"

		testServer := createTestServerWithOAuth2(bearer)
		defer testServer.Close()

		for _, tcase := range []struct {
			bearer string
		}{
			{bearer: ""},
			{bearer: "wrong"},
		} {
			ctx := context.Background()

			tripper := auth.NewAuth2()
			tripper.SetBearer(tcase.bearer)
			tripper.SetURL(testServer.URL)

			client, err := NewClient(ctx, tripper)
			if err != nil {
				t.Fatalf("error creating client: %v", err)
			}

			uri, err := url.Parse(testServer.URL)
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

	t.Run("empty url only in auth2 tripper", func(t *testing.T) {
		t.Parallel()

		const bearer = "AbCd1234"

		testServer := createTestServerWithOAuth2(bearer)
		defer testServer.Close()

		ctx := context.Background()

		// Don't set url for tripper.
		tripper := auth.NewAuth2()
		tripper.SetBearer(bearer)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
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

func TestFetchWithAPIKey(t *testing.T) {
	t.Parallel()

	t.Run("authorization success", func(t *testing.T) {
		t.Parallel()

		const key = "apikey"
		const passphrase = "passphrase"
		const secret = "secret01"

		testServer := createTestServerWithAPIKey(key, passphrase, secret)
		defer testServer.Close()

		ctx := context.Background()

		tripper := auth.NewAPIKey()
		tripper.SetKey(key)
		tripper.SetPassphrase(passphrase)
		tripper.SetSecret(secret)
		tripper.SetURL(testServer.URL)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		for _, tcase := range []struct {
			method string
		}{
			{method: http.MethodGet},
			{method: http.MethodPost},
			{method: http.MethodPut},
			{method: http.MethodDelete},
		} {
			_, err = Fetch(ctx, &FetchConfig{
				C:           client,
				Method:      tcase.method,
				URL:         uri,
				RateLimiter: rate.NewLimiter(1, 1),
			})
			if err != nil {
				t.Fatalf("fetch error: %v", err)
			}
		}
	})

	t.Run("authorization failed", func(t *testing.T) {
		t.Parallel()

		const key = "apikey"
		const passphrase = "passphrase"
		const secret = "secret01"

		testServer := createTestServerWithAPIKey(key, passphrase, secret)
		defer testServer.Close()

		for _, tcase := range []struct {
			key, passphrase, secret string
		}{
			{key: "wrong", passphrase: passphrase, secret: secret},
			{key: key, passphrase: "wrong", secret: secret},
			{key: key, passphrase: passphrase, secret: "wrong002"},
			{key: "", passphrase: "", secret: ""},
		} {
			ctx := context.Background()

			tripper := auth.NewAPIKey()
			tripper.SetKey(tcase.key)
			tripper.SetPassphrase(tcase.passphrase)
			tripper.SetSecret(tcase.secret)
			tripper.SetURL(testServer.URL)

			client, err := NewClient(ctx, tripper)
			if err != nil {
				t.Fatalf("error creating client: %v", err)
			}

			uri, err := url.Parse(testServer.URL)
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

	t.Run("empty url only in api key tripper", func(t *testing.T) {
		t.Parallel()

		const key = "apikey"
		const passphrase = "passphrase"
		const secret = "secret01"

		testServer := createTestServerWithAPIKey(key, passphrase, secret)
		defer testServer.Close()

		ctx := context.Background()

		// Don't set url for tripper.
		tripper := auth.NewAPIKey()
		tripper.SetKey(key)
		tripper.SetPassphrase(passphrase)
		tripper.SetSecret(secret)

		client, err := NewClient(ctx, tripper)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}

		uri, err := url.Parse(testServer.URL)
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

// createTestServerWithBasicAuth is a helper that creates a httptest.Server with a handler that has basic auth.
func createTestServerWithBasicAuth(username, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		reqUsername, reqPassword, ok := req.BasicAuth()
		if !ok || reqUsername != username || reqPassword != password {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
}

// createTestServerWithOAuth2 is a helper that creates a httptest.Server with a handler that has OAuth 2.
func createTestServerWithOAuth2(bearer string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		authHeader := strings.Split(req.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) == 2 && authHeader[1] == bearer { // authHeader[1] contains token.
			writer.WriteHeader(http.StatusOK)

			return
		}
		writer.WriteHeader(http.StatusUnauthorized)
	}))
}

// createTestServerWithAPIKey is a helper that creates a httptest.Server with a handler that has api key auth.
func createTestServerWithAPIKey(key, passphrase, secret string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		reqContentType := req.Header.Get("content-type")
		if reqContentType != "application/json" {
			writer.WriteHeader(http.StatusUnsupportedMediaType)

			return
		}

		reqKey := req.Header.Get("cb-access-key")
		reqPassphrase := req.Header.Get("cb-access-passphrase")
		reqSign := req.Header.Get("cb-access-sign")
		if reqKey != key || reqPassphrase != passphrase || reqSign == "" {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		// The httptest.Server handler has a trailing slash.
		// We have to remove it because otherwise the paths will be different,
		// and we will generate the wrong signature.
		req.URL.Path = strings.TrimRight(req.URL.Path, "/")

		// generate sign
		msg := tools.NewHTTPMessage(req, req.Header.Get("cb-access-timestamp"))

		sign, err := msg.Sign(secret)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		// compare signatures
		if reqSign != sign {
			writer.WriteHeader(http.StatusUnauthorized)

			return
		}

		writer.WriteHeader(http.StatusOK)
	}))
}
