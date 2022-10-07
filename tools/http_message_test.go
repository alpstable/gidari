package tools

import (
	"context"
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
			req, err := http.NewRequestWithContext(ctx, tcase.method, "https://someurl.com/path", nil)
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			timestamp := currentTimestamp()

			msg := NewHTTPMessage(req, timestamp)

			expected := HTTPMessage(fmt.Sprintf("%s%s/path", timestamp, tcase.method))
			if msg != expected {
				t.Fatalf("expected %s, got %s", expected, msg)
			}
		}
	})

	t.Run("path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://someurl.com/path", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		msg := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path", timestamp))
		if msg != expected {
			t.Fatalf("expected %s, got %s", expected, msg)
		}
	})

	t.Run("path with query params", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://someurl.com/path?param=1", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		msg := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path?param=1", timestamp))
		if msg != expected {
			t.Fatalf("expected %s, got %s", expected, msg)
		}
	})

	t.Run("body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := strings.NewReader(`{"username":"john doe"}`)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://someurl.com", body)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		msg := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET{\"username\":\"john doe\"}", timestamp))
		if msg != expected {
			t.Fatalf("expected %s, got %s", expected, msg)
		}
	})

	t.Run("path and body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		body := strings.NewReader(`{"username":"john doe"}`)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://someurl.com/path?param=1", body)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		msg := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET/path?param=1{\"username\":\"john doe\"}", timestamp))
		if msg != expected {
			t.Fatalf("expected %s, got %s", expected, msg)
		}
	})

	t.Run("empty path and body", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://someurl.com", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}

		timestamp := currentTimestamp()

		msg := NewHTTPMessage(req, timestamp)

		expected := HTTPMessage(fmt.Sprintf("%sGET", timestamp))
		if msg != expected {
			t.Fatalf("expected %s, got %s", expected, msg)
		}
	})
}

// currentTimestamp is a helper that returns the formatted current time as a string.
func currentTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
