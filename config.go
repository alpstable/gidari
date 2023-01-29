// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"net/http"
	"time"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/time/rate"
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

type Storage struct {
	Storage proto.Storage

	// Database is the name of the database to run operations against. This is an optional field and will not be
	// needed for every storage device. It is needed for storage like MongoDB, for instance, which needs a client
	// to make transactions. But it is not needed by PostgreSQL or a file system.
	Database string

	// Close indicates that the storage should be closed after the transport operation. It is not recommended to
	// set this to true unless you are running a single transport operation. The primary use case for this is
	// with CLI commands that create connections to a database given a connection string.
	Close bool `yaml:"-"`

	// Connection string is the URI to connect to the database. This is only valid using the CLI.
	ConnectionString *string `yaml:"connectionString"`
}

// RateLimitConfig is the data needed for constructing a rate limit for the HTTP requests.
type RateLimitConfig struct {
	// Burst represents the number of requests that we limit over a period frequency.
	Burst int `yaml:"burst"`

	// Period is the time between each burst.
	Period time.Duration `yaml:"period"`
}

func (rl RateLimitConfig) validate() error {
	if rl.Burst == 0 {
		return MissingRateLimitFieldError("burst")
	}

	if rl.Period == 0 {
		return MissingRateLimitFieldError("period")
	}

	return nil
}

// Timeseries is a struct that contains the information needed to query a web API for Timeseries data.
type Timeseries struct {
	StartName string `yaml:"startName"`
	EndName   string `yaml:"endName"`

	// Period is the size of each chunk in seconds for which we can query the API. Some API will not allow us to
	// query all data within the start and end range.
	Period int32 `yaml:"period"`

	// Layout is the time layout for parsing the "Start" and "End" values into "time.Time". The default is assumed
	// to be RFC3339.
	Layout string `yaml:"layout"`

	// Chunks are the time ranges for which we can query the API. These are broken up into pieces for API requests
	// that only return a limited number of results.
	Chunks [][2]time.Time
}

// Request is the information needed to query the web API for data to transport.
//type Request struct {
//*http.Request

//HttpResponseHandler IteratorWebResponseHandler

//HttpRequestHandler IteratorWebRequestHandler

// Timeseries indicates that the underlying data should be queries as a time series. This means that the
//Timeseries *Timeseries `yaml:"timeseries"`

// Table is the name of the table/collection to insert the data fetched from the web API.
//Table string `yaml:"table"`

// Truncate before upserting on single request
//Truncate *bool `yaml:"truncate"`

//ClobColumn string `yaml:"clobColumn"`

// Chunks of requests should share a rate limiter, probably all of them; inheriting the rate limiter from the
// root configuration.
// RateLimiter *rate.Limiter
//}

//func (req *Request) validate() error {
//	return nil
//}

// WebResult is a wrapper for the HTTP response body returned by fetching on an endpoint defined by a Request. It
// also holds other data that is requird for constructing a "proto.IteratorResult" slice for end-user consumption.
//type WebResult struct {
//	*http.Response
//
//	ClobColumn string
//	TableName  string
//	URL        *url.URL
//}

type Request struct {
	*http.Request

	Table string `yaml:"table"`

	Storage []*Storage `yaml:"storage"`

	RateLimiter *rate.Limiter `yaml:"-"`
}

type Config struct {
	//RawURL          string           `yaml:"url"`
	// Authentication  Authentication   `yaml:"authentication"`

	Requests []*Request `yaml:"requests"`

	Client Client `yaml:"-"`

	//RateLimiter *rate.Limiter `yaml:"-"`

	// RateLimitConfig *RateLimitConfig `yaml:"rateLimit"`
	//Storage []*Storage `yaml:"storage"`

	// Truncate bool `yaml:"-"`

	// URL is a required field and is used to construct the HTTP request to fetch data from the web API.
	//URL *url.URL `yaml:"-"`

	// Client is the HTTP client used to run the requests defined on the configuraiton. This is an optional field
	// and will default to http.DefaultClient if not set.
	//Client *http.Client `yaml:"-"`

	// HTTPResponseHAndler is an optional function that is used in the iterator process to transform the HTTP
	// response body into a slice of proto.IteratorResult objects. Note that this function cannot be set on a
	// Transport, it will be overwritten by the Transport's upserter's HTTPResponseHandler.
	//HTTPResponseHandler WebResultAssigner `yaml:"-"`
}

// validate will ensure that the configuration is valid for querying the web API.
func (cfg *Config) validate() error {
	//if cfg.URL == nil {
	//	return ErrMissingURL
	//}

	//for _, req := range cfg.Requests {
	//	if err := req.validate(); err != nil {
	//		return err
	//	}
	//}

	return nil
}
