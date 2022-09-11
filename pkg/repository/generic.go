package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/pkg/proto"
	"github.com/alpine-hodler/gidari/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

// Generic is the interface for the generic service.
type Generic interface {
	storage.Storage
	storage.Tx

	UpsertRawJSON(context.Context, *Raw, *proto.UpsertResponse) error
}

// GenericService is the implementation of the Generic service.
type GenericService struct {
	storage.Storage
	storage.Tx
}

// New returns a new Generic service.
func New(ctx context.Context, dns string) (Generic, error) {
	stg, err := storage.New(ctx, dns)
	return &GenericService{stg, nil}, err
}

// NewTx returns a new Generic service with an initialized transaction object that can be used to commit or rollback
// storage operations made by the repository layer.
func NewTx(ctx context.Context, dns string) (Generic, error) {
	stg, err := storage.New(ctx, dns)
	if err != nil {
		return nil, err
	}

	tx, err := stg.StartTx(ctx)
	if err != nil {
		return nil, err
	}
	return &GenericService{stg, tx}, nil

}

// UpsertRawJSON upserts a raw json document into the database and writes the resulting document to a
// "proto.CreateResponse" object.
func (svc *GenericService) UpsertRawJSON(ctx context.Context, raw *Raw, rsp *proto.UpsertResponse) error {
	var records []*structpb.Struct
	var data interface{}
	if err := json.Unmarshal(raw.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal raw data: %w", err)
	}

	if err := tools.MakeRecordsRequest(data, &records); err != nil {
		return fmt.Errorf("error making records request: %v", err)
	}

	// If there are no records to upsert, do nothing.
	if len(records) == 0 {
		return nil
	}

	req := new(proto.UpsertRequest)
	req.Table = raw.Table
	req.Records = records
	return svc.Storage.Upsert(ctx, req, rsp)
}
