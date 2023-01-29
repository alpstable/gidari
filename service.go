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

	"github.com/alpstable/gidari/internal/repository"
	"github.com/alpstable/gidari/proto"
	"golang.org/x/sync/errgroup"
)

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

type Service struct {
	storage []*Storage

	HTTP *HTTPService
}

type ServiceOption func(*Service)

// WithStorage will add storage to the service. The storage will be used to
// store data from the source API. For instance, if a client uses the HTTP
// service, then the data will be "transported" from a web API request to
// the storage devices defined by this method. There is no default storage
// device and if no storage is defined, then the data will be requested but
// disgarded.
func WithStorage(stgs ...*Storage) ServiceOption {
	return func(svc *Service) {
		svc.storage = stgs
	}
}

func NewService(ctx context.Context, opts ...ServiceOption) (*Service, error) {
	svc := &Service{}
	for _, opt := range opts {
		opt(svc)
	}

	svc.HTTP = NewHTTPService(svc)

	return svc, nil
}

type repo struct {
	repository.Generic

	// closable is a flag that indicates whether or not the repository is
	// closable. A repository is only closable if it is created by a
	// connection string. If a repository is created by a client or a
	// database, then it is the responsibility of the caller to close.
	closable bool
	database string
}

func (gen *repo) Close() {
	if gen.closable {
		gen.Generic.Close()
	}
}

func (svc *Service) repositories(ctx context.Context) ([]repo, error) {
	repositories := make([]repo, len(svc.storage))
	for idx, stg := range svc.storage {
		genRepository, err := repository.NewTx(ctx, stg.Storage)
		if err != nil {
			return nil, fmt.Errorf("failed to create repository: %w", err)
		}

		repositories[idx] = repo{
			Generic:  genRepository,
			database: stg.Database,
			closable: stg.Close,
		}
	}

	return repositories, nil
}

type upsertWorkerJob struct {
	table string
	data  []byte
}

type upsertWorkerConfig struct {
	jobs         <-chan upsertWorkerJob
	done         chan<- struct{}
	errs         *errgroup.Group
	repositories []repo
}

func upsert_(_ context.Context, req *proto.UpsertRequest, repo repo) error {
	txfn := func(sctx context.Context, repo repository.Generic) error {
		_, err := repo.Upsert(sctx, req)
		if err != nil {
			return fmt.Errorf("error upserting data: %w", err)
		}

		return nil
	}

	repo.Transact(txfn)

	return nil
}

// startUpsertWorker will start a worker to upsert data from HTTP responses into
// a database.
func startUpsertWorker(ctx context.Context, cfg upsertWorkerConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-cfg.jobs:
			table := job.table

			errs, ctx := errgroup.WithContext(ctx)
			for _, repo := range cfg.repositories {
				repo := repo

				errs.Go(func() error {
					req := &proto.UpsertRequest{
						Table: &proto.Table{
							Name:     table,
							Database: repo.database,
						},
						Data: job.data,
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
