package polygon

import (
	"context"
	"net/http"

	"github.com/alpine-hodler/web/pkg/transport"
)

// Client is a wrapper for http.Client.
type Client struct{ http.Client }

// NewClient will return a new HTTP client to interface with the polygon.io web API.
func NewClient(_ context.Context, roundtripper transport.T) (*Client, error) {
	client := new(Client)
	client.Transport = roundtripper
	return client, nil
}
