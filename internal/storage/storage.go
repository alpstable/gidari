package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/sherpa/pkg/proto"
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
	// ExecTx(context.Context, func(context.Context, Storage) (bool, error)) error
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error
	TruncateTables(context.Context, *proto.TruncateTablesRequest) error
	Upsert(context.Context, *proto.UpsertRequest, *proto.CreateResponse) error
	Type() uint8
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
