package storage

import (
	"context"
	"fmt"

	"github.com/alpine-hodler/sherpa/pkg/proto"
)

const (
	MongoType uint8 = iota
	PostgresType
)

type Storage interface {
	Close()
	// ExecTx(context.Context, func(context.Context, Storage) (bool, error)) error
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error
	TruncateTables(context.Context, *proto.TruncateTablesRequest) error
	Upsert(context.Context, *proto.UpsertRequest, *proto.CreateResponse) error
}

// DNSRoot takes a byte and returns the associated DNS root database resource.
func DNSRoot(t uint8) (string, error) {
	switch t {
	case MongoType:
		return "mongodb", nil
	case PostgresType:
		return "postgresql", nil
	default:
		return "", fmt.Errorf("type %q is not supported", t)
	}
}
