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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/alpstable/gidari/internal/repository"
	"github.com/alpstable/gidari/internal/web"
	"github.com/alpstable/gidari/internal/web/auth"
	"github.com/alpstable/gidari/proto"
)

// connect will attempt to connect to the web API client. Since there are multiple ways to build a transport given the
// authentication data, this method will exhaust every transport option in the "Authentication" struct.
func connect(ctx context.Context, cfg *Config) (*web.Client, error) {
	if apiKey := cfg.Authentication.APIKey; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAPIKey().
			SetURL(cfg.RawURL).
			SetKey(apiKey.Key).
			SetPassphrase(apiKey.Passphrase).
			SetSecret(apiKey.Secret))
		if err != nil {
			return nil, fmt.Errorf("failed to create API key client: %w", err)
		}

		return client, nil
	}

	if apiKey := cfg.Authentication.Auth2; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAuth2().SetBearer(apiKey.Bearer).SetURL(cfg.RawURL))
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}

		return client, nil
	}

	// In the case of no authentication, create a client without an auth transport.
	client, err := web.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

type repoCloser func()

type repo struct {
	repository.Generic

	// closable is a flag that indicates whether or not the repository is closable. A repository is only closable
	// if it is created by a connection string. If a repository is created by a client or a database, then it is
	// the responsibility of the caller to close.
	closable bool
	database string
}

// repos will return a slice of generic repositories along with associated transaction instances.
func repos(ctx context.Context, cfg *Config) ([]repo, repoCloser, error) {
	repos := []repo{}

	for _, stgOpts := range cfg.StorageOptions {
		stg := stgOpts.Storage

		genRepository, err := repository.NewTx(ctx, stg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create repository: %w", err)
		}

		repos = append(repos, repo{
			Generic:  genRepository,
			database: stgOpts.Database,
			closable: stgOpts.Close,
		})
	}

	return repos, func() {
		for _, repo := range repos {
			if !repo.closable {
				continue
			}

			repo.Close()
		}
	}, nil
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

type repoJob struct {
	req   http.Request
	b     []byte
	table string
}

type repoConfig struct {
	repos      []repo
	closeRepos func()
	jobs       chan *repoJob
	done       chan bool
}

func newRepoConfig(ctx context.Context, cfg *Config, volume int) (*repoConfig, error) {
	repos, closeRepos, err := repos(ctx, cfg)
	if err != nil {
		return nil, err
	}

	config := &repoConfig{
		repos:      repos,
		closeRepos: closeRepos,
		jobs:       make(chan *repoJob, volume*len(repos)),
		done:       make(chan bool, volume),
	}

	return config, nil
}

func repositoryWorker(_ context.Context, _ int, cfg *repoConfig) {
	for job := range cfg.jobs {
		if job == nil {
			cfg.done <- false

			continue
		}

		reqs := []*proto.UpsertRequest{
			{
				Table: &proto.Table{Name: job.table},
				Data:  job.b,
			},
		}

		for _, req := range reqs {
			for _, repo := range cfg.repos {
				req.Table.Database = repo.database

				txfn := func(sctx context.Context, repo repository.Generic) error {
					_, err := repo.Upsert(sctx, req)
					if err != nil {
						return fmt.Errorf("error upserting data: %w", err)
					}

					return nil
				}

				// Put the data onto the transaction channel for storage.
				repo.Transact(txfn)
			}
		}

		cfg.done <- true
	}
}

type webJob struct {
	*flattenedRequest
	repoJobs chan<- *repoJob
}

func newWebJob(_ *Config, req *flattenedRequest, repoJobs chan<- *repoJob) *webJob {
	return &webJob{
		flattenedRequest: req,
		repoJobs:         repoJobs,
	}
}

func webWorker(ctx context.Context, _ int, jobs <-chan *webJob) {
	for job := range jobs {
		rsp, err := web.Fetch(ctx, job.fetchConfig)
		if err != nil {
			log.Fatalf("error fetching data: %v", err)
		}

		bytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Fatalf("error reading response body: %v", err)
		}

		if !json.Valid(bytes) {
			if job.flattenedRequest.clobColumn == "" {
				job.repoJobs <- nil

				continue
			}

			data := make(map[string]string)
			data[job.flattenedRequest.clobColumn] = string(bytes)

			bytes, err = json.Marshal(data)
			if err != nil {
				job.repoJobs <- nil

				log.Fatalf("failed to marhsal data: %v", err)
			}
		}

		job.repoJobs <- &repoJob{b: bytes, req: *rsp.Request, table: job.table}
	}
}

