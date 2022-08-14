package transport

import (
	"fmt"
	"net/http"
	"net/url"
)

type Basic struct {
	email, password string
	url             *url.URL
}

// NewBasic will return an Basic http transport.
func NewBasic() *Basic {
	return new(Basic)
}

// SetEmail will set the email field on Basic.
func (auth *Basic) SetEmail(val string) *Basic {
	auth.email = val
	return auth
}

// SetPassword will set the password field on Basic.
func (auth *Basic) SetPassword(val string) *Basic {
	auth.password = val
	return auth
}

// SetURL will set the key field on Basic.
func (auth *Basic) SetURL(val string) *Basic {
	auth.url, _ = url.Parse(val)
	return auth
}

// RoundTrip authorizes the request with a signed OAuth1 Authorization header using the author and TokenSource.
func (auth *Basic) RoundTrip(req *http.Request) (*http.Response, error) {
	if auth.url == nil {
		return nil, fmt.Errorf("url is a required value on the HTTP client's transport")
	}

	req.URL.Scheme = auth.url.Scheme
	req.URL.Host = auth.url.Host
	req.SetBasicAuth(auth.email, auth.password)
	return http.DefaultTransport.RoundTrip(req)
}
