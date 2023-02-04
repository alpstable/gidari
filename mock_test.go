package gidari

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/time/rate"
)

type mockServiceOptions struct {
	upsertStgCount int
	reqCount       int
	rateLimiter    *rate.Limiter
}

func newMockService(opts mockServiceOptions) *Service {
	svc, err := NewService(context.Background())
	if err != nil {
		panic(err)
	}

	reqs := newHTTPRequests(opts.reqCount)

	svc.HTTP.
		Client(newMockHTTPClient(
			withMockHTTPClientRequests(reqs...),
		)).
		RateLimiter(opts.rateLimiter).
		Requests(reqs...).
		UpsertWriters(newMockUpsertStorage(opts.upsertStgCount)...)

	return svc
}

func newHTTPRequests(volume int) []*HTTPRequest {
	requests := make([]*HTTPRequest, volume)
	for i := 0; i < volume; i++ {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://example%d", i), nil)
		requests[i] = &HTTPRequest{
			Request: req,
			Table:   fmt.Sprintf("table%d", i),
		}
	}

	return requests
}

type mockHTTPClientResponseError struct {
	rsp *http.Response
	err error
}

type mockHTTPClient struct {
	mutex     sync.Mutex
	responses map[*http.Request]*mockHTTPClientResponseError
}

type mockHTTPClientOption func(*mockHTTPClient)

func newMockHTTPClient(opts ...mockHTTPClientOption) *mockHTTPClient {
	client := &mockHTTPClient{
		responses: make(map[*http.Request]*mockHTTPClientResponseError),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// withMockHTTPClientRequests will set the mockHTTPClient responses to the
// provided requests.
func withMockHTTPClientRequests(reqs ...*HTTPRequest) mockHTTPClientOption {
	return func(client *mockHTTPClient) {
		for _, req := range reqs {
			body := io.NopCloser(bytes.NewBufferString(""))
			code := http.StatusOK

			rsp := &http.Response{
				Body:       body,
				StatusCode: code,
				Request:    req.Request,
			}

			rspErr := &mockHTTPClientResponseError{rsp: rsp}

			// If the request has already been set, then just
			// update the response.
			if _, ok := client.responses[req.Request]; ok {
				client.responses[req.Request].rsp = rspErr.rsp

				continue
			}

			client.responses[req.Request] = rspErr
		}
	}
}

func withMockHTTPClientResponseError(req *HTTPRequest, err error) mockHTTPClientOption {
	return func(client *mockHTTPClient) {
		if req == nil {
			return
		}

		if err == nil {
			return
		}

		// If the request has already been set, then just
		// update the error.
		if _, ok := client.responses[req.Request]; ok {
			client.responses[req.Request].err = err

			return
		}

		client.responses[req.Request] = &mockHTTPClientResponseError{
			err: err,
		}
	}
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.responses) == 0 {
		return nil, nil
	}

	rsp := m.responses[req]

	// If the response has an error, return it.
	if rsp.err != nil {
		return nil, rsp.err
	}

	return rsp.rsp, nil
}

type mockUpsertWriter struct {
	count   int
	countMu sync.Mutex
}

func newMockUpsertStorage(volume int) []proto.UpsertWriter {
	stg := make([]proto.UpsertWriter, volume)
	for i := 0; i < volume; i++ {
		stg[i] = &mockUpsertWriter{}
	}

	return stg
}

func (m *mockUpsertWriter) Write(context.Context, *proto.UpsertRequest) error {
	m.countMu.Lock()
	defer m.countMu.Unlock()

	m.count++

	return nil
}
