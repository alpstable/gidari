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
	"testing"
	"time"

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

	t.Run("Iterator.Next", func(t *testing.T) {
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

				var itr *HTTPIteratorService

				// urlSet is a set of URLs that the iterator will request.
				urlSet := make(map[string]struct{})

				{
					// Create an HTTP Service to make the requests.
					svc, err := NewService(context.Background())
					if err != nil {
						t.Fatalf("failed to create service: %v", err)
					}

					var rlimiter *rate.Limiter
					if tcase.setRateLimiters {
						rlimiter = rate.NewLimiter(rate.Every(time.Second), 100)
					}

					httpSvc := svc.HTTP.RateLimiter(rlimiter).Client(&mockClient{})

					// Add the requests to the config.
					for i := 0; i < tcase.requestCount; i++ {
						url := fmt.Sprintf("https://example.com/%d", i)
						urlSet[url] = struct{}{}

						httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
						if err != nil {
							t.Fatalf("failed to create request: %v", err)
						}

						httpSvc.Requests(&HTTPRequest{Request: httpRequest})
					}

					itr = httpSvc.Iterator
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

	t.Run("Do", func(t *testing.T) {
		t.Parallel()

		for _, tcase := range []struct {
			name                              string
			expectedNumberOfUpsertsPerStorage int
			svc                               *Service
			err                               error
		}{
			{
				name: "single request with single storage",
				svc: newMockService(mockServiceOptions{
					reqCount: 1,
					stgCount: 1,
				}),
				expectedNumberOfUpsertsPerStorage: 1,
			},
			{
				name: "single request with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount: 1,
					stgCount: 3,
				}),
				expectedNumberOfUpsertsPerStorage: 1,
			},
			{
				name: "multiple requests with single storage",
				svc: newMockService(mockServiceOptions{
					reqCount: 3,
					stgCount: 1,
				}),
				expectedNumberOfUpsertsPerStorage: 3,
			},
			{
				name: "multiple requests with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount: 3,
					stgCount: 3,
				}),
				expectedNumberOfUpsertsPerStorage: 3,
			},
			{
				name: "voluminous requests with multiple storages",
				svc: newMockService(mockServiceOptions{
					reqCount:    10_000,
					stgCount:    3,
					rateLimiter: rate.NewLimiter(rate.Limit(1*time.Second), 10_000),
				}),
				expectedNumberOfUpsertsPerStorage: 10_000,
			},
		} {
			tcase := tcase

			t.Run(tcase.name, func(t *testing.T) {
				t.Parallel()

				_, err := tcase.svc.HTTP.Do(context.Background())
				if tcase.err != nil && !errors.Is(err, tcase.err) {
					t.Errorf("expected error %v, got %v", tcase.err, err)
				}

				if tcase.err == nil && err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				// If there is no mock storage then we can terminate
				// the test here.
				if len(tcase.svc.storage) == 0 {
					return
				}

				// We need to validate various operation for each
				// storage object.
				for _, stg := range tcase.svc.storage {
					mockStorage, ok := stg.Storage.(*mockStorage)
					if !ok {
						t.Errorf("expected mock storage, got %T", stg)
					}

					// The number of upserts should be equal to the
					// expected number of upserts. Note that there
					// can be less requests than upserts, for
					// example a timeseries request could be broken
					// into multiple flattened requests for upsert.
					if mockStorage.upsertCount != tcase.expectedNumberOfUpsertsPerStorage {
						t.Errorf("expected %d upserts, got %d", tcase.expectedNumberOfUpsertsPerStorage,
							mockStorage.upsertCount)
					}
				}
			})
		}
	})
}
