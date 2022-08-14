package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

// CoinbasePro are Coinbase Pro specific implementations.
type CoinbasePro interface {
	tools.GenericStorage

	UpsertJSON(context.Context, string, []byte, *proto.CreateResponse) error
	UpsertAccountsJSON(context.Context, []byte, *proto.CreateResponse) error
}

type cbp struct{ *storage }

// NewCoinbasePro will construct a repository for interacting with Coinbase Pro data.
func NewCoinbasePro(_ context.Context, r tools.GenericStorage) CoinbasePro {
	stg := new(storage)
	stg.r = newStorage(r)
	return &cbp{storage: stg}
}

// UpsertJSON will attempt to read a bytes buffer into the specified table.
func (svc *cbp) UpsertJSON(ctx context.Context, table string, b []byte, rsp *proto.CreateResponse) error {
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

// UpsertAccountsJSON will attempt to write the given bytes into the accounts table for the selected storage.
func (svc *cbp) UpsertAccountsJSON(ctx context.Context, b []byte, rsp *proto.CreateResponse) error {
	return svc.UpsertJSON(ctx, "accounts", b, rsp)
}
