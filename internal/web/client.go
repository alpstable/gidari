// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/alpine-hodler/gidari/internal/web/auth"
	"golang.org/x/time/rate"
)

var (
	// ErrCreatingRequest is returned when the request fails to create.
	ErrCreatingRequest = errors.New("failed to create request")

	// ErrFailedToCreateWebClient is returned when the client fails to create.
	ErrFailedToCreateWebClient = errors.New("failed to create web client")

	// ErrGettingResponse is returned when the response fails to get.
	ErrGettingResponse = errors.New("failed to get response")

	// ErrInvalidResponse is returned when the response is invalid.
	ErrInvalidResponse = errors.New("invalid response")

	// ErrMissingFetchConfigField is returned when a required field is missing.
	ErrMissingFetchConfigField = errors.New("missing required field on FetchConfig")
)

// CreateRequestError is returned when the request fails to create.
func CreateRequestError(err error) error {
	return fmt.Errorf("%w: %v", ErrCreatingRequest, err)
}

// FailedToCreateClientError is returned when the client fails to create.
func FailedToCreateClientError(err error) error {
	return fmt.Errorf("%w: %v", ErrFailedToCreateWebClient, err)
}

// MissingFetchConfigFieldError is returned when a required field is missing.
func MissingFetchConfigFieldError(field string) error {
	return fmt.Errorf("%w: %q", ErrMissingFetchConfigField, field)
}

// GettingResponseError is returned when the response fails to get.
func GettingResponseError(rsp *http.Response) error {
	if _, err := io.ReadAll(rsp.Body); err != nil {
		return fmt.Errorf("%w: %v", ErrGettingResponse, err)
	}

	return fmt.Errorf("%w: %v", ErrGettingResponse, rsp.Status)
}

// Client is a wrapper around the http.Client that will handle authentication and rate limiting.
type Client struct{ http.Client }

// NewClient will return a new client with the given options.
func NewClient(_ context.Context, roundtripper auth.Transport) (*Client, error) {
	c := new(Client)
	c.Client.Transport = roundtripper

	return c, nil
}

// newHTTPRequest will return a new request.  If the options are set, this function will encode a body if possible.
func newHTTPRequest(ctx context.Context, method string, uri fmt.Stringer) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, uri.String(), nil)
	if err != nil {
		return nil, CreateRequestError(err)
	}

	return req, nil
}

// validateResponse is a switch condition that parses an error response.
func validateResponse(res *http.Response) error {
	if res == nil {
		return ErrInvalidResponse
	}

	switch res.StatusCode {
	case
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusInternalServerError,
		http.StatusNotFound,
		http.StatusTooManyRequests,
		http.StatusForbidden:
		return GettingResponseError(res)
	}

	return nil
}

type FetchConfig struct {
	C           *Client
	Method      string
	URL         *url.URL
	RateLimiter *rate.Limiter
}

func (cfg *FetchConfig) validate() error {
	if cfg.C == nil {
		return MissingFetchConfigFieldError("Client")
	}

	if cfg.Method == "" {
		return MissingFetchConfigFieldError("Method")
	}

	if cfg.URL == nil {
		return MissingFetchConfigFieldError("URL")
	}

	if cfg.RateLimiter == nil {
		return MissingFetchConfigFieldError("RateLimiter")
	}

	return nil
}

// FetchResponse is a wrapper for the Fetch function's response data for an HTTP web request.
type FetchResponse struct {
	// Request is the request that was made to the server.
	Request *http.Request

	// Body is the response body from the server.
	Body io.ReadCloser
}

func newFetchResponse(req *http.Request, body io.ReadCloser) *FetchResponse {
	return &FetchResponse{
		Request: req,
		Body:    body,
	}
}

// Fetch will make an HTTP request using the underlying client and endpoint.
func Fetch(ctx context.Context, cfg *FetchConfig) (*FetchResponse, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// If the rate limiter is not set, set it with defaults.
	if err := cfg.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	fmt.Println("how are you not waiting?")

	req, err := newHTTPRequest(ctx, cfg.Method, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("rate limiter timeout: %w", err)
	}

	rsp, err := cfg.C.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if err := validateResponse(rsp); err != nil {
		rsp.Body.Close()

		return nil, fmt.Errorf("error validating response: %w", err)
	}

	return newFetchResponse(req, rsp.Body), nil
}
