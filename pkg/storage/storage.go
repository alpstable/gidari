package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/sherpa/internal/storage"
	"github.com/alpine-hodler/sherpa/pkg/proto"
	"github.com/alpine-hodler/sherpa/tools"
)

// S is an implementaiton of a generic storage, any storage device should be implementable on S.
type S struct{ Generic tools.GenericStorage }

// Close will close the connection to a storage device.
func (stg *S) Close() {
	stg.Generic.Close()
}

// ExecTx will attempt to run the storage device operation within a database transaction.
func (stg *S) ExecTx(ctx context.Context, fn func(context.Context, tools.GenericStorage) (bool, error)) error {
	return stg.Generic.ExecTx(ctx, fn)
}

// Read will take data from the read request and attempt to write the server response to the "rsp" input.
func (stg *S) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	return stg.Generic.Read(ctx, req, rsp)
}

// TruncateTables will send a truncate request to the server to drop all data within a table or collection.
func (stg *S) TruncateTables(ctx context.Context, req *proto.TruncateTablesRequest) error {
	return stg.Generic.TruncateTables(ctx, req)
}

// Upsert will take the data from the upsert request and attempt to insert the data into the storage device if it is
// not present, otherwise update it. After the write has been sent to the server, the server reply is written onto
// the response "rsp" input.
func (stg *S) Upsert(ctx context.Context, req *proto.UpsertRequest, rsp *proto.CreateResponse) error {
	return stg.Generic.Upsert(ctx, req, rsp)
}

// Type will return the underlying type of the generic storage device, i.e. mongo, postgres, etc.
func (stg *S) Type() uint8 {
	return stg.Generic.Type()
}

// New will attempt to return a generic storage object given a DNS.
func New(ctx context.Context, dns string) (tools.GenericStorage, error) {
	if strings.Contains(dns, "mongodb") {
		return storage.NewMongo(ctx, dns)
	}
	if strings.Contains(dns, "postgresql") {
		return storage.NewPostgres(ctx, dns)
	}
	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}
