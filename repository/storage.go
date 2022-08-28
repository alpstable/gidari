package repository

import (
	"context"

	"github.com/alpine-hodler/sherpa/proto"
	"github.com/alpine-hodler/sherpa/tools"
)

type storage struct{ r tools.GenericStorage }

func newStorage(r tools.GenericStorage) tools.GenericStorage {
	return &storage{r}
}

func (stg *storage) Close() {
	stg.r.Close()
}

func (stg *storage) ExecTx(ctx context.Context, fn func(context.Context, tools.GenericStorage) (bool, error)) error {
	return stg.r.ExecTx(ctx, fn)
}

func (stg *storage) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	return stg.r.Read(ctx, req, rsp)
}

func (stg *storage) TruncateTables(ctx context.Context, req *proto.TruncateTablesRequest) error {
	return stg.r.TruncateTables(ctx, req)
}

func (stg *storage) Upsert(ctx context.Context, req *proto.UpsertRequest, rsp *proto.CreateResponse) error {
	return stg.r.Upsert(ctx, req, rsp)
}

func (stg *storage) Type() uint8 {
	return stg.r.Type()
}
