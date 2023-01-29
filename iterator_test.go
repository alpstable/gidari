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

	//t.Run("NewIterator", func(t *testing.T) {
	//	t.Parallel()

	//	for _, tcase := range []struct {
	//		name string
	//		cfg  *Config
	//		err  error
	//	}{
	//		{
	//			name: "nil",
	//			err:  ErrNilConfig,
	//		},
	//	} {
	//		t.Run(tcase.name, func(t *testing.T) {
	//			t.Parallel()

	//			_, err := NewIterator(context.Background(), tcase.cfg)
	//			if tcase.err != nil && !errors.Is(err, tcase.err) {
	//				t.Errorf("expected error %v, got %v", tcase.err, err)
	//			}

	//			if tcase.err == nil && err != nil {
	//				t.Errorf("expected no error, got %v", err)
	//			}
	//		})
	//	}
	//})

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
}
