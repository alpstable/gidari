// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestTimeseries(t *testing.T) {
	t.Parallel()
	t.Run("chunks where end date is before last iteration", func(t *testing.T) {
		t.Parallel()

		timeseries := &Timeseries{
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

		timeseries := &Timeseries{
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
		timeseries := &Timeseries{
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

//func TestWebWorker(t *testing.T) {
//	t.Parallel()
//
//	t.Run("normal json response", func(t *testing.T) {
//		t.Parallel()
//
//		const (
//			baseurl  = "https://pokeapi.co"
//			endpoint = "/api/v2/pokemon/ditto"
//			table    = "test"
//		)
//
//		webWorkerJobs := make(chan *webJob, 1)
//		repoJobs := make(chan *repoJob, 1)
//
//		client, err := web.NewClient(context.Background(), nil)
//		if err != nil {
//			t.Fatalf("Error while creating client: %s", err)
//		}
//
//		url, err := url.Parse(baseurl)
//		if err != nil {
//			t.Fatalf("Error while parsing url: %s", err)
//		}
//
//		url.Path = path.Join(url.Path, endpoint)
//		rateLimiter := rate.NewLimiter(rate.Every(1), 1)
//
//		cfg := web.FetchConfig{
//			C:           client,
//			Method:      "GET",
//			URL:         url,
//			RateLimiter: rateLimiter,
//		}
//
//		req := flattenedRequest{
//			fetchConfig: &cfg,
//			table:       table,
//		}
//
//		job := webJob{
//			flattenedRequest: &req,
//			repoJobs:         repoJobs,
//		}
//
//		go webWorker(context.Background(), 1, webWorkerJobs)
//
//		webWorkerJobs <- &job
//
//		close(webWorkerJobs)
//
//		for i := 0; i < 1; i++ {
//			result := <-repoJobs
//
//			if result == nil {
//				t.Fatalf("Expected repoJob, go nil")
//			}
//			if result.table != table {
//				t.Fatalf("Expected table to be %s, instead got: %s", table, result.table)
//			}
//		}
//	})
//
//	t.Run("html response with clobColumn not set", func(t *testing.T) {
//		t.Parallel()
//
//		baseurl := "https://api.cryptonator.com"
//		endpoint := "/api/ticker/btc-usd"
//		table := "test"
//
//		webWorkerJobs := make(chan *webJob, 1)
//		repoJobs := make(chan *repoJob, 1)
//
//		client, err := web.NewClient(context.Background(), nil)
//		if err != nil {
//			t.Fatalf("Error while creating client: %s", err)
//		}
//
//		url, err := url.Parse(baseurl)
//		if err != nil {
//			t.Fatalf("Error while parsing url: %s", err)
//		}
//
//		url.Path = path.Join(url.Path, endpoint)
//		rateLimiter := rate.NewLimiter(rate.Every(1), 1)
//
//		cfg := web.FetchConfig{
//			C:           client,
//			Method:      "GET",
//			URL:         url,
//			RateLimiter: rateLimiter,
//		}
//
//		req := flattenedRequest{
//			fetchConfig: &cfg,
//			table:       table,
//		}
//
//		job := webJob{
//			flattenedRequest: &req,
//			repoJobs:         repoJobs,
//		}
//
//		go webWorker(context.Background(), 1, webWorkerJobs)
//
//		webWorkerJobs <- &job
//
//		close(webWorkerJobs)
//
//		for i := 0; i < 1; i++ {
//			result := <-repoJobs
//
//			if result != nil {
//				t.Fatalf("Expected repoJob to be nil")
//			}
//		}
//	})
//
//	t.Run("html response with clobColumn set", func(t *testing.T) {
//		t.Parallel()
//
//		baseurl := "https://api.cryptonator.com"
//		endpoint := "/api/ticker/btc-usd"
//		table := "test"
//		clobColumn := "data"
//
//		webWorkerJobs := make(chan *webJob, 1)
//		repoJobs := make(chan *repoJob, 1)
//
//		client, err := web.NewClient(context.Background(), nil)
//		if err != nil {
//			t.Fatalf("Error while creating client: %s", err)
//		}
//
//		url, err := url.Parse(baseurl)
//		if err != nil {
//			t.Fatalf("Error while parsing url: %s", err)
//		}
//
//		url.Path = path.Join(url.Path, endpoint)
//		rateLimiter := rate.NewLimiter(rate.Every(1), 1)
//
//		cfg := web.FetchConfig{
//			C:           client,
//			Method:      "GET",
//			URL:         url,
//			RateLimiter: rateLimiter,
//		}
//
//		req := flattenedRequest{
//			fetchConfig: &cfg,
//			table:       table,
//			clobColumn:  clobColumn,
//		}
//
//		job := webJob{
//			flattenedRequest: &req,
//			repoJobs:         repoJobs,
//		}
//
//		go webWorker(context.Background(), 1, webWorkerJobs)
//
//		webWorkerJobs <- &job
//
//		close(webWorkerJobs)
//
//		for i := 0; i < 1; i++ {
//			result := <-repoJobs
//
//			if result == nil {
//				t.Fatalf("Expected repoJob not to be nil")
//			}
//			if result.table != table {
//				t.Fatalf("Expected table to be %s, instead got: %s", table, result.table)
//			}
//			var data interface{}
//			if err := json.Unmarshal(result.b, &data); err != nil {
//				t.Fatalf("failed to unmarshal json data")
//			}
//			dataMap, ok := data.(map[string]interface{})
//			if !ok {
//				t.Fatalf("failed to convert data to map[string]interface{}")
//			}
//
//			if _, ok := dataMap[clobColumn]; !ok {
//				t.Fatalf("expected json data to have a key: %s, but got none", clobColumn)
//			}
//		}
//	})
//}

//func TestNewFetchConfig(t *testing.T) {
//	t.Parallel()
//
//	testURL, err := url.Parse("https://fuzz")
//	if err != nil {
//		t.Fatalf("error parsing url: %v", err)
//	}
//
//	testClient := &web.Client{}
//
//	testRateLimiter := rate.NewLimiter(rate.Every(time.Second), 2)
//
//	for _, tcase := range []struct {
//		req      *Request
//		expected *web.FetchConfig
//	}{
//		{
//			req: &Request{
//				Endpoint:    "/test",
//				Method:      "POST",
//				RateLimiter: testRateLimiter,
//				Query: map[string]string{
//					"hey": "ho",
//					"foo": "bar",
//				},
//			},
//			expected: &web.FetchConfig{
//				Method:      "POST",
//				URL:         testURL,
//				C:           testClient,
//				RateLimiter: testRateLimiter,
//			},
//		},
//	} {
//		fetchConfig := newFetchConfig(tcase.req, *testURL, testClient)
//
//		if fetchConfig.Method != tcase.expected.Method {
//			t.Error("unexpected fetchMethod")
//		}
//
//		if fetchConfig.C != tcase.expected.C {
//			t.Error("unexpected fetchClient")
//		}
//
//		if fetchConfig.RateLimiter != tcase.expected.RateLimiter {
//			t.Error("unexpected fetchRateLimiter")
//		}
//
//		// URL Path should be concatenated with request Endpoint
//		if fetchConfig.URL.Path != path.Join(testURL.Path, tcase.req.Endpoint) {
//			t.Errorf("unexpected fetchURL Path: %v", fetchConfig.URL.Path)
//		}
//
//		// Query parameters should be added to the URL Query
//		for k, v := range tcase.req.Query {
//			if fetchConfig.URL.Query().Get(k) != v {
//				t.Errorf("unexpected fetchURL Query for %v", k)
//			}
//		}
//	}
//}

//func Test_flattenConfigRequests(t *testing.T) {
//	t.Parallel()
//
//	type args struct {
//		cfg *Config
//	}
//
//	tests := []struct {
//		name     string
//		args     args
//		wantReqs []*flattenedRequest
//		wantErr  bool
//	}{
//		{
//			name: "successfully flatten arequest",
//			args: args{
//				cfg: &Config{
//					URL: &url.URL{
//						Path: "path",
//					},
//					Requests: []*Request{
//						{
//							Method:   "POST",
//							Endpoint: "endpoint",
//						},
//					},
//				},
//			},
//			wantReqs: []*flattenedRequest{
//				{
//					fetchConfig: &web.FetchConfig{
//						Method: "POST",
//						URL: &url.URL{
//							Path: "path/endpoint",
//						},
//					},
//				},
//			},
//			wantErr: false,
//		},
//		{
//			name: "fail to flattenrequests when none are passed",
//			args: args{
//				cfg: &Config{},
//			},
//			wantErr: true,
//		},
//		{
//			name: "fail to flatten request time series",
//			args: args{
//				cfg: &Config{
//					URL: &url.URL{},
//					Requests: []*Request{
//						{
//							Timeseries: &Timeseries{
//								StartName: "startName",
//							},
//						},
//					},
//				},
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tcase := range tests {
//		tcase := tcase
//
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			gotReqs, err := flattenConfigRequests(context.Background(), tcase.args.cfg)
//			if (err != nil) != tcase.wantErr {
//				t.Errorf("flattenConfigRequests() error = %v, wantErr %v", err, tcase.wantErr)
//
//				return
//			}
//			if len(gotReqs) != len(tcase.wantReqs) {
//				t.Errorf("flattenConfigRequests() got %d requests, want %d", len(gotReqs), len(tcase.wantReqs))
//			}
//			for _, gotReq := range gotReqs {
//				for _, wantReq := range tcase.wantReqs {
//					if !compareRequests(gotReq, wantReq) {
//						t.Errorf("flattenConfigRequests() request comparison failed")
//					}
//				}
//			}
//		})
//	}
//}
//
//func compareRequests(request1, request2 *flattenedRequest) bool {
//	if request1 == nil || request2 == nil {
//		return false
//	}
//
//	if request1.fetchConfig.URL.Path != request2.fetchConfig.URL.Path {
//		return false
//	}
//
//	if request1.fetchConfig.Method != request2.fetchConfig.Method {
//		return false
//	}
//
//	return true
//}
