// The MIT License (MIT)

// Copyright (c) 2015 Dalton Hubble

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package transport

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Auth1 is an http.RoundTripper used to authenticate using the OAuth 1.a algorithm defined by twitter:
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/creating-a-signature
type Auth1 struct {
	accessToken       string
	accessTokenSecret string
	consumerKey       string
	consumerSecret    string
	url               *url.URL
}

// NewAuth1 will return an OAuth1 http transpoauth.
func NewAuth1() *Auth1 {
	return new(Auth1)
}

// RoundTrip authorizes the request with a signed OAuth1 Authorization header.
func (auth *Auth1) RoundTrip(req *http.Request) (*http.Response, error) {
	if auth.url == nil {
		return nil, fmt.Errorf("url is a required value on the HTTP client's transport")
	}
	req.URL.Scheme = auth.url.Scheme
	req.URL.Host = auth.url.Host

	err := auth.setRequestAuthHeader(req)
	if err != nil {
		return nil, err
	}
	return http.DefaultTransport.RoundTrip(req)
}

// SetConsumerKey will set the consumerKey field on Auth1.
func (auth *Auth1) SetAccessToken(val string) *Auth1 {
	auth.accessToken = val
	return auth
}

// SetAccessTokenSecret will set the accessTokenSecret field on Auth1.
func (auth *Auth1) SetAccessTokenSecret(val string) *Auth1 {
	auth.accessTokenSecret = val
	return auth
}

// SetConsumerKey will set the consumerKey field on Auth1.
func (auth *Auth1) SetConsumerKey(val string) *Auth1 {
	auth.consumerKey = val
	return auth
}

// SetConsumerSecret will set the consumerSecret field on Auth1.
func (auth *Auth1) SetConsumerSecret(val string) *Auth1 {
	auth.consumerSecret = val
	return auth
}

// setRequestAuthHeader sets the OAuth1 header for making authenticated requests with an AccessToken (token credential)
// according to RFC 5849 3.1.
func (auth *Auth1) setRequestAuthHeader(req *http.Request) error {
	oauthParams := map[string]string{
		oauthConsumerKeyParam:     auth.consumerKey,
		oauthSignatureMethodParam: defaultOauthSignatureMethod,
		oauthTimestampParam:       strconv.FormatInt(time.Now().Unix(), 10),
		oauthNonceParam:           nonce(),
		oauthVersionParam:         oauthVersion1,
	}
	oauthParams[oauthTokenParam] = auth.accessToken
	params, err := collectParameters(req, oauthParams)
	if err != nil {
		return err
	}
	signatureBase := signatureBase(req, params)
	signature, err := hmacSign(auth.consumerSecret, auth.accessTokenSecret, signatureBase, sha1.New)
	if err != nil {
		return err
	}
	oauthParams[oauthSignatureParam] = signature
	req.Header.Set(authorizationHeaderParam, authHeaderValue(oauthParams))
	return nil
}

// SetURL will set the key field on APIKey.
func (auth *Auth1) SetURL(u string) *Auth1 {
	auth.url, _ = url.Parse(u)
	return auth
}

// baseURI returns the base string URI of a request according to RFC 5849 3.4.1.2. The scheme and host are lowercased,
// the port is dropped if it is 80 or 443, and the path minus query parameters is included.
func baseURI(req *http.Request) string {
	scheme := strings.ToLower(req.URL.Scheme)
	host := strings.ToLower(req.URL.Host)
	if hostPort := strings.Split(host, ":"); len(hostPort) == 2 && (hostPort[1] == "80" || hostPort[1] == "443") {
		host = hostPort[0]
	}
	// TODO: use req.URL.EscapedPath() once Go 1.5 is more generally adopted
	// For now, hacky workaround accomplishes the same internal escaping mode escape(u.Path, encodePath) for proper
	// compliance with the OAuth1 spec.
	path := req.URL.Path
	if path != "" {
		path = strings.Split(req.URL.RequestURI(), "?")[0]
	}
	return fmt.Sprintf("%v://%v%v", scheme, host, path)
}

