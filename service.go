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

	structpb "google.golang.org/protobuf/types/known/structpb"
)

// Service is the main service for Gidari. It is responsible for providing the
// services for transporting and processing data.
type Service struct {
	// HTTP is used for transporting and processing HTTP requests and
	// responses.
	HTTP *HTTPService
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

	return svc, nil
}

// ListWriter is use to write data to io, storage, whatever, from a list of
// structpb.Values.
type ListWriter interface {
	Write(cxt context.Context, list *structpb.ListValue) error
}

type listWriterJob struct {
	table    string
	database string
	dataType DecodeType
	data     []byte
	writer   ListWriter
}

type listWriterConfig struct {
	// id is a unique identifier for the worker. This value MUST be set in
	// order to start a web worker. One and only one web worker
	// configuration MUST have an ID of 1 in order to close the response
	// channel.
	id    int
	jobs  <-chan listWriterJob
	done  chan<- struct{}
	errCh chan<- error
}

func writeList(ctx context.Context, job *listWriterJob) <-chan error {
	errs := make(chan error, 1)

	go func() {
		// Decode the data into a structpb.ListValue.
		list, err := Decode(job.dataType, job.data)
		if err != nil {
			errs <- err

			return
		}

		if err := job.writer.Write(ctx, list); err != nil {
			errs <- err
		}

		close(errs)
	}()

	return errs
}

// startListWriter will start a worker to upsert data from HTTP responses into
// a database.
func startListWriter(ctx context.Context, cfg listWriterConfig) {
	for job := range cfg.jobs {
		errs := writeList(ctx, &job)
		if err := <-errs; err != nil {
			cfg.errCh <- err
		}

		cfg.done <- struct{}{}
	}

	if cfg.id == 1 {
		close(cfg.done)
		close(cfg.errCh)
	}
}
