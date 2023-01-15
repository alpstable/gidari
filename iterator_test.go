// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/alpstable/gidari/internal/web"
	"golang.org/x/time/rate"
)

type mockClient struct{}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	rsp := &http.Response{
		Request: req,
	}

	return rsp, nil
}

func TestIterator(t *testing.T) {
	t.Parallel()

	t.Run("NewIterator", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name string
			cfg  *Config
			err  error
		}{
			{
				name: "nil",
				err:  ErrNilConfig,
			},
		} {
			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				_, err := NewIterator(context.Background(), tcase.cfg)
				if tcase.err != nil && !errors.Is(err, tcase.err) {
					t.Errorf("expected error %v, got %v", tcase.err, err)
				}

				if tcase.err == nil && err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			})
		}
	})

	t.Run("Next", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name string

			// requestCount are the number of requests that the iterator will make.
			requestCount int

			// setRateLimiters will set the rate limiter on the requests, if true.
			setRateLimiters bool

			// err is the error that the iterator will return.
			err error
		}{
			{
				name:            "many healthy requests",
				requestCount:    5,
				setRateLimiters: true,
			},
		} {
			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				var cfg *Config
				var itr *Iterator

				// urlSet is a set of URLs that the iterator will request.
				urlSet := make(map[string]struct{})

				{
					// Setup the configuration and the iterator.
					cfg = &Config{
						Client: &mockClient{},
					}

					if tcase.setRateLimiters {
						cfg.RateLimiter = rate.NewLimiter(rate.Every(time.Second), 100)
					}

					// Add the requests to the config.
					for i := 0; i < tcase.requestCount; i++ {
						url := fmt.Sprintf("https://example.com/%d", i)
						urlSet[url] = struct{}{}

						httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
						if err != nil {
							t.Fatalf("failed to create request: %v", err)
						}

						cfg.Requests = append(cfg.Requests, &Request{
							Request: httpRequest,
						})
					}

					var err error
					itr, err = NewIterator(context.Background(), cfg)

					if err != nil {
						t.Fatalf("unexpected error %v", err)
					}
				}

				{
					// Iterate over the request responses.
					for itr.Next(context.Background()) {
						url := itr.Current.Response.Request.URL.String()
						if _, ok := urlSet[url]; !ok {
							t.Errorf("unexpected url %s", url)
						}

						delete(urlSet, url)
					}

					// Ensure that no unexpected errors occur.
					if err := itr.Err(); err != nil {
						t.Errorf("unexpected error %v", err)
					}

					// Ensure that all requests were made.
					if len(urlSet) > 0 {
						t.Errorf("expected all urls to be returned, got %d", len(urlSet))
					}
				}
			})
		}
	})
}

type webWorkerMockClient struct {
	mutex sync.Mutex

	count    int
	idErrors map[int]error
}

func (w *webWorkerMockClient) Do(req *http.Request) (*http.Response, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.idErrors == nil {
		return nil, nil
	}

	w.count++

	if err, ok := w.idErrors[w.count]; ok {
		return nil, err
	}

	return nil, nil
}

