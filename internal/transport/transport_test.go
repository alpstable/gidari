//go:build utests

// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package transport

import (
	"context"
	"net/url"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/alpstable/gidari/config"
	"github.com/alpstable/gidari/internal/web"
	"golang.org/x/time/rate"
)

func TestTimeseries(t *testing.T) {
	t.Parallel()
	t.Run("chunks where end date is before last iteration", func(t *testing.T) {
		t.Parallel()

		timeseries := &config.Timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T00:00:00Z")
		testURL.RawQuery = query.Encode()

		err = chunkTimeseries(timeseries, *testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 0, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.Chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.Chunks)
		}
	})

	t.Run("chunks where end date is equal to last iteration", func(t *testing.T) {
		t.Parallel()

		timeseries := &config.Timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T01:00:00Z")
		testURL.RawQuery = query.Encode()

		err = chunkTimeseries(timeseries, *testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.Chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.Chunks)
		}
	})

	t.Run("chunks where end date is after last iteration", func(t *testing.T) {
		t.Parallel()
		timeseries := &config.Timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		testURL, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := testURL.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T02:00:00Z")
		testURL.RawQuery = query.Encode()

		err = chunkTimeseries(timeseries, *testURL)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 0o5, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 0o5, 11, 1, 0, 0, 0, time.UTC),
				time.Date(2022, 0o5, 11, 2, 0, 0, 0, time.UTC),
			},
		}

		if !reflect.DeepEqual(expChunks, timeseries.Chunks) {
			t.Fatalf("unexpected chunks: %v", timeseries.Chunks)
		}
	})
}

func TestNewFetchConfig(t *testing.T) {
	t.Parallel()

	testURL, err := url.Parse("https://fuzz")
	if err != nil {
		t.Fatalf("error parsing url: %v", err)
	}

	testClient := &web.Client{}

	testRateLimiter := rate.NewLimiter(rate.Every(time.Second), 2)

	for _, tcase := range []struct {
		req      *config.Request
		expected *web.FetchConfig
	}{
		{
			req: &config.Request{
				Endpoint:    "/test",
				Method:      "POST",
				RateLimiter: testRateLimiter,
				Query: map[string]string{
					"hey": "ho",
					"foo": "bar",
				},
			},
			expected: &web.FetchConfig{
				Method:      "POST",
				URL:         testURL,
				C:           testClient,
				RateLimiter: testRateLimiter,
			},
		},
	} {
		fetchConfig := newFetchConfig(tcase.req, *testURL, testClient)

		if fetchConfig.Method != tcase.expected.Method {
			t.Error("unexpected fetch config Method")
		}

		if fetchConfig.C != tcase.expected.C {
			t.Error("unexpected fetch config Client")
		}

		if fetchConfig.RateLimiter != tcase.expected.RateLimiter {
			t.Error("unexpected fetch config RateLimiter")
		}

		// URL Path should be concatenated with request Endpoint
		if fetchConfig.URL.Path != path.Join(testURL.Path, tcase.req.Endpoint) {
			t.Errorf("unexpected fetch config URL Path: %v", fetchConfig.URL.Path)
		}

		// Query parameters should be added to the URL Query
		for k, v := range tcase.req.Query {
			if fetchConfig.URL.Query().Get(k) != v {
				t.Errorf("unexpected fetch config URL Query for %v", k)
			}
		}
	}
}

func Test_flattenConfigRequests(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	tests := []struct {
		name     string
		args     args
		wantReqs []*flattenedRequest
		wantErr  bool
	}{
		{
			name: "successfully flatten a config request",
			args: args{
				cfg: &config.Config{
					URL: &url.URL{
						Path: "path",
					},
					Requests: []*config.Request{
						{
							Method:   "POST",
							Endpoint: "endpoint",
						},
					},
				},
			},
			wantReqs: []*flattenedRequest{
				{
					fetchConfig: &web.FetchConfig{
						Method: "POST",
						URL: &url.URL{
							Path: "path/endpoint",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail to flatten config requests when none are passed",
			args: args{
				cfg: &config.Config{},
			},
			wantErr: true,
		},
		{
			name: "fail to flatten request time series",
			args: args{
				cfg: &config.Config{
					URL: &url.URL{},
					Requests: []*config.Request{
						{
							Timeseries: &config.Timeseries{
								StartName: "startName",
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReqs, err := flattenConfigRequests(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("flattenConfigRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotReqs) != len(tt.wantReqs) {
				t.Errorf("flattenConfigRequests() got %d requests, want %d", len(gotReqs), len(tt.wantReqs))
			}
			for _, gotReq := range gotReqs {
				for _, wantReq := range tt.wantReqs {
					if !compareRequests(gotReq, wantReq) {
						t.Errorf("flattenConfigRequests() request comparison failed")
					}
				}
			}
		})
	}
}

func compareRequests(request1, request2 *flattenedRequest) bool {
	if request1 == nil || request2 == nil {
		return false
	}
	if request1.fetchConfig.URL.Path != request2.fetchConfig.URL.Path {
		return false
	}
	if request1.fetchConfig.Method != request2.fetchConfig.Method {
		return false
	}
	return true
}
