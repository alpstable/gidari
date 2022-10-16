//go:build utest

// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPMessage(t *testing.T) {
	t.Parallel()

	t.Run("http methods", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		for _, tcase := range []struct {
			method string
		}{
			{method: http.MethodGet},
			{method: http.MethodPost},
			{method: http.MethodPut},
			{method: http.MethodDelete},
		} {
			req, err := http.NewRequestWithContext(ctx, tcase.method, "https://foo/path", nil)
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			timestamp := currentTimestamp()

			message := NewHTTPMessage(req, timestamp)

			expected := HTTPMessage(fmt.Sprintf("%s%s/path", timestamp, tcase.method))
			if message != expected {
				t.Fatalf("expected %s, got %s", expected, message)
			}
		}
	})

	t.Run("path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://foo/path", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		message := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path", timestamp))
		if message != expected {
			t.Fatalf("expected %s, got %s", expected, message)
		}
	})

	t.Run("path with query params", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://foo/path?param=1", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		message := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path?param=1", timestamp))
		if message != expected {
			t.Fatalf("expected %s, got %s", expected, message)
		}
	})

	t.Run("body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := strings.NewReader(`{"username":"john doe"}`)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://foo", body)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		message := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET{\"username\":\"john doe\"}", timestamp))
		if message != expected {
			t.Fatalf("expected %s, got %s", expected, message)
		}
	})

	t.Run("path and body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := strings.NewReader(`{"username":"john doe"}`)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://foo/path?param=1", body)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		message := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path?param=1{\"username\":\"john doe\"}", timestamp))
		if message != expected {
			t.Fatalf("expected %s, got %s", expected, message)
		}
	})

	t.Run("empty path and body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://foo", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		message := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET", timestamp))
		if message != expected {
			t.Fatalf("expected %s, got %s", expected, message)
		}
	})
}

func TestHTTPMessage_Sign(t *testing.T) {
	t.Parallel()

	t.Run("sign", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			message HTTPMessage
			secret  string
		}{
			{message: HTTPMessage("first"), secret: "1234"},
			{message: HTTPMessage("second"), secret: "12341234"},
			{message: HTTPMessage("third"), secret: "123412341234"},
		} {
			signature, err := tcase.message.Sign(tcase.secret)
			if err != nil {
				t.Fatalf("signature creation error: %v", err)
			}
			if signature == "" {
				t.Fatal("expected signature, got empty string")
			}
			if _, err = base64.StdEncoding.DecodeString(signature); err != nil {
				t.Fatalf("signature decoding error: %v", err)
			}
		}
	})

	t.Run("invalid secret", func(t *testing.T) {
		t.Parallel()

		message := HTTPMessage("message")

		for _, tcase := range []struct {
			secret string
		}{
			{secret: "1"},
			{secret: "12345"},
			{secret: "@@@@"},
			{secret: "\\\\\\\\"},
		} {
			if _, err := message.Sign(tcase.secret); err == nil {
				t.Fatalf("expected error, got nil")
			}
		}
	})
}

// currentTimestamp is a helper that returns the formatted current time as a string.
func currentTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
