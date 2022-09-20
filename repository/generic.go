package repository

import (
	"context"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/proto"
)

// Generic is the interface for the generic service.
type Generic interface {
	storage.Storage
	storage.Tx

	Transact(fn func(ctx context.Context, repo Generic) error)
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

// Transact is a helper function that wraps a function in a transaction and commits or rolls back the transaction. If
// svc is not a transaction, the function will be executed without executing.
func (svc *GenericService) Transact(fn func(ctx context.Context, repo Generic) error) {
	if svc.Tx != nil {
		svc.Send(func(ctx context.Context, stg storage.Storage) error {
			return fn(ctx, svc)
		})
	}
}

// Truncate truncates a table.
func (svc *GenericService) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	return svc.Storage.Truncate(ctx, req)
}
