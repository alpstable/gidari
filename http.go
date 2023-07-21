// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0

package gidari

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync"

	"github.com/alpstable/gidari/third_party/accept"
	"golang.org/x/time/rate"
)

// Request represents a request to be made by the service to the client.
// This object wraps the "net/http" package request object.
type Request struct {
	http *http.Request

	auth    func(*http.Request) (*http.Response, error) // round tripper
	writers []ListWriter
}

// RequestOption is used to set an option on a request.
type RequestOption func(*Request)

// NewHTTPRequest will create a new HTTP request.
func NewHTTPRequest(req *http.Request, opts ...RequestOption) *Request {
	hreq := &Request{http: req}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(hreq)
	}

	return hreq
}

// WithAuth will set a round tripper to be used by the service to authenticate
// the request during the http transport.
func WithAuth(auth func(*http.Request) (*http.Response, error)) RequestOption {
	return func(req *Request) {
		req.auth = auth
	}
}

// WithWriters sets optional writers to be used by the HTTP Service upsert
// method to write the data from the response.
func WithWriters(w ...ListWriter) RequestOption {
	return func(req *Request) {
		req.writers = append(req.writers, w...)
	}
}

// Client is an interface that wraps the "Do" method of the "net/http" package's
// "client" type.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPService is used process response data from requests sent to an HTTP
// client. "Processing" includes upserting data into a database, or concurrently
// iterating over the response data using a "Next" pattern.
type HTTPService struct {
	client Client
	svc    *Service

	// Iterator is a service that provides the functionality to
	// asynchronously iterate over a set of requests, handling them with a
	// custom handler. Each response in the request is achieved by calling
	// the Iterator's "Next" method, returning the "http.Response" object
	// defined by the "net/http" package.
	Iterator *HTTPIteratorService

	rlimiter *rate.Limiter
	requests []*Request
}

// NewHTTPService will create a new HTTPService.
func NewHTTPService(svc *Service) *HTTPService {
	httpSvc := &HTTPService{svc: svc, client: http.DefaultClient}
	httpSvc.Iterator = NewHTTPIteratorService(httpSvc)

	return httpSvc
}

// RateLimiter sets the optional rate limiter for the service. A rate limiter
// will limit the request to a set of bursts per period, avoiding 429 errors.
func (svc *HTTPService) RateLimiter(rlimiter *rate.Limiter) *HTTPService {
	svc.rlimiter = rlimiter

	return svc
}

// Client sets the optional client to be used by the service. If no client is
// set, the default "http.DefaultClient" defined by the "net/http" package
// will be used.
func (svc *HTTPService) Client(client Client) *HTTPService {
	svc.client = client

	return svc
}

// Requests sets the option requests to be made by the service to the client.
// If no client has been set for the service, the default "http.DefaultClient"
// defined by the "net/http" package will be used.
func (svc *HTTPService) Requests(reqs ...*Request) *HTTPService {
	svc.requests = reqs

	return svc
}

// isDecodeTypeJSON will check if the provided "accept" struct is typed for
// decoding into JSON.
func isDecodeTypeJSON(acceptHeader accept.Accept) bool {
	return acceptHeader.Typ == "application" &&
		(acceptHeader.Subtype == "json" || acceptHeader.Subtype == "*") ||
		acceptHeader.Typ == "*" && acceptHeader.Subtype == "*"
}

// bestFitDecodeType will parse the provided Accept(-Charset|-Encoding|-Language)
// header and return the header that best fits the decoding algorithm. If the
// "Accept" header is not set, then this method will return a decodeTypeJSON.
// If the "Accept" header is set, but no match is found, then this method will
// return a decodeTypeUnkown.
//
// See the "acceptSlice.Less" method in the "third_party/accept" package for
// more informaiton on how the "best fit" is determined.
func bestFitDecodeType(header string) DecodeType {
	decodeType := DecodeTypeUnknown

	for _, acceptHeader := range accept.ParseAcceptHeader(header) {
		if isDecodeTypeJSON(acceptHeader) {
			decodeType = DecodeTypeJSON

			break
		}
	}

	return decodeType
}

