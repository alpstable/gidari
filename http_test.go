// Copyright 2023 The Gidari Authors.
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
	"testing"
	"time"

	"golang.org/x/time/rate"
)

var errMissingURL = errors.New("missing URL")

func TestIterator(t *testing.T) {
	t.Parallel()

	t.Run("Iterator.Next", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name string

			// requestCount are the number of requests that the
			// iterator will make.
			requestCount int

			// setRateLimiters will set the rate limiter on the
			// requests, if true.
			setRateLimiters bool

			// forceErrOn will force the iterator to return an
			// error on the nth request. To make the iterator
			// return on the first request, set forceErrOn to -1.
			forceErrOn int

			// forceErr will be the error that the iterator will
			// return on the nth request. If this value is not
			// set, then the iterator will not return an error.
			forceErr error

			// wantErr will be the error that the iterator is
			// expected to return.
			wantErr error
		}{
			{
				name:            "many healthy requests",
				requestCount:    5,
				setRateLimiters: true,
			},
			{
				name:            "error on first request",
				requestCount:    5,
				setRateLimiters: true,
				forceErrOn:      -1,
				forceErr:        errMissingURL,
				wantErr:         errMissingURL,
			},
			{
				name:            "error on second request",
				requestCount:    5,
				setRateLimiters: true,
				forceErrOn:      1,
				forceErr:        errMissingURL,
				wantErr:         errMissingURL,
			},
			{
				name:            "error on last request",
				requestCount:    5,
				setRateLimiters: true,
				forceErrOn:      4,
				forceErr:        errMissingURL,
				wantErr:         errMissingURL,
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				// If there is no requestCount, then skip the
				// test.
				if tcase.requestCount == 0 {
					t.Skip("no request count")
				}

				var itr *HTTPIteratorService

				// urlSet is a set of URLs that the iterator
				// will request.
				urlSet := make(map[string]struct{})
				{
					// Create an HTTP Service to make the
					// requests.
					svc, err := NewService(context.Background())
					if err != nil {
						t.Fatalf("failed to create service: %v", err)
					}

					var rlimiter *rate.Limiter
					if tcase.setRateLimiters {
						rlimiter = rate.NewLimiter(rate.Every(time.Second), 100)
					}

					reqs := newHTTPRequests(tcase.requestCount)

					// If the forceErrOn is set, then set
					// the error on the nth request.

					var errReq *Request
					if tcase.forceErrOn >= 0 {
						errReq = reqs[tcase.forceErrOn]
					}

					// If forceErrOn is -1 then set the
					// error on the first request.
					if tcase.forceErrOn == -1 {
						errReq = reqs[0]
					}

					httpSvc := svc.HTTP.RateLimiter(rlimiter).Requests(reqs...)
					httpSvc.client = newMockHTTPClient(
						withMockHTTPClientResponseError(errReq, tcase.forceErr),
						withMockHTTPClientRequests(reqs...))

					// Set the urlSet using the requests.
					for _, req := range reqs {
						urlSet[req.http.URL.String()] = struct{}{}
					}

					itr = httpSvc.Iterator
				}

				{
					// Iterate over the request responses.
					for itr.Next(context.Background()) {
						rsp := itr.Current.Response

						url := rsp.Request.URL.String()
						if _, ok := urlSet[url]; !ok {
							t.Fatalf("unexpected url %s", url)
						}

						delete(urlSet, url)
					}

					err := itr.Err()

					if err == nil && tcase.wantErr != nil {
						t.Fatalf("expected error %v, got nil", tcase.wantErr)
					}

					if !errors.Is(err, tcase.wantErr) {
						t.Fatalf("expected error %v, got %v", tcase.wantErr, err)
					}

					// Ensure that all requests were made
					// unless there was an error.
					if len(urlSet) > 0 && err == nil {
						t.Errorf("expected all urls to be returned, got %d", len(urlSet))
					}
				}
			})
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name                              string
			expectedNumberOfUpsertsPerStorage int
			svc                               *Service

			// count is the number of times to call Do.
			count int

			// err is the expected error from Do.
			err error
		}{
			{
				name: "single request with single storage",
				svc: newMockService(mockServiceOptions{
					reqCount:       1,
					upsertStgCount: 1,
				}),
				expectedNumberOfUpsertsPerStorage: 1,
			},
			{
				name: "single request with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount:       1,
					upsertStgCount: 3,
				}),
				expectedNumberOfUpsertsPerStorage: 1,
			},
			{
				name: "multiple requests with single storage",
				svc: newMockService(mockServiceOptions{
					reqCount:       3,
					upsertStgCount: 1,
				}),
				expectedNumberOfUpsertsPerStorage: 3,
			},
			{
				name: "multiple requests with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount:       3,
					upsertStgCount: 3,
				}),
				expectedNumberOfUpsertsPerStorage: 3,
			},
			{
				name: "voluminous requests with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount:       10_000,
					upsertStgCount: 3,
					rateLimiter:    rate.NewLimiter(rate.Limit(1*time.Second), 10_000),
				}),
				expectedNumberOfUpsertsPerStorage: 10_000,
			},
			{
				name: "call Do many times",
				svc: newMockService(mockServiceOptions{
					reqCount:       1,
					upsertStgCount: 1,
				}),
				count: 10,
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				count := tcase.count
				if count == 0 {
					count = 1
				}

				for i := 0; i < count; i++ {
					err := tcase.svc.HTTP.Store(context.Background())
					if tcase.err != nil && !errors.Is(err, tcase.err) {
						t.Errorf("expected error %v, got %v", tcase.err, err)
					}

					if tcase.err == nil && err != nil {
						t.Errorf("expected no error, got %v", err)
					}
				}

				// If ccount > 1, then we stop here.
				if count > 1 {
					return
				}

				// We need to validate various operation for
				// the upsert storage.
				for _, req := range tcase.svc.HTTP.requests {
					for _, w := range req.writers {
						stg, ok := w.(*mockUpsertWriter)
						if !ok {
							t.Errorf("expected mock storage, got %T", w)
						}

						// The number of upserts should be equal to the
						// expected number of upserts. Note that there
						// can be less requests than upserts, for
						// example a timeseries request could be broken
						// into multiple flattened requests for upsert.
						if stg.count != tcase.expectedNumberOfUpsertsPerStorage {
							t.Errorf("expected %d upserts, got %d",
								tcase.expectedNumberOfUpsertsPerStorage,
								stg.count)
						}
					}
				}
			})
		}
	})
}

func BenchmarkIterator(b *testing.B) {
	// Create a new service.
	svc := newMockService(mockServiceOptions{
		reqCount:       10_000,
		upsertStgCount: 3,
		rateLimiter:    rate.NewLimiter(rate.Limit(1*time.Second), 10_000),
	})

	// Create a new iterator.
	itr := svc.HTTP.Iterator

	// Reset the benchmark timer.
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Iterate over the iterator.
		for itr.Next(context.Background()) {
			// Do nothing.
		}
	}

	// Check for errors.
	if err := itr.Err(); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkHTTPServiceDo(b *testing.B) {
	// Create a new service.
	svc := newMockService(mockServiceOptions{
		reqCount:       1,
		upsertStgCount: 3,
	})

	// Reset the benchmark timer.
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Execute the service.
		if err := svc.HTTP.Store(context.Background()); err != nil {
			b.Fatal(err)
		}
	}
}
