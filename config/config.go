// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package config

import (
	"net/url"

	"github.com/alpstable/gidari/proto"
	"github.com/alpstable/gidari/tools"
	"github.com/sirupsen/logrus"
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