func (svc *HTTPService) store(ctx context.Context, jobs chan<- listWriterJob, done <-chan struct{}) error {
	for svc.Iterator.Next(ctx) {
		rsp := svc.Iterator.Current.Response

		// If there is no response, then do nothing.
		if rsp == nil {
			continue
		}

		// If response status code is not 200 (OK) return with an error
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf("%w: %d", ErrBadResponse, rsp.StatusCode)
		}

		job := &listWriterJob{writers: svc.Iterator.Current.writers}

		// Get the best fit type for decoding the response body. If the
		// best fit is "Unknown", then return an error.
		switch bestFitDecodeType(rsp.Header.Get("Accept")) {
		case DecodeTypeJSON:
			job.decFunc = decodeFuncJSON(rsp)
		case DecodeTypeUnknown:
			return fmt.Errorf("%w: %q", ErrUnsupportedDecodeType, rsp.Request.URL.String())
		}

		jobs <- *job
	}

	if err := svc.Iterator.Err(); err != nil {
		return fmt.Errorf("error iterating over requests: %w", err)
	}

	for w := 1; w <= len(svc.requests); w++ {
		<-done
	}

	// Close the jobs channel.
	close(jobs)

	if err := svc.Iterator.Close(); err != nil {
		return fmt.Errorf("failed to close iterator: %w", err)
	}

	return nil
}

// Store will concurrently make the requests to the client and store the data
// from the responses in the provided storage. If no storage is provided, then
// the data will be discarded.
func (svc *HTTPService) Store(ctx context.Context) error {
	reqCount := len(svc.requests)

	// If there are no requests, do nothing.
	if reqCount == 0 {
		return nil
	}

	// Reset the iterator.
	svc.Iterator = NewHTTPIteratorService(svc)

	listWriterCh := startListWriter(ctx, reqCount)

	upsertWorkerJobs := listWriterCh.jobs
	done := listWriterCh.done
	errCh := listWriterCh.err

	if err := svc.store(ctx, upsertWorkerJobs, done); err != nil {
		return fmt.Errorf("failed to upsert data: %w", err)
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("error in upsert worker: %w", err)
	}

	return nil
}

// Current is a struct that represents the most recent response by calling the
// "Next" method on the HTTPIteratorService.
type Current struct {
	Response *http.Response // HTTP response from the request.
	writers  []ListWriter   // Writer for storage.
}

// HTTPIteratorService is a service that will iterate over the requests defined
// for the HTTPService and return the response from each request.
type HTTPIteratorService struct {
	svc *HTTPService

	// Current is the most recent response from the iterator. This value is
	// set and blocked by the "Next" method, updating with each iteration.
	Current *Current

	currentChan chan *Current
	errCh       chan error

	// closemu prevents the iterator from closing while there is an active
	// streaming  result. It is held for read during non-close operations
	// and exclusively during close.
	//
	// closemu guards lasterr and closed.
	closemu sync.RWMutex
	closed  bool
	lasterr error
}

// NewHTTPIteratorService will return a new HTTPIteratorService.
func NewHTTPIteratorService(svc *HTTPService) *HTTPIteratorService {
	iter := &HTTPIteratorService{svc: svc, errCh: make(chan error, 1)}

	return iter
}

// Close closes the iterator.
func (iter *HTTPIteratorService) Close() error {
	iter.closemu.Lock()
	defer iter.closemu.Unlock()

	if iter.closed {
		return nil
	}

	iter.closed = true

	return nil
}

// Err returns any error encountered by the iterator.
func (iter *HTTPIteratorService) Err() error {
	iter.closemu.RLock()
	defer iter.closemu.RUnlock()

	// If the error is EOF or nil, return nil.
	if errors.Is(iter.lasterr, io.EOF) || iter.lasterr == nil {
		return nil
	}

	return iter.lasterr
}

type webWorkerJob struct {
	req      *Request
	client   Client
	rlimiter *rate.Limiter
}

type webWorkerConfig struct {
	// id is a unique identifier for the worker. This value MUST be set in
	// order to start a web worker. One and only one web worker
	// configuration MUST have an ID of 1 in order to close the response
	// channel.
	id int

	jobs      chan webWorkerJob
	currentCh chan *Current
	done      chan bool
	errCh     chan error
}

type authRoundTripper struct {
	rt func(*http.Request) (*http.Response, error)
}

// RoundTrip will execute the request and return the response.
func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return a.rt(req)
}

func fetch(ctx context.Context, job *webWorkerJob) (<-chan *http.Response, <-chan error) {
	out := make(chan *http.Response, 1)
	errs := make(chan error, 1)

	go func() {
		// If the rate limiter is not set, set it with defaults.
		if rlimiter := job.rlimiter; rlimiter != nil {
			if err := job.rlimiter.Wait(ctx); err != nil {
				errs <- fmt.Errorf("rate limiter error: %w", err)
			}
		}

		// Copy the client in case it is modified.
		client := job.client

		// If the client is an *http.Client, then set the auth
		// round-tripper.
		if client, ok := client.(*http.Client); ok && job.req.auth != nil {
			client.Transport = &authRoundTripper{rt: job.req.auth}
		}

		//nolint:bodyclose
		rsp, err := client.Do(job.req.http)
		if err != nil {
			errs <- fmt.Errorf("failed to make request: %w", err)
		}

		out <- rsp

		close(out)
		close(errs)
	}()

	return out, errs
}

