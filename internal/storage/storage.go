package storage

import (
	"context"
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

// Storage is an interface that defines the methods that a storage device should implement.
type Storage interface {
	Close()
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error

	// StartTx will start a transaction and return a "Tx" object that can be used to put operations on a channel,
	// commit the result of all operations sent to the transaction, or rollback the result of all operations sent
	// to the transaction.
	StartTx(context.Context) (Tx, error)
	Truncate(context.Context, *proto.TruncateRequest) (*proto.TruncateResponse, error)
	Upsert(context.Context, *proto.UpsertRequest) (*proto.UpsertResponse, error)
	Type() uint8
}

// Tx is an interface that defines the methods that a transaction object should implement.
type Tx interface {
	Commit() error
	Rollback() error
	Send(TXChanFn)
}

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

// New will attempt to return a generic storage object given a DNS.
func New(ctx context.Context, dns string) (Storage, error) {
	if strings.Contains(dns, Scheme(MongoType)) {
		return NewMongo(ctx, dns)
	}

	if strings.Contains(dns, Scheme(PostgresType)) {
		return NewPostgres(ctx, dns)
	}

	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}