// collectParameters collects request parameters from the request query, OAuth parameters (which should exclude
// oauth_signature), and the request body provided the body is single part, form encoded, and the form content type
// header is set. The returned map of collected parameter keys and values follow RFC 5849 3.4.1.3, except duplicate
// parameters are not supported.
func collectParameters(req *http.Request, oauthParams map[string]string) (map[string]string, error) {
	// add oauth, query, and body parameters into params
	params := map[string]string{}
	for key, value := range req.URL.Query() {
		// most backends do not accept duplicate query keys
		params[key] = value[0]
	}
	if req.Body != nil && req.Header.Get(contentType) == formContentType {
		// reads data to a []byte, draining req.Body
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		values, err := url.ParseQuery(string(b))
		if err != nil {
			return nil, err
		}
		for key, value := range values {
			// not supporting params with duplicate keys
			params[key] = value[0]
		}
		// reinitialize Body with ReadCloser over the []byte
		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}
	for key, value := range oauthParams {
		params[key] = value
	}
	return params, nil
}

func hmacSign(consumerSecret, tokenSecret, message string, algo func() hash.Hash) (string, error) {
	signingKey := strings.Join([]string{consumerSecret, tokenSecret}, "&")
	mac := hmac.New(algo, []byte(signingKey))
	mac.Write([]byte(message))
	signatureBytes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signatureBytes), nil
}

// nonce provides a random nonce string.
func nonce() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

// shouldEscape returns false if the byte is an unreserved character that should not be escaped and true otherwise,
// according to RFC 3986 2.1.
func shouldEscape(c byte) bool {
	// RFC3986 2.3 unreserved characters
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '.', '_', '~':
		return false
	}
	// all other bytes must be escaped
	return true
}

// percentEncode percent encodes a string according to RFC 3986 2.1.
func percentEncode(input string) string {
	var buf bytes.Buffer
	for _, b := range []byte(input) {
		// if in unreserved set
		if shouldEscape(b) {
			buf.Write([]byte(fmt.Sprintf("%%%02X", b)))
		} else {
			// do not escape, write byte as-is
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// encodeParameters percent encodes parameter keys and values according to RFC5849 3.6 and RFC3986 2.1 and returns a new
// map.
func encodeParameters(params map[string]string) map[string]string {
	encoded := map[string]string{}
	for key, value := range params {
		encoded[percentEncode(key)] = percentEncode(value)
	}
	return encoded
}

// sortParameters sorts parameters by key and returns a slice of key/value pairs formatted with the given format string
// (e.g. "%s=%s").
func sortParameters(params map[string]string, format string) []string {
	// sort by key
	keys := make([]string, len(params))
	i := 0
	for key := range params {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// parameter join
	pairs := make([]string, len(params))
	for i, key := range keys {
		pairs[i] = fmt.Sprintf(format, key, params[key])
	}
	return pairs
}

// authHeaderValue formats OAuth parameters according to RFC 5849 3.5.1. OAuth params are percent encoded, sorted by key
// (for testability), and joined by "=" into pairs. Pairs are joined with a ", " comma separator into a header string.
// The given OAuth params should include the "oauth_signature" key.
func authHeaderValue(oauthParams map[string]string) string {
	pairs := sortParameters(encodeParameters(oauthParams), `%s="%s"`)
	return fmt.Sprintf("%s %s", authorizationPrefix, strings.Join(pairs, ", "))
}

// parameterString normalizes collected OAuth parameters (which should exclude oauth_signature) into a parameter string
// as defined in RFC 5894 3.4.1.3.2. The parameters are encoded, sorted by key, keys and values joined with "&", and
// pairs joined with "=" (e.g. foo=bar&q=gopher).
func normalizedParameterString(params map[string]string) string {
	return strings.Join(sortParameters(encodeParameters(params), "%s=%s"), "&")
}

// signatureBase combines the uppercase request method, percent encoded base string URI, and normalizes the request
// parameters int a parameter string. Returns the OAuth1 signature base string according to RFC5849 3.4.1.
func signatureBase(req *http.Request, params map[string]string) string {
	method := strings.ToUpper(req.Method)
	baseURL := baseURI(req)
	parameterString := normalizedParameterString(params)
	// signature base string constructed accoding to 3.4.1.1
	baseParts := []string{method, percentEncode(baseURL), percentEncode(parameterString)}
	return strings.Join(baseParts, "&")
}