// startWebWorker will start a worker upto the given specifications of the
// configuration. The worker will listen for jobs defined by the confirugation,
// asynchronous make web requests, and then propagate them onto the response
// channel.
//
// This function should be the only function that sends to the response channel
// (i.e. "rspCh"). Because this function is meant to be used as a worker pool,
// it is important that the response channel is not closed until all workers
// have finished. Therefore, this function will close the response channel ONLY
// when the worker with ID 1 has finished. This works because all workers will
// be blocked from the "defer" method until the "jobs" channel is closed.
//
// If an error is encountered, the worker will push the error onto the error
// channel and then exit. Note that only the  most recent error will be
// propagated to the "errCh" channel, per the rules of "errgroup.Group". Also,
// regardless of errors encountered, the worker will always continue to process
// jobs until the jobs channel is closed.
func startWebWorker(ctx context.Context, cfg *webWorkerConfig) {
	for job := range cfg.jobs {
		go func(job webWorkerJob) {
			defer func() {
				cfg.done <- true
			}()

			//nolint:bodyclose
			rspCh, errCh := fetch(ctx, &job)

			err := <-errCh
			if err != nil {
				cfg.errCh <- err
			}

			cfg.currentCh <- &Current{
				Response: <-rspCh,
				writers:  job.req.writers,
			}
		}(job)
	}

	if cfg.id == 1 {
		close(cfg.currentCh)
		close(cfg.done)
		close(cfg.errCh)
	}
}

// startWorkers will start the iterator's web workers and response workers. This
// method can be used to lazy load the underlying buffered channels.
func (iter *HTTPIteratorService) startWorkers(ctx context.Context) {
	reqCount := len(iter.svc.requests)
	iter.currentChan = make(chan *Current, reqCount)

	// webWorkerJobChan is responsible for making HTTP requests and pushing
	// the response body onto the responseWorkerJobChan. This channel is
	// buffered to be equal to the number of requests made.
	webWorkerJobChan := make(chan webWorkerJob, reqCount)
	done := make(chan bool, reqCount)

	// Start the web workers.
	for i := 0; i < runtime.NumCPU(); i++ {
		go startWebWorker(ctx, &webWorkerConfig{
			id:        i + 1,
			jobs:      webWorkerJobChan,
			currentCh: iter.currentChan,
			done:      done,
			errCh:     iter.errCh,
		})
	}

	go func() {
		// Send the flattened requests to the web workers for processing.
		for _, req := range iter.svc.requests {
			webWorkerJobChan <- webWorkerJob{
				req:      req,
				client:   iter.svc.client,
				rlimiter: iter.svc.rlimiter,
			}
		}
	}()

	go func() {
		// Wait for all the web workers to finish.
		for i := 0; i < reqCount; i++ {
			<-done
		}

		close(webWorkerJobChan)
	}()
}

func (iter *HTTPIteratorService) next(ctx context.Context) error {
	for {
		select {
		// If the context has timed out or been canceled, then we return
		// false.
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case result, ok := <-iter.currentChan:
			if !ok || result.Response == nil {
				// If we don't get a response, then we know
				// something is wrong and we need to wait for
				// the error channel to be closed.
				if err := <-iter.errCh; err != nil {
					return err
				}

				// Return an EOF error to indicate that we have
				// reached the end of the iterator.
				return io.EOF
			}

			iter.Current = result

			return nil
		}
	}
}

// Next will push the next response as a byte slice onto the Iterator. If there
// are no more responses, the returned boolean will be false. The user is
// responsible for decoding the response.
//
// The HTTP requests used to define the configuration will be fetched
// concurrently once the "Next" method is called for the first time.
func (iter *HTTPIteratorService) Next(ctx context.Context) bool {
	iter.closemu.RLock()
	defer iter.closemu.RUnlock()

	// If the current channel is nil, then we need to start the workers.
	// This will lazy load the web workers and the response workers, each
	// buffered by the number of requests.
	if iter.currentChan == nil {
		iter.startWorkers(ctx)
	}

	iter.lasterr = iter.next(ctx)

	return iter.lasterr == nil
}
