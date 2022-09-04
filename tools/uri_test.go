package tools

import (
	"net/http"
	"net/url"
	"testing"
)

func TestURLHelpers(t *testing.T) {
	t.Run("TableFromHTTPRequest", func(t *testing.T) {
		t.Run("should return the table name from the request", func(t *testing.T) {
			u, _ := url.Parse("http://test.com/v1/tables/test")
			req := http.Request{URL: u}
			table, err := TableFromHTTPRequest(req)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if table != "test" {
				t.Errorf("expected table to be %s, got %s", "test", table)
			}
		})

		t.Run("should return an error when the table name is not present", func(t *testing.T) {
			u, _ := url.Parse("http://test")
			req := http.Request{URL: u}
			_, err := TableFromHTTPRequest(req)
			if err == nil {
				t.Error("expected error to not be nil")
			}
		})
	})

	t.Run("EndpointPartsFromHTTPRequest", func(t *testing.T) {
		t.Run("should return the parts of the endpoint", func(t *testing.T) {
			u, _ := url.Parse("http://test.com/v1/tables/test")
			req := http.Request{URL: u}
			parts := EndpointPartsFromHTTPRequest(req)
			if parts[0] != "v1" {
				t.Errorf("expected parts[0] to be %s, got %s", "v1", parts[0])
			}
			if parts[1] != "tables" {
				t.Errorf("expected parts[1] to be %s, got %s", "tables", parts[1])
			}
			if parts[2] != "test" {
				t.Errorf("expected parts[2] to be %s, got %s", "test", parts[2])
			}
		})

		t.Run("should return an empty slice when the endpoint is not present", func(t *testing.T) {
			u, _ := url.Parse("http://test")
			req := http.Request{URL: u}
			parts := EndpointPartsFromHTTPRequest(req)
			if len(parts) != 0 {
				t.Errorf("expected parts to be empty, got %v", parts)
			}
		})
	})
}
