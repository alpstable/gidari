package storage

import (
	"context"

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
}

// DNSRoot takes a byte and returns the associated DNS root database resource.
func DNSRoot(t uint8) string {
	switch t {
	case MongoType:
		return "mongodb"
	case PostgresType:
		return "postgresql"
	default:
		return "unknown"
	}
}
