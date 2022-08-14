package storage

import (
	"context"

	"github.com/alpine-hodler/driver/proto"
)

type Storage interface {
	Close()
	ExecTx(context.Context, func(context.Context, Storage) (bool, error)) error
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error
	TruncateTables(context.Context, *proto.TruncateTablesRequest) error
	Upsert(context.Context, *proto.UpsertRequest, *proto.CreateResponse) error
}
