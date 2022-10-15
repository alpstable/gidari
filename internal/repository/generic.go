// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package repository

import (
	"context"
	"fmt"

	"github.com/alpstable/gidari/internal/mongo"
	"github.com/alpstable/gidari/internal/postgres"
	"github.com/alpstable/gidari/internal/proto"
)

// ErrFailedToCreateRepository is returned when the repository layer fails to create a new repository.
var ErrFailedToCreateRepository = fmt.Errorf("failed to create repository")

// FailedToCreateRepositoryError is a helper function that returns a new error with the ErrFailedToCreateRepository
// error wrapped.
func FailedToCreateRepositoryError(err error) error {
	return fmt.Errorf("%w: %v", ErrFailedToCreateRepository, err)
}

// Generic is the interface for the generic service.
type Generic interface {
	proto.Storage
	proto.Transactor

	Transact(fn func(ctx context.Context, repo Generic) error)
}

// GenericService is the implementation of the Generic service.
type GenericService struct {
	proto.Storage
	*proto.Txn
}

// NewStorage returns a new storage service.
func NewStorage(ctx context.Context, dns string) (*proto.StorageService, error) {
	var stg *proto.StorageService

	scheme := proto.SchemeFromConnectionString(dns)
	switch scheme {
	case proto.SchemeFromStorageType(proto.MongoType):
		mdb, err := mongo.New(ctx, dns)
		if err != nil {
			return nil, err
		}

		stg = &proto.StorageService{Storage: mdb}
	case proto.SchemeFromStorageType(proto.PostgresType):
		pdb, err := postgres.New(ctx, dns)
		if err != nil {
			return nil, err
		}

		stg = &proto.StorageService{Storage: pdb}
	default:
		return nil, fmt.Errorf("unknown scheme: %s", scheme)
	}

	return stg, nil
}

// New returns a new Generic service.
func New(ctx context.Context, dns string) (*GenericService, error) {
	stg, err := NewStorage(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("failed to construct storage: %w", err)
	}

	return &GenericService{stg, nil}, nil
}

// NewTx returns a new Generic service with an initialized transaction object that can be used to commit or rollback
// storage operations made by the repository layer.
func NewTx(ctx context.Context, dns string) (*GenericService, error) {
	stg, err := NewStorage(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("failed to construct storage: %w", err)
	}

	tx, err := stg.StartTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return &GenericService{stg, tx}, nil
}

// Transact is a helper function that wraps a function in a transaction and commits or rolls back the transaction. If
// svc is not a transaction, the function will be executed without executing.
func (svc *GenericService) Transact(fn func(ctx context.Context, repo Generic) error) {
	svc.Txn.Send(func(ctx context.Context, stg proto.Storage) error {
		err := fn(ctx, svc)
		if err != nil {
			return fmt.Errorf("error executing transaction: %w", err)
		}

		return nil
	})
}

// Truncate truncates a table.
func (svc *GenericService) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	rsp, err := svc.Storage.Truncate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error truncating table: %w", err)
	}

	return rsp, nil
}
