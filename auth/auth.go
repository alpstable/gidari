// Package auth contains a non-exhaustive list of custom authentication round
// trippers to be used as authentication middleware with a gidari HTTP Service.
package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// RoundTrip is an HTTP round tripper that acts as a middleware to add
// auth requirements to HTTP requests.
type RoundTrip func(*http.Request) (*http.Response, error)

// NewCoinbaseRoundTrip will return a "RoundTrip" function that can be used as
// a "RoundTrip" function in an "http.RoundTripper" interface to authenticate
// requests to the Coinbase Cloud API.
func NewCoinbaseRoundTrip(key, secret, passphrase string) RoundTrip {
	return func(req *http.Request) (*http.Response, error) {
		var body []byte
		if req.Body != nil {
			body, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		requestPath := req.URL.Path
		if req.URL.RawQuery != "" {
			requestPath = fmt.Sprintf("%s?%s", req.URL.Path, req.URL.RawQuery)
		}

		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		msg := fmt.Sprintf("%s%s%s%s", timestamp, req.Method, requestPath, string(body))

		skey, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, fmt.Errorf("error decoding secret: %w", err)
		}

		signature := hmac.New(sha256.New, skey)

		// Don't handle error because hash.Write method never returns an error.
		signature.Write([]byte(msg))
		sig := base64.StdEncoding.EncodeToString(signature.Sum(nil))

		req.Header.Set("content-type", "application/json")
		req.Header.Add("cb-access-key", key)
		req.Header.Add("cb-access-passphrase", passphrase)
		req.Header.Add("cb-access-sign", sig)
		req.Header.Add("cb-access-timestamp", timestamp)

		rsp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}

		return rsp, nil
	}
}
