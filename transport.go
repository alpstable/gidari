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
	"io/ioutil"
	"net/url"
	"path"
	"time"

	"github.com/alpstable/gidari/internal/repository"
	"github.com/alpstable/gidari/internal/web"
	"github.com/alpstable/gidari/proto"
)

// Transport will construct the transport operation using a "transport.Config" object.
func Transport(ctx context.Context, cfg *Config) error {
	if cfg == nil {
		return ErrNilConfig
	}

	upserter, err := newUpserter(ctx, cfg)
	if err != nil {
		return fmt.Errorf("unable to create upserter: %w", err)
	}

	cfg.HTTPResponseHandler = upserter.HTTPResponseHandler

	// Create an iterator that will iterate over the requests in the configuration file.
	iter, err := NewIterator(ctx, cfg)
	if err != nil {
		return fmt.Errorf("unable to create iterator: %w", err)
	}

	defer iter.Close(ctx)

	for iter.Next(ctx) {
		// TODO - do something with the response.
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error iterating over requests: %w", err)
	}

	return nil
}

type genericRepository struct {
	repository.Generic

	// closable is a flag that indicates whether or not the repository is closable. A repository is only closable
	// if it is created by a connection string. If a repository is created by a client or a database, then it is
	// the responsibility of the caller to close.
	closable bool
	database string
}

func (gen *genericRepository) Close() {
	if gen.closable {
		gen.Generic.Close()
	}
}

func newGenericRepositories(ctx context.Context, cfg *Config) ([]genericRepository, error) {
	genericRepositories := make([]genericRepository, len(cfg.StorageOptions))

	for idx, stgOpts := range cfg.StorageOptions {
		stg := stgOpts.Storage

		genRepository, err := repository.NewTx(ctx, stg)
		if err != nil {
			return nil, fmt.Errorf("failed to create repository: %w", err)
		}

		genericRepositories[idx] = genericRepository{
			Generic:  genRepository,
			database: stgOpts.Database,
			closable: stgOpts.Close,
		}
	}

	return genericRepositories, nil
}

// newFetchConfig will construct a new HTTP request from the transport request.
func newFetchConfig(req *Request, rurl url.URL, client *web.Client) *web.FetchConfig {
	rurl.Path = path.Join(rurl.Path, req.Endpoint)

	// Add the query params to the URL.
	if req.Query != nil {
		query := rurl.Query()
		for key, value := range req.Query {
			query.Set(key, value)
		}

		rurl.RawQuery = query.Encode()
	}

	return &web.FetchConfig{
		Method:      req.Method,
		URL:         &rurl,
		C:           client,
		RateLimiter: req.RateLimiter,
	}
}

// flattenedRequest contains all of the request information to create a web job. The number of flattened request  for an
// operation should be 1-1 with the number of requests to the web API.
type flattenedRequest struct {
	// fetchConfig is the configuration for the HTTP request. Each request gets it's own connection to ensure that
	// the web worker can process concurrently without locking. Despite this, however, all of the requests should
	// share a common rate limiter to prevent overloading the web API and gettig a 429 response.
	fetchConfig *web.FetchConfig
	table       string
	clobColumn  string
}

// flatten will compress the request information into a "web.FetchConfig" request and a "table" name for storage
// interaction.
func flattenRequest(req *Request, rurl url.URL, client *web.Client) *flattenedRequest {
	fetchConfig := newFetchConfig(req, rurl, client)

	// If the table is not set on the request, then set it using the last path segment of the endpoint.
	if req.Table == "" {
		req.Table = path.Base(req.Endpoint)
	}

	return &flattenedRequest{
		fetchConfig: fetchConfig,
		table:       req.Table,
		clobColumn:  req.ClobColumn,
	}
}

