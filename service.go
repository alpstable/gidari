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
	"sync"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

// Service is the main service for Gidari. It is responsible for providing the
// services for transporting and processing data.
type Service struct {
	// HTTP is used for transporting and processing HTTP requests and
	// responses.
	HTTP *HTTPService

	// Socket is used for transporting and processing data over a socket
	// connection.
	Socket *SocketService
}

// ServiceOption is a function for configuring a Service.
type ServiceOption func(*Service)

// NewService will create a new Service.
func NewService(ctx context.Context, opts ...ServiceOption) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}

	svc.HTTP = NewHTTPService(svc)
	svc.Socket = NewSocketService(svc)

	return svc, nil
}

// ListWriter is use to write data to io, storage, whatever, from a list of
// structpb.Values.
type ListWriter interface {
	Write(cxt context.Context, list *structpb.ListValue) error
}

type listWriterJob struct {
	decFunc DecodeFunc
	writers []ListWriter
}

func writeList(ctx context.Context, job *listWriterJob) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		list := &structpb.ListValue{}
		if err := job.decFunc(list); err != nil {
			errs <- err

			return
		}

		wg := &sync.WaitGroup{}
		wg.Add(len(job.writers))

		for _, writer := range job.writers {
			writer := writer

			go func(writer ListWriter) {
				defer wg.Done()

				if err := writer.Write(ctx, list); err != nil {
					errs <- err
				}
			}(writer)
		}

		wg.Wait()
	}()

	return errs
}

// upsertWorkerJobs := make(chan listWriterJob)
// go startListWriter(ctx, listWriterConfig{
//	id:    1,
//	jobs:  upsertWorkerJobs,
//	errCh: errs,
// })

type listWriterChan struct {
	done <-chan struct{}
	err  <-chan error
	jobs chan<- listWriterJob
}

// startListWriter will start a worker to upsert data from HTTP responses into
// a database.
func startListWriter(ctx context.Context, numJobs int) listWriterChan {
	if numJobs == 0 {
		numJobs = 1
	}

	var (
		jobs  chan listWriterJob
		done  chan struct{}
		errCh chan error
	)

	if numJobs > 0 {
		done = make(chan struct{}, numJobs)
		jobs = make(chan listWriterJob, numJobs)
	} else {
		jobs = make(chan listWriterJob)
	}

	errCh = make(chan error, 1)

	go func() {
		defer close(errCh)

		if done != nil {
			defer close(done)
		}

		for job := range jobs {
			errs := writeList(ctx, &job)
			if err := <-errs; err != nil {
				errCh <- err
			}

			if done != nil {
				done <- struct{}{}
			}
		}
	}()

	return listWriterChan{
		done: done,
		err:  errCh,
		jobs: jobs,
	}
}
