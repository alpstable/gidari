package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alpine-hodler/gidari/proto"
)

const (
	// MongoType is the byte representation of a mongo database.
	MongoType uint8 = iota

	// PostgresType is the byte representation of a postgres database.
	PostgresType
)

var (
	ErrDNSNotSupported     = fmt.Errorf("dns is not supported")
	ErrTransactionNotFound = fmt.Errorf("transaction not found")
	ErrNoTables            = fmt.Errorf("no tables found")
	ErrTransactionAborted  = fmt.Errorf("transaction aborted")
)

// DNSNotSupported wraps an error with ErrDNSNotSupported.
func DNSNotSupportedError(dns string) error {
	return fmt.Errorf("%w: %s", ErrDNSNotSupported, dns)
}

// Storage is an interface that defines the methods that a storage device should implement.
type Storage interface {
	// Close will disconnect the storage device.
	Close()

	// ListPrimaryKeys will return a list of primary keys for all tables in the database.
	ListPrimaryKeys(ctx context.Context) (*proto.ListPrimaryKeysResponse, error)

	// ListTables will return a list of all tables in the database.
	ListTables(ctx context.Context) (*proto.ListTablesResponse, error)

	// IsNoSQL will return true if the storage device is a NoSQL database.
	IsNoSQL() bool

	// StartTx will start a transaction and return a "Tx" object that can be used to put operations on a channel,
	// commit the result of all operations sent to the transaction, or rollback the result of all operations sent
	// to the transaction.
	StartTx(context.Context) (*Txn, error)

	// Truncate will delete all data from the storage device for ast list of tables.
	Truncate(context.Context, *proto.TruncateRequest) (*proto.TruncateResponse, error)

	// Type returns the type of storage device.
	Type() uint8

	// Upsert will insert or update a batch of records in the storage device.
	Upsert(context.Context, *proto.UpsertRequest) (*proto.UpsertResponse, error)
}

// sqlPrepareContextFn can be used to prepare a statement and return the result.
type sqlPrepareContextFn func(context.Context, string) (*sql.Stmt, error)

// Scheme takes a byte and returns the associated DNS root database resource.
func Scheme(t uint8) string {
	switch t {
	case MongoType:
		return "mongodb"
	case PostgresType:
		return "postgresql"
	default:
		return "unknown"
	}
}

// Service is a wrapper for a Storage implementation.
type Service struct {
	Storage
}

// New will attempt to return a generic storage object given a DNS.
func New(ctx context.Context, dns string) (*Service, error) {
	if strings.Contains(dns, Scheme(MongoType)) {
		svc, err := NewMongo(ctx, dns)
		if err != nil {
			return nil, fmt.Errorf("failed to construct mongo storage: %w", err)
		}

		return &Service{svc}, nil
	}

	if strings.Contains(dns, Scheme(PostgresType)) {
		svc, err := NewPostgres(ctx, dns)
		if err != nil {
			return nil, fmt.Errorf("failed to construct postgres storage: %w", err)
		}

		return &Service{svc}, nil
	}

	return nil, DNSNotSupportedError(dns)
}
