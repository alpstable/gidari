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
	"net/http"
	"runtime"

	"github.com/alpstable/gidari/internal/web"
	"golang.org/x/sync/errgroup"
)

type Handler interface {
	ResponseCallback(ctx context.Context, rsp *http.Response, cfg *flattenedRequestConfig) error
}

type Current struct {
	Response *http.Response
	Table    string
}

// Iterator holds the request state of a gidari configuration and can be used to iterate over the request results.
type Iterator struct {
	// Current is a byte slice of the most recent request pushed onto the iterator by the "Next" method.
	Current *Current

	//webResultAssigner WebResultAssigner

	// cfg is the configuration for the iterator.
	//cfg *Config

	// currentChan is a channel that holds the end-result of the response worker. The iterator's "Next" method is
	// used to push data from the currentChan onto the Current field.
	//
	// The size of the currentChan is partially non-deterministic. That is, the buffer size should be equal to the
	// number of results in an HTTP JSON response. However, the number of results is not known until the response
	// is received and decoded by the response worker. Therefore, this channel must remain unbuffered and the
	// closure of this channel is left to the response worker.
	currentChan chan *Current

	// errCh is a channel that holds any errors encountered by the iterator. It can be propagated to the user
	// by the "Err" method.
	errCh chan error

	// err is the error that was encountered by the iterator and is propagated to the user with the "Err" method.
	err error

	// flattenRequests is a slice of flattened requests that will be used to make HTTP requests. The number of
	// web requests made will be equal to the length of this slice.
	flattenedRequests []*flattenedRequest

	handler Handler
}