func TestStartWebWorker(t *testing.T) {
	t.Parallel()

	errTest := fmt.Errorf("test error")

	for _, tcase := range []struct {
		name string
		ctx  context.Context
		cfg  *webWorkerConfig

		// poolCount is the number of workers to start.
		poolCount int

		// jobcount is the number of jobs to send to the workers for req/rsp processing.
		jobCount int

		// idToError is a map of job IDs to throw the mapping error on. If this is nil, then no planned errors
		// will be thrown.
		idToError map[int]error

		// expectedErr is the error expected to be returned by the worker.
		expectedError error
	}{
		{
			name:          "empty",
			cfg:           &webWorkerConfig{},
			expectedError: context.DeadlineExceeded,
		},
		{
			name:          "empty with high pool count",
			cfg:           &webWorkerConfig{},
			poolCount:     100,
			expectedError: context.DeadlineExceeded,
		},
		{
			name:     "one job on a buffered channel",
			jobCount: 1,
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 1),
			},
		},
		{
			name:      "one job on a buffered channel with high pool count",
			jobCount:  1,
			poolCount: 100,
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 1),
			},
		},
		{
			name:     "many jobs on a buffered channel",
			jobCount: 10,
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 10),
			},
		},
		{
			name:     "many jobs with error on first",
			jobCount: 10,
			idToError: map[int]error{
				1: errTest,
			},
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 10),
			},
			expectedError: errTest,
		},
		{
			name:     "many jobs with error in middle",
			jobCount: 10,
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 10),
			},
			idToError: map[int]error{
				3: errTest,
			},
			expectedError: errTest,
		},
		{
			name:     "many jobs with error on last",
			jobCount: 10,
			cfg: &webWorkerConfig{
				currentCh: make(chan *Current, 1),
				jobs:      make(chan webWorkerJob, 10),
			},
			idToError: map[int]error{
				10: errTest,
			},
			expectedError: errTest,
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			// If the test context is available use it, otherwise create a context with a timeout. If the
			// timeout occurs, then the test will fail.
			ctx := tcase.ctx
			if ctx == nil {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

				defer cancel()
			}

			// If the error channel is not set, then create a channel with a buffer size of 1.
			if tcase.cfg.errCh == nil {
				tcase.cfg.errCh = make(chan error, 1)
			}

			// If the pool count is not set, then set it to 1.
			if tcase.poolCount == 0 {
				tcase.poolCount = 1
			}

			// If the done channel is not set, then create a channel with a buffer size of the number of
			// jobs.
			if tcase.cfg.done == nil {
				tcase.cfg.done = make(chan bool, tcase.jobCount)
			}

			for w := 0; w < tcase.poolCount; w++ {
				go startWebWorker(ctx, tcase.cfg)
			}

			// Send empty requests to the jobs channel.
			for i := 0; i < tcase.jobCount; i++ {
				job := &webWorkerJob{
					req: &flattenedRequest{
						fetchConfig: &web.FetchConfig{
							Request: &http.Request{},
						},
						client: &webWorkerMockClient{
							idErrors: tcase.idToError,
						},
					},
				}

				// Add the default request handler to the configuration.
				//job.req.reqHandler = func(_ context.Context,
				//	_ *web.FetchConfig) (*http.Response, error) {

				//	return &http.Response{}, nil
				//}

				// Add the artifical errors to the configuration.
				//if tcase.idToError != nil {
				//	//handlerCount := 0
				//	//handlerCountMtx := &sync.Mutex{}
				//	//job.req.client.Do = func(_ *http.Response) (*http.Response, error) {
				//	//	handlerCountMtx.Lock()
				//	//	defer handlerCountMtx.Unlock()

				//	//	handlerCount++

				//	//	if err, ok := tcase.idToError[handlerCount]; ok {
				//	//		return err
				//	//	}

				//	//	return nil
				//	//}
				//}

				tcase.cfg.jobs <- *job
			}

			if tcase.cfg.jobs != nil {
				close(tcase.cfg.jobs)
			}

			for i := 0; i < tcase.jobCount; i++ {
				<-tcase.cfg.done
			}

			if err := <-tcase.cfg.errCh; err != nil {
				expErr := tcase.expectedError
				if expErr != nil {
					if !errors.Is(err, expErr) {
						t.Fatalf("expected error %v, got %v", expErr, err)
					}
				} else {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

//func newIteratorConfig(t *testing.T, uri string) *Config {
//	t.Helper()
//
//	return &Config{
//		URL: func() *url.URL {
//			u, _ := url.Parse("http://localhost")
//			return u
//		}(),
//	}
//}
//
//func newIteratorConfigWithHandler(t *testing.T, uri string, h WebResultAssigner) *Config {
//	t.Helper()
//
//	cfg := newIteratorConfig(t, uri)
//	cfg.HTTPResponseHandler = h
//
//	return cfg
//}
//
//func newIteratorHandler(t *testing.T) WebResultAssigner {
//	t.Helper()
//
//	return func(ctx context.Context, httpResponse WebResult) ([]*proto.IteratorResult, error) {
//		return []*proto.IteratorResult{
//			{
//				URL:  "http://localhost",
//				Data: []byte("test"),
//			},
//		}, nil
//	}
//}

//func TestNewIterator(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name string
//		cfg  *Config
//		want *Iterator
//		err  error
//	}{
//		{
//			name: "nil config",
//			cfg:  nil,
//			err:  ErrNilConfig,
//		},
//		{
//			name: "empty config",
//			cfg:  &Config{},
//			err:  ErrMissingURL,
//		},
//		{
//			name: "default HTTP Handler",
//			cfg:  newIteratorConfig(t, "http://localhost"),
//			want: &Iterator{
//				cfg:               newIteratorConfig(t, "http://localhost"),
//				webResultAssigner: assignWebResult,
//			},
//		},
//		{
//			name: "custom HTTP Handler",
//			cfg:  newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),
//			want: &Iterator{
//				cfg: newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),
//
//				webResultAssigner: newIteratorHandler(t),
//			},
//		},
//	} {
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, err := NewIterator(context.Background(), tcase.cfg)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if tcase.err != nil {
//				return
//			}
//
//			httpResponse := WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("test"))),
//				},
//			}
//
//			gotResults, err := got.webResultAssigner(context.Background(), httpResponse)
//			if err != nil {
//				t.Fatalf("got error %v, want nil", err)
//			}
//
//			wantResults, err := tcase.want.webResultAssigner(context.Background(), httpResponse)
//			if err != nil {
//				t.Fatalf("got error %v, want nil", err)
//			}
//
//			if len(gotResults) != len(wantResults) {
//				t.Fatalf("got %d results, want %d", len(gotResults), len(wantResults))
//			}
//
//			if !reflect.DeepEqual(gotResults, wantResults) {
//				t.Fatalf("got %v, want %v", gotResults[0], wantResults[0])
//			}
//		})
//	}
//}
//
//func TestSanitizeInvalidJSON(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name       string
//		data       []byte
//		clobColumn string
//		want       []byte
//		ok         bool
//		err        error
//	}{
//		{
//			name: "empty",
//			data: []byte{},
//			want: []byte{},
//		},
//		{
//			name: "invalid, no clob column",
//			data: []byte("invalid"),
//			ok:   false,
//		},
//		{
//			name:       "invalid, clob column",
//			data:       []byte("invalid"),
//			clobColumn: "test",
//			want:       []byte(`{"test":"invalid"}`),
//			ok:         true,
//		},
//		{
//			name: "valid, no clob column",
//			data: []byte(`{"test": "test"}`),
//			want: []byte(`{"test": "test"}`),
//			ok:   true,
//		},
//		{
//			name:       "valid, clob column",
//			data:       []byte(`{"test": "test"}`),
//			clobColumn: "doesnt matter",
//			want:       []byte(`{"test": "test"}`),
//			ok:         true,
//		},
//	} {
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, ok, err := sanitizeJSON(tcase.data, tcase.clobColumn)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if ok != tcase.ok {
//				t.Fatalf("got ok %v, want %v", ok, tcase.ok)
//			}
//
//			if !bytes.Equal(got, tcase.want) {
//				t.Fatalf("got %s, want %s", got, tcase.want)
//			}
//		})
//	}
//}
//
//func TestAssignWebResult(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name string
//		rsp  WebResult
//		err  error
//		want []*proto.IteratorResult
//	}{
//		{
//			name: "nil body",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: nil,
//				},
//			},
//			want: []*proto.IteratorResult{},
//		},
//		{
//			name: "empty body",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
//				},
//			},
//			want: []*proto.IteratorResult{},
//		},
//		{
//			name: "invalid JSON",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("invalid"))),
//				},
//			},
//			want: nil,
//		},
//		{
//			name: "valid JSON",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"test": "test"}`))),
//				},
//
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "invalid JSON, clob column",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("invalid"))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//				ClobColumn: "test",
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"invalid"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "valid JSON, clob column",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"test": "test"}`))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//				ClobColumn: "doesn't matter",
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "valid JSON array",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`[{"test": "test"},
//{"test2": "test"}]`))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//				{
//					Data: []byte(`{"test2":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//	} {
//		tcase := tcase
//
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, err := assignWebResult(context.Background(), tcase.rsp)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if len(got) != len(tcase.want) {
//				t.Fatalf("got %d results, want %d", len(got), len(tcase.want))
//			}
//
//			for idx, result := range got {
//				if result.GetURL() != tcase.want[idx].GetURL() {
//					t.Fatalf("got %s, want %s", result.GetURL(), tcase.want[idx].GetURL())
//				}
//
//				if !bytes.Equal(result.GetData(), tcase.want[idx].GetData()) {
//					t.Fatalf("got %s, want %s", result.GetData(), tcase.want[idx].GetData())
//				}
//			}
//		})
//	}
//}

//func mockResponseBodyArray(b *testing.B, size int) []byte {
//	b.Helper()
//
//	var buf bytes.Buffer
//
//	buf.WriteString("[")
//
//	for i := 0; i < size; i++ {
//		buf.WriteString(`{"test": "test"}`)
//
//		if i < size-1 {
//			buf.WriteString(",")
//		}
//	}
//
//	buf.WriteString("]")
//
//	return buf.Bytes()
//}
//
//func mockResponseBodyObject(b *testing.B) []byte {
//	b.Helper()
//
//	return []byte(`{"test": "test"}`)
//}
//
//func BenchmarkAssignWebResult(b *testing.B) {
//	b.ReportAllocs()
//
//	for _, tcase := range []struct {
//		name string
//		rsp  WebResult
//	}{
//		{
//			name: "valid JSON array",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader(mockResponseBodyArray(b, 100))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//		},
//		{
//			name: "valid JSON object",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader(mockResponseBodyObject(b))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//		},
//	} {
//		tcase := tcase
//
//		b.Run(tcase.name, func(b *testing.B) {
//			b.ReportAllocs()
//
//			for i := 0; i < b.N; i++ {
//				assignWebResult(context.Background(), tcase.rsp)
//			}
//		})
//	}
//}
