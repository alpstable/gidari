// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package web

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/alpine-hodler/gidari/internal/web/auth"
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

		// generate sign
		msg := generageMsg(req, req.Header.Get("cb-access-timestamp"))

		sign, err := generateSign(secret, msg)
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

// generateSign is a helper to generate signature for request.
func generateSign(secret, message string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("error decoding secret: %w", err)
	}

	signature := hmac.New(sha256.New, key)

	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("error writing signature: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}

// generageMsg makes the message to be signed.
func generageMsg(req *http.Request, timestamp string) string {
	postAuthority := strings.TrimSuffix(req.URL.String(), "/")

	return fmt.Sprintf("%s%s%s%s", timestamp, req.Method, postAuthority, string(parsebytes(req)))
}

// parsebytes will return the byte stream for the body.
func parsebytes(req *http.Request) []byte {
	if req.Body == nil {
		return []byte{}
	}

	body, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return body
}
