// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alpstable/gidari/proto"
	"github.com/alpstable/gidari/tools"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
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

// Config is the configuration used to query data from the web using HTTP requests and storing that data using
// the repositories defined by the "ConnectionStrings" list.
type Config struct {
	RawURL            string           `yaml:"url"`
	Authentication    Authentication   `yaml:"authentication"`
	ConnectionStrings []string         `yaml:"connectionStrings"`
	Requests          []*Request       `yaml:"requests"`
	RateLimitConfig   *RateLimitConfig `yaml:"rateLimit"`

	Logger         *logrus.Logger
	StgConstructor proto.Constructor
	Storage        []proto.Storage
	Truncate       bool

	URL *url.URL `yaml:"-"`
}

// New takes a YAML byte slice and returns a new transport configuration for upserting data to storage.
//
// For web requests defined on the transport configuration, the default HTTP Request Method is "GET". Furthermore,
// if rate limit data has not been defined for a request it will inherit the rate limit data from the transport config.
func New(_ context.Context, file *os.File) (*Config, error) {
	var cfg Config

	cfg.Logger = logrus.New()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cfg.URL, err = url.Parse(cfg.RawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %w", err)
	}

	// create a rate limiter to pass to all "flattenedRequest". This has to be defined outside of the scope of
	// individual "flattenedRequest"s so that they all share the same rate limiter, even concurrent requests to
	// different endpoints could cause a rate limit error on a web API.
	rateLimiter := rate.NewLimiter(rate.Every(*cfg.RateLimitConfig.Period), *cfg.RateLimitConfig.Burst)

	// Update default request data.
	for _, req := range cfg.Requests {
		if req.Method == "" {
			req.Method = http.MethodGet
		}

		if req.Table == "" {
			endpointParts := strings.Split(req.Endpoint, "/")
			req.Table = endpointParts[len(endpointParts)-1]
		}

		req.RateLimiter = rateLimiter
	}

	return &cfg, nil
}

// Validate will ensure that the configuration is valid for querying the web API.
func (cfg *Config) Validate() error {
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
