package twitter

import (
	"context"
	"net/http"

	"github.com/alpine-hodler/web/pkg/transport"
)

const (
	defaultTwitterAPIURL = "https://api.twitter.com"
	identifier           = "Twitter"
	TwitterISO8601       = "2006-01-02T15:04:05.000Z"
)

// Client is a wrapper for http.Client.
type Client struct{ http.Client }

// NewClient will return an HTTP client to interface with the Twitter API.
func NewClient(_ context.Context, transport transport.T) (*Client, error) {
	client := new(Client)
	client.Transport = transport
	return client, nil
}