// setTruncateRequestTables will set the tables data to truncate from the configuration structure. If "Truncate" is set
// on the configuration, then all tables in the request will be truncated before the data is upserted.
func setTruncateRequestTables(treq *proto.TruncateRequest, cfg *Config) {
	for _, req := range cfg.Requests {
		if (req.Truncate != nil && *req.Truncate && req.Table != "") || cfg.Truncate {
			treq.Tables = append(treq.Tables, &proto.Table{Name: req.Table})
		}
	}
}

// Truncate will truncate the defined tables in the configuration.
func truncate(ctx context.Context, cfg *Config, repoConfig *repoConfig) error {
	// truncateRequest is a special request that will truncate the table before upserting data.
	truncateRequest := new(proto.TruncateRequest)
	setTruncateRequestTables(truncateRequest, cfg)

	for _, repo := range repoConfig.repos {
		for _, table := range truncateRequest.Tables {
			table.Database = repo.database
		}

		_, err := repo.Truncate(ctx, truncateRequest)
		if err != nil {
			return fmt.Errorf("unable to truncate tables: %w", err)
		}

		tables := truncateRequest.GetTables()

		tableNames := make([]string, len(tables))
		for idx, table := range tables {
			tableNames[idx] = table.GetName()
		}
	}

	return nil
}

// upsert will use the configuration file to upsert data from the
//
// For each DNS entry in the configuration file, a repository will be created and used to upsert data. For each
// repository, a transaction will be created and used to upsert data. The transaction will be committed at the end
// of the upsert operation. If the transaction fails, the transaction will be rolled back. Note that it is possible
// for some repository transactions to succeed and others to fail.
func upsert(ctx context.Context, cfg *Config) error {
	threads := runtime.NumCPU()

	flattenedRequests, err := flattenConfigRequests(ctx, cfg)
	if err != nil {
		return err
	}

	repoConfig, err := newRepoConfig(ctx, cfg, len(flattenedRequests))
	if err != nil {
		return err
	}

	if err := truncate(ctx, cfg, repoConfig); err != nil {
		return err
	}

	defer func() {
		repoConfig.closeRepos()
	}()

	// Start the repository workers.
	for id := 1; id <= threads; id++ {
		go repositoryWorker(ctx, id, repoConfig)
	}

	webWorkerJobs := make(chan *webJob, len(cfg.Requests))

	// Start the same number of web workers as the cores on the machine.
	for id := 1; id <= threads; id++ {
		go webWorker(ctx, id, webWorkerJobs)
	}

	// Enqueue the worker jobs
	for _, req := range flattenedRequests {
		webWorkerJobs <- newWebJob(cfg, req, repoConfig.jobs)
	}

	// Wait for all of the data to flush.
	for a := 1; a <= len(flattenedRequests); a++ {
		<-repoConfig.done
	}

	// Commit the transactions and check for errors.
	for _, repo := range repoConfig.repos {
		if err := repo.Commit(); err != nil {
			return fmt.Errorf("unable to commit transaction: %w", err)
		}
	}

	return nil
}

// Transport will construct the transport operation using a "transport.Config" object.
func Transport(ctx context.Context, cfg *Config) error {
	if err := upsert(ctx, cfg); err != nil {
		return fmt.Errorf("unable to upsert the config: %w", err)
	}

	return nil
}
