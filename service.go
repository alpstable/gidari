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
	"sync"

	"github.com/alpstable/gidari/proto"
)

type Service struct {
	HTTP *HTTPService
}

type ServiceOption func(*Service)

func NewService(ctx context.Context, opts ...ServiceOption) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}

	svc.HTTP = NewHTTPService(svc)

	return svc, nil
}

type upsertWorkerJob struct {
	table    string
	database string
	dataType proto.DecodeType
	data     []byte
}

type upsertWorkerConfig struct {
	// id is a unique identifier for the worker. This value MUST be set in
	// order to start a web worker. One and only one web worker
	// configuration MUST have an ID of 1 in order to close the response
	// channel.
	id      int
	jobs    <-chan upsertWorkerJob
	done    chan<- struct{}
	errCh   chan<- error
	writers []proto.UpsertWriter
}

func upsert(ctx context.Context, stg proto.UpsertWriter, job *upsertWorkerJob) <-chan error {
	errs := make(chan error, 1)

	go func() {
		req := &proto.UpsertRequest{
			Table: &proto.Table{
				Name:     job.table,
				Database: job.database,
			},
			Data: job.data,
		}

		if err := stg.Write(ctx, req); err != nil {
			errs <- err
		}

		close(errs)
	}()

	return errs
}

// startUpsertWorker will start a worker to upsert data from HTTP responses into
// a database.
func startUpsertWorker(ctx context.Context, cfg upsertWorkerConfig) {
	for job := range cfg.jobs {
		wg := sync.WaitGroup{}
		wg.Add(len(cfg.writers))

		for _, stg := range cfg.writers {
			go func(stg proto.UpsertWriter, job upsertWorkerJob) {
				defer wg.Done()

				errs := upsert(ctx, stg, &job)
				if err := <-errs; err != nil {
					cfg.errCh <- err
				}
			}(stg, job)
		}

		wg.Wait()
		cfg.done <- struct{}{}
	}

	if cfg.id == 1 {
		close(cfg.done)
		close(cfg.errCh)
	}
}
