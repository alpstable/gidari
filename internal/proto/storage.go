// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"context"
	"fmt"
	"strings"
)

const (
	// PostgresType is the byte representation of a postgres database.
	PostgresType = 0x01

	// MongoType is the byte representation of a mongo database.
	MongoType = 0x02
)

var ErrDNSNotSupported = fmt.Errorf("dns is not supported")

// DNSNotSupported wraps an error with ErrDNSNotSupported.
func DNSNotSupportedError(dns string) error {
	return fmt.Errorf("%w: %s", ErrDNSNotSupported, dns)
}

// Storage is an interface that defines the methods that a storage device should implement.
type Storage interface {
	// Close will disconnect the storage device.
	Close()

	// ListPrimaryKeys will return a list of primary keys for all tables in the database.
	ListPrimaryKeys(ctx context.Context) (*ListPrimaryKeysResponse, error)

	// ListTables will return a list of all tables in the database.
	ListTables(ctx context.Context) (*ListTablesResponse, error)

	// IsNoSQL will return true if the storage device is a NoSQL database.
	IsNoSQL() bool

	// StartTx will start a transaction and return a "Tx" object that can be used to put operations on a channel,
	// commit the result of all operations sent to the transaction, or rollback the result of all operations sent
	// to the transaction.
	StartTx(context.Context) (*Txn, error)

	// Truncate will delete all data from the storage device for ast list of tables.
	Truncate(context.Context, *TruncateRequest) (*TruncateResponse, error)

	// Type returns the type of storage device.
	Type() uint8

	// Upsert will insert or update a batch of records in the storage device.
	Upsert(context.Context, *UpsertRequest) (*UpsertResponse, error)

	// UpsertBinary will insert or update a batch of records that are part of a "property bag"-like structure that
	// containers binary data in the storage device.
	UpsertBinary(context.Context, *UpsertBinaryRequest) (*UpsertBinaryResponse, error)
}

type StorageService struct{ Storage }

// Constructor is a constructor method for a storage package.
type Constructor func(context.Context, string) (*StorageService, error)

// SchemeFromStorageType takes a byte and returns the associated DNS root database resource.
func SchemeFromStorageType(t uint8) string {
	switch t {
	case MongoType:
		return "mongodb"
	case PostgresType:
		return "postgresql"
	default:
		return "unknown"
	}
}

// SchemeFromConnectionString will return the scheme of a DNS.
func SchemeFromConnectionString(dns string) string {
	return strings.Split(dns, "://")[0]
}

// Service is a wrapper for a Storage implementation.
type Service struct {
	Storage
}