// chunkTimeseries will attempt to use the query string of a URL to partition the timeseries into "Chunks" of time for
// queying a web API.
func chunkTimeseries(timeseries *Timeseries, rurl url.URL) error {
	// If layout is not set, then default it to be RFC3339
	if timeseries.Layout == "" {
		timeseries.Layout = time.RFC3339
	}

	query := rurl.Query()

	startSlice := query[timeseries.StartName]
	if len(startSlice) != 1 {
		return ErrInvalidStartTimeSize
	}

	start, err := time.Parse(timeseries.Layout, startSlice[0])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}

	endSlice := query[timeseries.EndName]
	if len(endSlice) != 1 {
		return ErrInvalidEndTimeSize
	}

	end, err := time.Parse(timeseries.Layout, endSlice[0])
	if err != nil {
		return fmt.Errorf("unable to parse end time: %w", err)
	}

	for start.Before(end) {
		next := start.Add(time.Second * time.Duration(timeseries.Period))
		if next.Before(end) {
			timeseries.Chunks = append(timeseries.Chunks, [2]time.Time{start, next})
		} else {
			timeseries.Chunks = append(timeseries.Chunks, [2]time.Time{start, end})
		}

		start = next
	}

	return nil
}

// flattenRequestTimeseries will compress the request information into a "web.FetchConfig" request and a "table" name
// for storage interaction. This function will create a flattened request for each time series in the request. If no
// timeseries are defined, this function will return a single flattened request.
func flattenRequestTimeseries(req *Request, rurl url.URL, client *web.Client) ([]*flattenedRequest, error) {
	timeseries := req.Timeseries
	if timeseries == nil {
		flatReq := flattenRequest(req, rurl, client)

		return []*flattenedRequest{flatReq}, nil
	}

	requests := make([]*flattenedRequest, 0, len(timeseries.Chunks))

	// Add the query params to the URL.
	if req.Query != nil {
		query := rurl.Query()
		for key, value := range req.Query {
			query.Set(key, value)
		}

		rurl.RawQuery = query.Encode()
	}

	if err := chunkTimeseries(timeseries, rurl); err != nil {
		return nil, fmt.Errorf("failed to set time series chunks: %w", err)
	}

	for _, chunk := range timeseries.Chunks {
		// copy the request and update it to reflect the partitioned timeseries
		chunkReq := req
		chunkReq.Query[timeseries.StartName] = chunk[0].Format(timeseries.Layout)
		chunkReq.Query[timeseries.EndName] = chunk[1].Format(timeseries.Layout)

		fetchConfig := newFetchConfig(chunkReq, rurl, client)

		requests = append(requests, &flattenedRequest{
			fetchConfig: fetchConfig,
			table:       req.Table,
			clobColumn:  req.ClobColumn,
		})
	}

	return requests, nil
}

// flattenConfigRequests will flatten the requests into a single slice for HTTP requests.
func flattenConfigRequests(ctx context.Context, cfg *Config) ([]*flattenedRequest, error) {
	client, err := connect(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to web API: %w", err)
	}

	var flattenedRequests []*flattenedRequest

	for _, req := range cfg.Requests {
		flatReqs, err := flattenRequestTimeseries(req, *cfg.URL, client)
		if err != nil {
			return nil, err
		}

		flattenedRequests = append(flattenedRequests, flatReqs...)
	}

	if len(flattenedRequests) == 0 {
		return nil, ErrNoRequests
	}

	return flattenedRequests, nil
}

func upsert_(_ context.Context, req *proto.UpsertRequest, genericRepository genericRepository) error {
	txfn := func(sctx context.Context, repo repository.Generic) error {
		_, err := repo.Upsert(sctx, req)
		if err != nil {
			return fmt.Errorf("error upserting data: %w", err)
		}

		return nil
	}

	genericRepository.Transact(txfn)

	return nil
}

type upserter struct {
	// repositories is a list of repositories to upsert data into, defined by a configuration.
	repositories []genericRepository

	repositoryChan chan genericRepository
}

func newUpserter(ctx context.Context, cfg *Config) (*upserter, error) {
	upserter := new(upserter)

	var err error

	upserter.repositories, err = newGenericRepositories(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return upserter, nil
}

func (upserter *upserter) HTTPResponseHandler(ctx context.Context, rsp HTTPResponse) ([]*proto.IteratorResult, error) {
	errChan := make(chan error, len(upserter.repositories))

	// Get bytes from response body.
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	defer rsp.Body.Close()

	for _, repo := range upserter.repositories {
		go func(repo genericRepository) {
			req := &proto.UpsertRequest{
				Table: &proto.Table{Name: rsp.TableName, Database: repo.database},
				Data:  body,
			}

			err := upsert_(ctx, req, repo)
			if err != nil {
				errChan <- err
			}

			errChan <- nil
		}(repo)
	}

	for idx := 0; idx < len(upserter.repositories); idx++ {
		err := <-errChan
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