// NewIterator returns an Iterator object for the given configuration.
func NewIterator(ctx context.Context, cfg *Config) (*Iterator, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	iter := &Iterator{
		errCh: make(chan error, 1),
		//handler: handler,
		//webResultAssigner: assignWebResult,
	}

	var err error

	iter.flattenedRequests, err = flattenConfigRequests(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return iter, nil
}

// sanitizeJSON will sanitize invalid JSON by convering the source bytes into a string and using the clob column
// on the request to create a valid JSON object. This function will return "false" if there is an error creating the
// new JSON object or if the source is invalid JSON and no clob column is defined for the request.
func sanitizeJSON(src []byte, clobColumn string) ([]byte, bool, error) {
	// If the source is valid JSON, then we can return the source bytes.
	if json.Valid(src) {
		return src, true, nil
	}

	//// If the source is invalid JSON and there is no clob column defined, then we cannot sanitize the JSON.
	//if clobColumn == "" {
	//	return nil, false, nil
	//}

	// If the source is invalid JSON and there is a clob column defined, then we can create a new JSON object
	// using the clob column.
	//obj := map[string]string{
	//	clobColumn: string(src),
	//}

	//fmt.Println("obj", obj)

	return nil, false, nil

	//newSrc, err := json.Marshal(obj)
	//if err != nil {
	//	return nil, false, fmt.Errorf("error marshaling new JSON object: %w", err)
	//}

	//return newSrc, true, nil
}

// assignWebResult is the default response job function that is used to assign HTTP response data to
// proto.IteratorResult objects.
//func assignWebResult(ctx context.Context, webResult WebResult) ([]*proto.IteratorResult, error) {
//	// If the response body is nil, then we can return an empty slice of results.
//	if webResult.Body == nil {
//		return []*proto.IteratorResult{}, nil
//	}
//
//	// read the bytes from the response body.
//	src, err := io.ReadAll(webResult.Response.Body)
//	if err != nil {
//		return nil, fmt.Errorf("error reading response body: %w", err)
//	}
//
//	// sanitize the bytes if they are invalid JSON.
//	//src, ok, err := sanitizeJSON(src, webResult.ClobColumn)
//	//if err != nil {
//	//	return nil, fmt.Errorf("error sanitizing invalid JSON: %w", err)
//	//}
//
//	//// if the bytes are invalid JSON and we cannot sanitize them, then we cannot process the response.
//	//if !ok {
//	//	return nil, nil
//	//}
//
//	return proto.DecodeIteratorResults(webResult.URL.String(), src)
//}

// Close will close the iterator and release any resources.
func (iter *Iterator) Close(ctx context.Context) {}

// Err returns any error encountered by the iterator.
func (iter *Iterator) Err() error {
	return iter.err
}

type responseWorkerConfig struct {
	coreNum int

	// jobs are a channel of HTTP Response objects sent by the web worker.
	jobs <-chan *http.Response

	// responseChan
	responseChan chan<- *http.Response

	done  chan<- bool
	errCh chan<- error
}

// responseWorkerResults will process the response worker job and push the results onto the current channel.
//func responseWorkerResults(ctx context.Context, cfg responseWorkerConfig, rsp WebResult) error {
//	defer func() {
//		cfg.done <- true
//	}()
//
//	if rsp.Body == nil {
//		return nil
//	}
//
//	results, err := cfg.assigner(context.Background(), rsp)
//	if err != nil {
//		return err
//	}
//
//	for _, result := range results {
//		cfg.currentIteratorResultJobChan <- result
//	}
//
//	return nil
//}

// startResponseWorker ...
// The goal of this worker is to process the HTTP responses either by sending the response to the response channel for
// the user to process, or by calling the response handler function to process the response.
//func startResponseWorker(_ context.Context, cfg responseWorkerConfig) {
//	defer func() {
//		if cfg.coreNum == 1 {
//			close(cfg.currentIteratorResultJobChan)
//			close(cfg.done)
//		}
//	}()
//
//	for job := range cfg.jobs {
//		if err := responseWorkerResults(context.Background(), cfg, job); err != nil {
//			cfg.errCh <- err
//		}
//	}
//}

type webWorkerJob struct {
	req *flattenedRequest
}

type webWorkerConfig struct {
	// id is a unique identifier for the worker. This value MUST be set in order to start a web worker. One and
	// only one web worker configuration MUST have an ID of 1 in order to close the response channel.
	id int

	jobs      chan webWorkerJob
	currentCh chan *Current
	done      chan bool
	errCh     chan error

	hasClosed  bool
	jobCounter int

	handler Handler
}

// startWebWorker will start a worker that will make a request to the given URL and push the response onto the
// the repository channel. If the request fails, the error will be pushed onto the error channel which will
// propagate to the iterator's Next method.
//
// This function should be the only function that sends to the response channel (i.e. "rspCh"). Because this function
// is meant to be used as a worker pool, it is important that the response channel is not closed until all workers have
// finished. Therefore, this function will close the response channel ONLY when the worker with ID 1 has finished. This
// works because all workers will be blocked from the "defer" method until the "jobs" channel is closed.
//
// If an error is encountered, the worker will push the error onto the error channel and then exit. Note that only the
// most recent error will be propagated to the "errCh" channel, per the rules of "errgroup.Group". Also, regardless of
// errors encountered, the worker will always continue to process jobs until the jobs channel is closed.
func startWebWorker(ctx context.Context, cfg *webWorkerConfig) {
	errs, ctx := errgroup.WithContext(ctx)

	for job := range cfg.jobs {
		job := job

		errs.Go(func() error {
			select {
			case <-ctx.Done():
				cfg.jobCounter++

				return ctx.Err()
			default:
				//cfg.jobCounter++
				var httpRsp *http.Response
				defer func() {
					var table string
					if cfg := job.req.cfg; cfg != nil {
						table = cfg.table
					}

					// Push the response onto the response channel, if it exists.
					cfg.currentCh <- &Current{
						Response: httpRsp,
						Table:    table,
					}

					// Alert the done channel that the job is complete, if it exists.
					cfg.done <- true
				}()

				rsp, err := web.Fetch(ctx, job.req.fetchConfig)
				if err != nil {
					return err
				}

				httpRsp = rsp

				//// Check if there is a web handler set on the configuration. If there is, then send the
				//// response to the web handler.
				//if cfg.handler != nil {
				//	err := cfg.handler.ResponseCallback(ctx, httpRsp, job.req.cfg)
				//	if err != nil {
				//		return err
				//	}
				//}

				return nil
			}
		})
	}

	if err := errs.Wait(); err != nil {
		cfg.errCh <- err
	}

	if cfg.id == 1 {
		close(cfg.currentCh)
		close(cfg.done)
		close(cfg.errCh)
	}
}

// startWorkers will start the iterator's web workers and response workers. This method can be used to lazy load the
// underlying buffered channels.
func (iter *Iterator) startWorkers(ctx context.Context) {
	reqCount := len(iter.flattenedRequests)
	iter.currentChan = make(chan *Current, reqCount)

	// webWorkerJobChan is responsible for making HTTP requests and pushing the response body onto the
	// responseWorkerJobChan. This channel is buffered to be equal to the number of requests made.
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
			handler:   iter.handler,
		})
	}

	go func() {
		// Send the flattened requests to the web workers for processing.
		for _, req := range iter.flattenedRequests {
			webWorkerJobChan <- webWorkerJob{req: req}
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

// Next will push the next response as a byte slice onto the Iterator. If there are no more responses, the
// returned boolean will be false. The user is responsible for decoding the response.
//
// The HTTP requests used to define the configuration will be fetched concurrently once the "Next" method is called for
// the first time.
func (iter *Iterator) Next(ctx context.Context) bool {
	// If the current channel is nil, then we need to start the workers. This will lazy load the web workers and
	// the response workers, each buffered by the number of flattened requests.
	if iter.currentChan == nil {
		iter.startWorkers(ctx)
	}

	for {
		select {
		// If the context has timed out or been canceled, then we return false.
		case <-ctx.Done():
			return false
		case err := <-iter.errCh:
			if err != nil {
				iter.err = err

				return false
			}
		case result, ok := <-iter.currentChan:
			if !ok {
				return false
			}

			// If the result is "nil," then something is wrong with the response. We will skip this response
			// and continue to encounter the underlying error.
			if result == nil {
				continue
			}

			iter.Current = result

			return true
		}
	}
}
