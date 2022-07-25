package coinbasepro

import (
	"context"
	"net/http"

	"github.com/alpine-hodler/web/pkg/transport"
)

const (
	coinbaseTimeLayout1 = "2006-01-02 15:04:05.999999999+00" // Some dumbass coinbase thing
	websocketURL        = "wss://ws-feed.pro.coinbase.com"
)

// Client is a wrapper for http.Client.
type Client struct{ http.Client }

// NewClient will return a new HTTP client to interface with the Coinbase Pro API.
func NewClient(_ context.Context, roundtripper transport.T) (*Client, error) {
	client := new(Client)
	client.Transport = roundtripper
	return client, nil
}
