package auth

import (
	"fmt"
	"net/http"
)

var (
	ErrURLRequired   = fmt.Errorf("url is a required value on the HTTP client's transport")
	ErrRequestFailed = fmt.Errorf("request failed")
)

type Transport interface {
	http.RoundTripper
}
