package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpine-hodler/sherpa/pkg/proto"
	"github.com/alpine-hodler/sherpa/pkg/storage"
	"github.com/alpine-hodler/sherpa/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

type Generic interface {
	storage.S

	UpsertJSON(context.Context, string, []byte, *proto.CreateResponse) error
}

type generic struct{ storage.S }

func New(_ context.Context, r storage.S) Generic {
	return &generic{r}
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
	return svc.S.Upsert(ctx, req, rsp)
}
