package auth

import "net/http"

type Transport interface {
	http.RoundTripper
}
