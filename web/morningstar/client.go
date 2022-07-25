package morningstar

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alpine-hodler/web/pkg/transport"
)

// Client is a wrapper for http.Client.
type Client struct{ http.Client }

// NewBearerToken will generate a bearer token with base64-encrypted username and password.
func NewBearerToken(url, username, password string) (*Token, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	token := new(Token)
	if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
		return nil, err
	}
	return token, nil
}

// NewClient will return a new HTTP client to interface with the Morningstar API.
func NewClient(_ context.Context, roundtripper transport.T) (*Client, error) {
	client := new(Client)
	client.Transport = roundtripper
	return client, nil
}
