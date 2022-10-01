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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/internal/web"
	"github.com/alpine-hodler/gidari/internal/web/auth"
	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/repository"
	"github.com/alpine-hodler/gidari/tools"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	ErrFetchingTimeseriesChunks = fmt.Errorf("failed to fetch timeseries chunks")
	ErrInvalidRateLimit         = fmt.Errorf("invalid rate limit configuration")
	ErrMissingConfigField       = fmt.Errorf("missing config field")
	ErrMissingRateLimitField    = fmt.Errorf("missing rate limit field")
	ErrMissingTimeseriesField   = fmt.Errorf("missing timeseries field")
	ErrSettingTimeseriesChunks  = fmt.Errorf("failed to set timeseries chunks")
	ErrUnableToParse            = fmt.Errorf("unable to parse")
	ErrNoRequests               = fmt.Errorf("no requests defined")
)

// MissingConfigFieldError is returned when a configuration field is missing.
func MissingConfigFieldError(field string) error {
	return fmt.Errorf("%w: %s", ErrMissingConfigField, field)
}

// MissingRateLimitFieldError is returned when the rate limit configuration is missing a field.
func MissingRateLimitFieldError(field string) error {
	return fmt.Errorf("%w: %s", ErrMissingRateLimitField, field)
}

// MissingTimeseriesFieldError is returned when the timeseries is missing from the configuration.
func MissingTimeseriesFieldError(field string) error {
	return fmt.Errorf("%w: %s", ErrMissingTimeseriesField, field)
}

// UnableToParseError is returned when a parser is unable to parse the data.
func UnableToParseError(name string) error {
	return fmt.Errorf("%s %w", name, ErrUnableToParse)
}

// WrapRepositoryError will wrap an error from the repository with a message.
func WrapRepositoryError(err error) error {
	return fmt.Errorf("repository: %w", err)
}

// WrapWebError will wrap an error from the web package with a message.
func WrapWebError(err error) error {
	return fmt.Errorf("web: %w", err)
}

// APIKey is one method of HTTP(s) transport that requires a passphrase, key, and secret.
type APIKey struct {
	Passphrase string `yaml:"passphrase"`
	Key        string `yaml:"key"`
	Secret     string `yaml:"secret"`
}

// Auth2 is a struct that contains the authentication data for a web API that uses OAuth2.
type Auth2 struct {
	Bearer string `yaml:"bearer"`
}

// Authentication is the credential information to be used to construct an HTTP(s) transport for accessing the API.
type Authentication struct {
	APIKey *APIKey `yaml:"apiKey"`
	Auth2  *Auth2  `yaml:"auth2"`
}

// timeseries is a struct that contains the information needed to query a web API for timeseries data.
type timeseries struct {
	StartName string `yaml:"startName"`
	EndName   string `yaml:"endName"`

	// Period is the size of each chunk in seconds for which we can query the API. Some API will not allow us to
	// query all data within the start and end range.
	Period int32 `yaml:"period"`

	// Layout is the time layout for parsing the "Start" and "End" values into "time.Time". The default is assumed
	// to be RFC3339.
	Layout *string `yaml:"layout"`

	// chunks are the time ranges for which we can query the API. These are broken up into pieces for API requests
	// that only return a limited number of results.
	chunks [][2]time.Time
}

// chunk will attempt to use the query string of a URL to partition the timeseries into "chunks" of time for queying
// a web API.
func (ts *timeseries) chunk(rurl url.URL) error {
	// If layout is not set, then default it to be RFC3339
	if ts.Layout == nil {
		str := time.RFC3339
		ts.Layout = &str
	}

	query := rurl.Query()

	startSlice := query[ts.StartName]
	if len(startSlice) != 1 {
		return MissingTimeseriesFieldError("startName")
	}

	start, err := time.Parse(*ts.Layout, startSlice[0])
	if err != nil {
		return UnableToParseError("startTime")
	}

	endSlice := query[ts.EndName]
	if len(endSlice) != 1 {
		return MissingTimeseriesFieldError("endName")
	}

	end, err := time.Parse(*ts.Layout, endSlice[0])
	if err != nil {
		return UnableToParseError("endTime")
	}

	for start.Before(end) {
		next := start.Add(time.Second * time.Duration(ts.Period))
		if next.Before(end) {
			ts.chunks = append(ts.chunks, [2]time.Time{start, next})
		} else {
			ts.chunks = append(ts.chunks, [2]time.Time{start, end})
		}

		start = next
	}

	return nil
}

