package tools

import (
	"context"

	"github.com/alpine-hodler/gidari/pkg/proto"
)

// GenericStorage provides CRUD methods for interacting with an arbitrary DB.
type GenericStorage interface {
	Close()
	ExecTx(context.Context, func(context.Context, GenericStorage) (bool, error)) error
	Read(context.Context, *proto.ReadRequest, *proto.ReadResponse) error
	TruncateTables(context.Context, *proto.TruncateTablesRequest) error
	Upsert(context.Context, *proto.UpsertRequest, *proto.CreateResponse) error
	Type() uint8
}
