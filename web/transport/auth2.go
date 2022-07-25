package transport

import (
	"fmt"
	"net/http"
	"net/url"
)

// TODO
type Auth2 struct {
	bearer string
	url    *url.URL
}

// NewAuth2 will return an OAuth2 http transport.
func NewAuth2() *Auth2 {
	return new(Auth2)
}

// SetBearer will set the bearer field on Auth2.
func (auth *Auth2) SetBearer(val string) *Auth2 {
	auth.bearer = val
	return auth
}

// SetURL will set the key field on APIKey.
func (auth *Auth2) SetURL(u string) *Auth2 {
	auth.url, _ = url.Parse(u)
	return auth
}

// RoundTrip authorizes the request with a signed OAuth1 Authorization header using the author and TokenSource.
func (auth *Auth2) RoundTrip(req *http.Request) (*http.Response, error) {
	if auth.url == nil {
		return nil, fmt.Errorf("url is a required value on the HTTP client's transport")
	}

	req.URL.Scheme = auth.url.Scheme
	req.URL.Host = auth.url.Host

	req.Header.Set(authorizationHeaderParam, fmt.Sprintf("%s %s", bearerHeaderPrefix, auth.bearer))
	return http.DefaultTransport.RoundTrip(req)
}
