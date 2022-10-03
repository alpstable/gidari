package web

import (
	"net/http"
	"net/http/httptest"
)

// createTestServerWithBasicAuth is a helper that creates a httptest.Server with a handler that has basic authentication.
func createTestServerWithBasicAuth(username, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUsername, reqPassword, ok := r.BasicAuth()
		if !ok || reqUsername != username || reqPassword != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
}