// RateLimitConfig is the data needed for constructing a rate limit for the HTTP requests.
type RateLimitConfig struct {
	// Burst represents the number of requests that we limit over a period frequency.
	Burst *int `yaml:"burst"`

	// Period is the number of times to allow a burst per second.
	Period *time.Duration `yaml:"period"`
}

func (rl RateLimitConfig) validate() error {
	if rl.Burst == nil {
		return MissingRateLimitFieldError("burst")
	}

	if rl.Period == nil {
		return MissingRateLimitFieldError("period")
	}

	return nil
}

// Config is the configuration used to query data from the web using HTTP requests and storing that data using
// the repositories defined by the "ConnectionStrings" list.
type Config struct {
	RawURL            string           `yaml:"url"`
	Authentication    Authentication   `yaml:"authentication"`
	ConnectionStrings []string         `yaml:"connectionStrings"`
	Requests          []*Request       `yaml:"requests"`
	RateLimitConfig   *RateLimitConfig `yaml:"rateLimit"`
	Logger            *logrus.Logger
	Truncate          bool

	URL *url.URL `yaml:"-"`
}

// New config takes a YAML byte slice and returns a new transport configuration for upserting data to storage.
//
// For web requests defined on the transport configuration, the default HTTP Request Method is "GET". Furthermore,
// if rate limit data has not been defined for a request it will inherit the rate limit data from the transport config.
func NewConfig(yamlBytes []byte) (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(yamlBytes, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// Parse the raw URL
	var err error

	cfg.URL, err = url.Parse(cfg.RawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %w", err)
	}

	// Update default request data.
	for _, req := range cfg.Requests {
		if req.Method == "" {
			req.Method = http.MethodGet
		}

		if req.RateLimitConfig == nil {
			req.RateLimitConfig = cfg.RateLimitConfig
		}

		if req.Table == "" {
			endpointParts := strings.Split(req.Endpoint, "/")
			req.Table = endpointParts[len(endpointParts)-1]
		}
	}

	return &cfg, nil
}

// connect will attempt to connect to the web API client. Since there are multiple ways to build a transport given the
// authentication data, this method will exhaust every transport option in the "Authentication" struct.
func (cfg *Config) connect(ctx context.Context) (*web.Client, error) {
	if apiKey := cfg.Authentication.APIKey; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAPIKey().
			SetURL(cfg.RawURL).
			SetKey(apiKey.Key).
			SetPassphrase(apiKey.Passphrase).
			SetSecret(apiKey.Secret))
		if err != nil {
			return nil, WrapWebError(web.FailedToCreateClientError(err))
		}

		return client, nil
	}

	if apiKey := cfg.Authentication.Auth2; apiKey != nil {
		client, err := web.NewClient(ctx, auth.NewAuth2().SetBearer(apiKey.Bearer).SetURL(cfg.RawURL))
		if err != nil {
			return nil, WrapWebError(web.FailedToCreateClientError(err))
		}

		return client, nil
	}

	// In the case of no authentication, create a client without an auth transport.
	client, err := web.NewClient(ctx, nil)
	if err != nil {
		return nil, WrapWebError(web.FailedToCreateClientError(err))
	}

	return client, nil
}

type repoCloser func()

// repos will return a slice of generic repositories along with associated transaction instances.
func (cfg *Config) repos(ctx context.Context) ([]repository.Generic, repoCloser, error) {
	repos := []repository.Generic{}

	for _, dns := range cfg.ConnectionStrings {
		repo, err := repository.NewTx(ctx, dns)
		if err != nil {
			return nil, nil, WrapRepositoryError(repository.FailedToCreateRepositoryError(err))
		}

		logInfo := tools.LogFormatter{
			Msg: fmt.Sprintf("created repository for %q", dns),
		}
		cfg.Logger.Info(logInfo.String())

		repos = append(repos, repo)
	}

	return repos, func() {
		for _, repo := range repos {
			repo.Close()

			logInfo := tools.LogFormatter{
				Msg: fmt.Sprintf("closed repository for %q", storage.Scheme(repo.Type())),
			}
			cfg.Logger.Info(logInfo.String())
		}
	}, nil
}

// validate will ensure that the configuration is valid for querying the web API.
func (cfg *Config) validate() error {
	if cfg.RateLimitConfig == nil {
		return MissingConfigFieldError("rateLimit")
	}

	if err := cfg.RateLimitConfig.validate(); err != nil {
		return ErrInvalidRateLimit
	}

	return nil
}

