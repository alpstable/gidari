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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/alpstable/gidari/config"
	"github.com/alpstable/gidari/internal/repository"
	"github.com/alpstable/gidari/internal/web"
	"github.com/alpstable/gidari/internal/web/auth"
	"github.com/alpstable/gidari/proto"
	"github.com/alpstable/gidari/tools"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidEndTimeSize   = fmt.Errorf("invalid end time size, expected 1")
	ErrInvalidStartTimeSize = fmt.Errorf("invalid start time size, expected 1")
)

// connect will attempt to connect to the web API client. Since there are multiple ways to build a transport given the
// authentication data, this method will exhaust every transport option in the "Authentication" struct.
func connect(ctx context.Context, cfg *config.Config) (*web.Client, error) {
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
	database *string
}

// repos will return a slice of generic repositories along with associated transaction instances.
func repos(ctx context.Context, cfg *config.Config) ([]repo, repoCloser, error) {
	repos := []repo{}

	for _, stgOpts := range cfg.StorageOptions {
		stg := stgOpts.Storage
		genRepository, err := repository.NewTx(ctx, stg)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create repository: %w", err)
		}

		logInfo := tools.LogFormatter{
			Msg: fmt.Sprintf("created repository for %q", stg.Type()),
		}
		cfg.Logger.Info(logInfo.String())

		repos = append(repos, repo{
			Generic:  genRepository,
			database: stgOpts.Database,
			closable: stgOpts.ConnectionString == nil,
		})
	}

	return repos, func() {
		for _, repo := range repos {
			if !repo.closable {
				continue
			}

			repo.Close()

			logInfo := tools.LogFormatter{
				Msg: fmt.Sprintf("closed repository for %q", proto.SchemeFromStorageType(repo.Type())),
			}
			cfg.Logger.Info(logInfo.String())
		}
	}, nil
}

