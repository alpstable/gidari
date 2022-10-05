package tools

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

// HTTPMessage represents data exchanged between server and client, typed for encryption purposes.
type HTTPMessage string

// NewHTTPMessage generates the base64-encoded message required to make API-Key-Authenticated requests.
func NewHTTPMessage(req *http.Request, timestamp string) HTTPMessage {
	if req.Body == nil {
		return HTTPMessage([]byte{})
	}

	body, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	requestPath := req.URL.Path
	if req.URL.RawQuery != "" {
		requestPath = fmt.Sprintf("%s?%s", req.URL.Path, req.URL.RawQuery)
	}

	return HTTPMessage(fmt.Sprintf("%s%s%s%s", timestamp, req.Method, requestPath, string(body)))
}

// Sign generates the base64-encoded signature required to make requests. In particular, the signed header is generated
// by creating a sha256 HMAC using the base64-decoded secret key on the prehash string
// timestamp + method + requestPath + body (where + represents string concatenation) and base64-encode the output. The
// timestamp value is the same as the timestamp header.
func (msg HTTPMessage) Sign(secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("error decoding secret: %w", err)
	}

	signature := hmac.New(sha256.New, key)

	_, err = signature.Write([]byte(msg))
	if err != nil {
		return "", fmt.Errorf("error writing signature: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