// flattenRequests will flatten the requests into a single slice for HTTP requests.
func (cfg *Config) flattenRequests(ctx context.Context) ([]*flattenedRequest, error) {
	client, err := cfg.connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to web API: %w", err)
	}

	var flattenedRequests []*flattenedRequest

	for _, req := range cfg.Requests {
		flatReqs, err := req.flattenTimeseries(*cfg.URL, client)
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
	repos      []repository.Generic
	closeRepos repoCloser
	jobs       chan *repoJob
	done       chan bool
	logger     *logrus.Logger
}

func newRepoConfig(ctx context.Context, cfg *Config, volume int) (*repoConfig, error) {
	repos, closeRepos, err := cfg.repos(ctx)
	if err != nil {
		return nil, err
	}

	return &repoConfig{
		repos:      repos,
		closeRepos: closeRepos,
		jobs:       make(chan *repoJob, volume*len(repos)),
		done:       make(chan bool, volume),
		logger:     cfg.Logger,
	}, nil
}

func repositoryWorker(_ context.Context, workerID int, cfg *repoConfig) {
	for job := range cfg.jobs {
		reqs := []*proto.UpsertRequest{
			{
				Table:    job.table,
				Data:     job.b,
				DataType: int32(tools.UpsertDataJSON),
			},
		}

		for _, req := range reqs {
			for _, repo := range cfg.repos {
				txfn := func(sctx context.Context, repo repository.Generic) error {
					start := time.Now()

					rsp, err := repo.Upsert(sctx, req)
					if err != nil {
						cfg.logger.Fatalf("error upserting data: %v", err)

						return fmt.Errorf("error upserting data: %w", err)
					}

					rt := repo.Type()

					msg := fmt.Sprintf("partial upsert completed: %s.%s", storage.Scheme(rt), req.Table)
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

func newWebJob(cfg *Config, req *flattenedRequest, repoJobs chan<- *repoJob) *webJob {
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

		job.repoJobs <- &repoJob{b: bytes, req: *rsp.Request, table: job.table}

		logInfo := tools.LogFormatter{
			WorkerID:   workerID,
			WorkerName: "web",
			Duration:   time.Since(start),
			Msg:        fmt.Sprintf("web request completed: %s", rsp.Request.URL.Path),
		}
		job.logger.Infof(logInfo.String())
	}
}

// Truncate will truncate the defined tables in the configuration.
func Truncate(ctx context.Context, cfg *Config) error {
	if !cfg.Truncate {
		return nil
	}

	start := time.Now()

	repos, closeRepos, err := cfg.repos(ctx)
	if err != nil {
		return err
	}

	defer closeRepos()

	// truncateRequest is a special request that will truncate the table before upserting data.
	truncateRequest := new(proto.TruncateRequest)

	for _, req := range cfg.Requests {
		// Add the table to the list of tables to truncate.
		truncateRequest.Tables = append(truncateRequest.Tables, req.Table)
	}

	for _, repo := range repos {
		start := time.Now()

		_, err := repo.Truncate(ctx, truncateRequest)
		if err != nil {
			return fmt.Errorf("unable to truncate tables: %w", err)
		}

		rt := repo.Type()
		tables := strings.Join(truncateRequest.Tables, ", ")
		msg := fmt.Sprintf("truncated tables on %q: %v", storage.Scheme(rt), tables)

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

// Upsert will use the configuration file to upsert data from the
//
// For each DNS entry in the configuration file, a repository will be created and used to upsert data. For each
// repository, a transaction will be created and used to upsert data. The transaction will be committed at the end
// of the upsert operation. If the transaction fails, the transaction will be rolled back. Note that it is possible
// for some repository transactions to succeed and others to fail.
func Upsert(ctx context.Context, cfg *Config) error {
	start := time.Now()

	if err := Truncate(ctx, cfg); err != nil {
		return err
	}

	flattenedRequests, err := cfg.flattenRequests(ctx)
	if err != nil {
		return err
	}

	repoConfig, err := newRepoConfig(ctx, cfg, len(flattenedRequests))
	if err != nil {
		return err
	}

	defer repoConfig.closeRepos()

	// Start the repository workers.
	for id := 1; id <= runtime.NumCPU(); id++ {
		go repositoryWorker(ctx, id, repoConfig)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "repository workers started"}.String())

	webWorkerJobs := make(chan *webJob, len(cfg.Requests))

	// Start the same number of web workers as the cores on the machine.
	for id := 1; id <= runtime.NumCPU(); id++ {
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

	logInfo := tools.LogFormatter{
		Duration: time.Since(start),
		Msg:      "upsert completed",
	}
	cfg.Logger.Info(logInfo.String())

	return nil
}