// newFetchConfig will construct a new HTTP request from the transport request.
func newFetchConfig(req *config.Request, rurl url.URL, client *web.Client) *web.FetchConfig {
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
func flattenRequest(req *config.Request, rurl url.URL, client *web.Client) *flattenedRequest {
	fetchConfig := newFetchConfig(req, rurl, client)

	return &flattenedRequest{
		fetchConfig: fetchConfig,
		table:       req.Table,
		clobColumn:  req.ClobColumn,
	}
}

// chunkTimeseries will attempt to use the query string of a URL to partition the timeseries into "Chunks" of time for
// queying a web API.
func chunkTimeseries(timeseries *config.Timeseries, rurl url.URL) error {
	// If layout is not set, then default it to be RFC3339
	if timeseries.Layout == nil {
		str := time.RFC3339
		timeseries.Layout = &str
	}

	query := rurl.Query()

	startSlice := query[timeseries.StartName]
	if len(startSlice) != 1 {
		return ErrInvalidStartTimeSize
	}

	start, err := time.Parse(*timeseries.Layout, startSlice[0])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %w", err)
	}

	endSlice := query[timeseries.EndName]
	if len(endSlice) != 1 {
		return ErrInvalidEndTimeSize
	}

	end, err := time.Parse(*timeseries.Layout, endSlice[0])
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
func flattenRequestTimeseries(req *config.Request, rurl url.URL, client *web.Client) ([]*flattenedRequest, error) {
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
		chunkReq.Query[timeseries.StartName] = chunk[0].Format(*timeseries.Layout)
		chunkReq.Query[timeseries.EndName] = chunk[1].Format(*timeseries.Layout)

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
func flattenConfigRequests(ctx context.Context, cfg *config.Config) ([]*flattenedRequest, error) {
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
		return nil, config.ErrNoRequests
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
	logger     *logrus.Logger
	database   *string
}

func newRepoConfig(ctx context.Context, cfg *config.Config, volume int) (*repoConfig, error) {
	repos, closeRepos, err := repos(ctx, cfg)
	if err != nil {
		return nil, err
	}

	config := &repoConfig{
		repos:      repos,
		closeRepos: closeRepos,
		jobs:       make(chan *repoJob, volume*len(repos)),
		done:       make(chan bool, volume),
		logger:     cfg.Logger,
	}

	return config, nil

}

func repositoryWorker(_ context.Context, workerID int, cfg *repoConfig) {
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
				// If the database is set, then set it on the request.
				if database := repo.database; database != nil {
					req.Table.Database = *database
				}

				txfn := func(sctx context.Context, repo repository.Generic) error {
					start := time.Now()

					rsp, err := repo.Upsert(sctx, req)
					if err != nil {
						cfg.logger.Fatalf("error upserting data: %v", err)

						return fmt.Errorf("error upserting data: %w", err)
					}

					rt := repo.Type()

					msg := fmt.Sprintf("partial upsert completed: %s.%s", proto.SchemeFromStorageType(rt), req.Table)
					logInfo := tools.LogFormatter{
						WorkerID:      workerID,
						WorkerName:    "repository",
						Duration:      time.Since(start),
						Msg:           msg,
						UpsertedCount: rsp.UpsertedCount,
						MatchedCount:  rsp.MatchedCount,
					}

					cfg.logger.Infof(logInfo.String())

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
	logger   *logrus.Logger
}

func newWebJob(cfg *config.Config, req *flattenedRequest, repoJobs chan<- *repoJob) *webJob {
	return &webJob{
		flattenedRequest: req,
		repoJobs:         repoJobs,
		logger:           cfg.Logger,
	}
}

func webWorker(ctx context.Context, workerID int, jobs <-chan *webJob) {
	for job := range jobs {
		start := time.Now()

		rsp, err := web.Fetch(ctx, job.fetchConfig)
		if err != nil {
			job.logger.Fatal(err)
		}

		bytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			job.logger.Fatal(err)
		}

		if !json.Valid(bytes) {
			if job.flattenedRequest.clobColumn == "" {
				job.repoJobs <- nil
				msg := fmt.Sprintf("response body for %s was invalid JSON, "+
					"discarding data since no 'clobColumn' was defined in the configuration file",
					job.fetchConfig.URL)
				logInfo := tools.LogFormatter{Msg: msg}
				job.logger.Warnf(logInfo.String())

				continue
			}

			data := make(map[string]string)
			data[job.flattenedRequest.clobColumn] = string(bytes)

			bytes, err = json.Marshal(data)
			if err != nil {
				job.repoJobs <- nil
				job.logger.Errorf("failed to marhsal data: %s", err)
			}
		}

		job.repoJobs <- &repoJob{b: bytes, req: *rsp.Request, table: job.table}

		// strings.Replace is used to ensure no line endings are present in the user input.
		escapedPath := strings.ReplaceAll(rsp.Request.URL.Path, "\n", "")
		escapedPath = strings.ReplaceAll(escapedPath, "\r", "")

		escapedHost := strings.ReplaceAll(rsp.Request.URL.Host, "\n", "")
		escapedHost = strings.ReplaceAll(escapedHost, "\r", "")

		logInfo := tools.LogFormatter{
			WorkerID:   workerID,
			WorkerName: "web",
			Duration:   time.Since(start),
			Host:       escapedHost,
			Msg:        fmt.Sprintf("web request completed: %s", escapedPath),
		}
		job.logger.Infof(logInfo.String())
	}
}

func truncate(ctx context.Context, cfg *config.Config, truncateRequest *proto.TruncateRequest) error {
	start := time.Now()

	repos, closeRepos, err := repos(ctx, cfg)
	if err != nil {
		return err
	}

	defer closeRepos()

	for _, repo := range repos {
		start := time.Now()

		// If the database is set, then set it on the request.
		if database := repo.database; database != nil {
			for _, table := range truncateRequest.Tables {
				table.Database = *database
			}
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

		rt := repo.Type()
		tablesCSV := strings.Join(tableNames, ", ")
		msg := fmt.Sprintf("truncated tables on %q: %v", proto.SchemeFromStorageType(rt), tablesCSV)

		logInfo := tools.LogFormatter{
			Duration: time.Since(start),
			Msg:      msg,
		}
		cfg.Logger.Infof(logInfo.String())
	}

	logInfo := tools.LogFormatter{
		Duration: time.Since(start),
		Msg:      "truncate completed",
	}
	cfg.Logger.Info(logInfo.String())

	return nil
}

// Truncate will truncate the defined tables in the configuration.
func Truncate(ctx context.Context, cfg *config.Config) error {
	// truncateRequest is a special request that will truncate the table before upserting data.
	truncateRequest := new(proto.TruncateRequest)

	if cfg.Truncate {
		for _, req := range cfg.Requests {
			// Add the table to the list of tables to truncate.
			if req.Truncate != nil && *req.Truncate {
				truncateRequest.Tables = append(truncateRequest.Tables, &proto.Table{Name: req.Table})
			}
		}
	} else {
		// checking for request-specific truncate
		for _, req := range cfg.Requests {
			if table := req.Table; req.Truncate != nil && *req.Truncate && table != "" {
				truncateRequest.Tables = append(truncateRequest.Tables, &proto.Table{Name: table})
			}
		}
	}

	return truncate(ctx, cfg, truncateRequest)
}

// Upsert will use the configuration file to upsert data from the
//
// For each DNS entry in the configuration file, a repository will be created and used to upsert data. For each
// repository, a transaction will be created and used to upsert data. The transaction will be committed at the end
// of the upsert operation. If the transaction fails, the transaction will be rolled back. Note that it is possible
// for some repository transactions to succeed and others to fail.
func Upsert(ctx context.Context, cfg *config.Config) error {
	start := time.Now()
	threads := runtime.NumCPU()

	if err := Truncate(ctx, cfg); err != nil {
		return err
	}

	flattenedRequests, err := flattenConfigRequests(ctx, cfg)
	if err != nil {
		return err
	}

	repoConfig, err := newRepoConfig(ctx, cfg, len(flattenedRequests))
	if err != nil {
		return err
	}

	defer func() {
		repoConfig.closeRepos()
	}()

	// Start the repository workers.
	for id := 1; id <= threads; id++ {
		go repositoryWorker(ctx, id, repoConfig)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "repository workers started"}.String())

	webWorkerJobs := make(chan *webJob, len(cfg.Requests))

	// Start the same number of web workers as the cores on the machine.
	for id := 1; id <= threads; id++ {
		go webWorker(ctx, id, webWorkerJobs)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "web workers started"}.String())

	// Enqueue the worker jobs
	for _, req := range flattenedRequests {
		webWorkerJobs <- newWebJob(cfg, req, repoConfig.jobs)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "web worker jobs enqueued"}.String())

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

	logInfo := tools.LogFormatter{Duration: time.Since(start), Msg: "upsert completed"}
	cfg.Logger.Info(logInfo.String())

	return nil
}
