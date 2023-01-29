// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import "fmt"

var (
	ErrDeadlineExceeded         = fmt.Errorf("deadline exceeded")
	ErrFetchingTimeseriesChunks = fmt.Errorf("failed to fetch timeseries chunks")
	ErrInvalidRateLimit         = fmt.Errorf("invalid rate limit configuration")
	ErrMissingConfigField       = fmt.Errorf("missing config field")
	ErrMissingRateLimitField    = fmt.Errorf("missing rate limit field")
	ErrMissingTimeseriesField   = fmt.Errorf("missing timeseries field")
	ErrMissingURL               = fmt.Errorf("missing URL")
	ErrSettingTimeseriesChunks  = fmt.Errorf("failed to set timeseries chunks")
	ErrUnableToParse            = fmt.Errorf("unable to parse")
	ErrNilConfig                = fmt.Errorf("config is nil")
	ErrNoRequests               = fmt.Errorf("no requests defined")
	ErrInvalidEndTimeSize       = fmt.Errorf("invalid end time size, expected 1")
	ErrInvalidStartTimeSize     = fmt.Errorf("invalid start time size, expected 1")
)

func wrapErrDeadlineExceeded(err error) error {
	return fmt.Errorf("%w: %v", ErrDeadlineExceeded, err)
}

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
