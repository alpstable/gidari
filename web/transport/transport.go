package transport

import "net/http"

type T interface {
	http.RoundTripper
}
