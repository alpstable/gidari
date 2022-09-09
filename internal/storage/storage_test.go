package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/alpine-hodler/gidari/pkg/proto"
	"github.com/alpine-hodler/gidari/tools"
)

// storageTestCase is a test case for generic storage operations.
type storageTestCase struct {
	ctx context.Context
	dns string
}

func TestStartTx(t *testing.T) {
	testCases := []storageTestCase{
		{context.Background(), "mongodb://mongo-coinbasepro:27017/coinbasepro"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("tx should commit %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			tx := stg.StartTx(tc.ctx)

			req := new(proto.UpsertRequest)
			req.Table = "test"

			data := map[string]interface{}{"test": "test"}
			err = tools.MakeRecordsRequest(data, &req.Records)
			if err != nil {
				t.Fatalf("failed to make records request: %v", err)
			}

			var rsp proto.CreateResponse
			tx.Transact(func(sctx context.Context) error {
				return stg.Upsert(sctx, req, &rsp)
			})

			if err := tx.Commit(); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}

			// TODO: check if the data was actually inserted

			// Truncate the test table
			truncateReq := new(proto.TruncateTablesRequest)
			truncateReq.Tables = []string{"test"}
			err = stg.TruncateTables(tc.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			tx := stg.StartTx(tc.ctx)

			req := new(proto.UpsertRequest)
			req.Table = "rollback_test"

			data := map[string]interface{}{"test": "test"}
			err = tools.MakeRecordsRequest(data, &req.Records)
			if err != nil {
				t.Fatalf("failed to make records request: %v", err)
			}

			var rsp proto.CreateResponse
			tx.Transact(func(sctx context.Context) error {
				return stg.Upsert(sctx, req, &rsp)
			})

			if err := tx.Rollback(); err != nil {
				t.Fatalf("failed to rollback transaction: %v", err)
			}

			// TODO: check if the data was actually inserted

			// Truncate the test table
			truncateReq := new(proto.TruncateTablesRequest)
			truncateReq.Tables = []string{"rollback_test"}
			err = stg.TruncateTables(tc.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			tx := stg.StartTx(tc.ctx)

			req := new(proto.UpsertRequest)
			req.Table = "rollback_err_test"

			tx.Transact(func(sctx context.Context) error {
				return fmt.Errorf("test error")
			})

			tx.Transact(func(sctx context.Context) error {
				return nil
			})

			if err := tx.Commit(); err == nil {
				t.Fatalf("expected error, got nil")
			}

			// TODO check if the data was actually not inserted
		})

	}
}
