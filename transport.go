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
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/alpstable/gidari/internal/repository"
	"github.com/alpstable/gidari/internal/web"
	"github.com/alpstable/gidari/proto"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

// Transport will construct the transport operation using a "transport.Config"
// object.
//
// Configuration request and response handlers will be overwriten by the
// Transport function, they are used as listeners for when an HTTP request has
// been made and when a response has been received.
func Transport(ctx context.Context, cfg *Config) error {
	if cfg == nil {
		return ErrNilConfig
	}

	// upsertHandler is responsible for sending the HTTP responses to the
	// upserter.
	//upsertHandler, err := newUpsertHandler(ctx, cfg)
	//if err != nil {
	//	return err
	//}

	// Create an iterator that will iterate over the requests in the
	// configuration file.
	iter, err := NewIterator(ctx, cfg)
	if err != nil {
		return fmt.Errorf("unable to create iterator: %w", err)
	}

	defer iter.Close(ctx)

	reqCount := len(iter.flattenedRequests)

	// Create a channel to send requests to the worker.
	currentCh := make(chan *Current, reqCount)

	// done is a channel that will be closed when the worker is done.
	done := make(chan struct{}, reqCount)

	// Create repositories
	repositories, err := newGenericRepositories(ctx, cfg)
	if err != nil {
		return err
	}

	// Start the upsert worker.
	for i := 1; i <= runtime.NumCPU(); i++ {
		go startUpsertWorker(ctx, upsertWorkerConfig{
			jobs:         currentCh,
			repositories: repositories,
			done:         done,
		})
	}

	for iter.Next(ctx) {
		currentCh <- iter.Current
	}

	// TODO: if the iterator experiences an error, we need to make sure that the upser worker does not actually
	// TODO: persist any data. Upserts should be done as a transaction, so this is definitely possible.
	if err := iter.Err(); err != nil {
		return fmt.Errorf("error iterating over requests: %w", err)
	}

	for w := 1; w <= reqCount; w++ {
		<-done
	}

	return nil
}

type genericRepository struct {
	repository.Generic

	// closable is a flag that indicates whether or not the repository is closable. A repository is only closable
	// if it is created by a connection string. If a repository is created by a client or a database, then it is
	// the responsibility of the caller to close.
	closable bool
	database string
}

func (gen *genericRepository) Close() {
	if gen.closable {
		gen.Generic.Close()
	}
}

func newGenericRepositories(ctx context.Context, cfg *Config) ([]genericRepository, error) {
	genericRepositories := make([]genericRepository, len(cfg.Storage))

	for idx, stg := range cfg.Storage {
		genRepository, err := repository.NewTx(ctx, stg.Storage)
		if err != nil {
			return nil, fmt.Errorf("failed to create repository: %w", err)
		}

		genericRepositories[idx] = genericRepository{
			Generic:  genRepository,
			database: stg.Database,
			closable: stg.Close,
		}
	}

	return genericRepositories, nil
}

type upsertWorkerConfig struct {
	jobs         <-chan *Current
	done         chan<- struct{}
	errs         *errgroup.Group
	repositories []genericRepository
}

func upsert_(_ context.Context, req *proto.UpsertRequest, genericRepository genericRepository) error {
	txfn := func(sctx context.Context, repo repository.Generic) error {
		_, err := repo.Upsert(sctx, req)
		if err != nil {
			return fmt.Errorf("error upserting data: %w", err)
		}

		return nil
	}

	genericRepository.Transact(txfn)

	return nil
}

// startUpsertWorker will start a worker to upsert data from HTTP responses into a database.
func startUpsertWorker(ctx context.Context, cfg upsertWorkerConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-cfg.jobs:
			table := job.Table
			rsp := job.Response

			// If the response is nil, then we can skip this job as
			// there is no data to upsert.
			if rsp == nil {
				continue
			}

			// Get bytes from response body.
			body, _ := ioutil.ReadAll(rsp.Body)

			defer rsp.Body.Close()

			errs, ctx := errgroup.WithContext(ctx)
			for _, repo := range cfg.repositories {
				repo := repo

				errs.Go(func() error {
					req := &proto.UpsertRequest{
						Table: &proto.Table{
							Name:     table,
							Database: repo.database,
						},
						Data: body,
					}

					err := upsert_(ctx, req, repo)
					if err != nil {
						return err
					}

					return err
				})

			}

			if err := errs.Wait(); err != nil {
				panic(err)
			}

			cfg.done <- struct{}{}
		}
	}
}

// newFetchConfig will construct a new HTTP request from the transport request.

func newFetchConfig(req *http.Request, client web.Client, rateLimiter *rate.Limiter) *web.FetchConfig {
	///rurl.Path = path.Join(rurl.Path, req.Endpoint)

	///// Add the query params to the URL.
	///if req.Query != nil {
	///	query := rurl.Query()
	///	for key, value := range req.Query {
	///		query.Set(key, value)
	///	}

	///	rurl.RawQuery = query.Encode()
	///}

	return &web.FetchConfig{
		Client:      client,
		RateLimiter: rateLimiter,
		Request:     req,
	}
}
