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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/internal/web"
	"github.com/alpine-hodler/gidari/internal/web/auth"
	"github.com/alpine-hodler/gidari/repository"
	"github.com/alpine-hodler/gidari/tools"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

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

	cfg.Logger = logrus.New()

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

	if cfg.ConnectionStrings == nil {
		logWarn := tools.LogFormatter{
			Msg: "no connectionStrings specified in the config file",
		}
		cfg.Logger.Warn(logWarn.String())
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
