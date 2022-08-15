package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpine-hodler/driver/proto"
	"github.com/alpine-hodler/driver/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

type Generic interface {
	tools.GenericStorage

	UpsertJSON(context.Context, string, []byte, *proto.CreateResponse) error
}

type generic struct{ *storage }

func New(_ context.Context, r tools.GenericStorage) Generic {
	stg := new(storage)
	stg.r = newStorage(r)
	return &generic{storage: stg}
}

// UpsertJSON will attempt to read a bytes buffer into the specified table.
func (svc *generic) UpsertJSON(ctx context.Context, table string, b []byte, rsp *proto.CreateResponse) error {
	var records []*structpb.Struct
	var data interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if err := tools.MakeRecordsRequest(data, &records); err != nil {
		return fmt.Errorf("error making records request: %v", err)
	}
	req := new(proto.UpsertRequest)
	req.Table = table
	req.Records = records
	return svc.r.Upsert(ctx, req, rsp)
}
