package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
)

// storageTestCase is a test case for generic storage operations.
type storageTestCase struct {
	ctx context.Context
	dns string
}

func TestTruncate(t *testing.T) {
	testCases := []storageTestCase{
		{context.Background(), "mongodb://mongo1:27017/coinbasepro"},
		{context.Background(), "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("empty case: %s", tc.dns), func(t *testing.T) {
			s, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer s.Close()

			if _, err := s.Truncate(tc.ctx, &proto.TruncateRequest{}); err != nil {
				t.Fatalf("failed to truncate storage: %v", err)
			}
		})
		t.Run(tc.dns, func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create new storage service: %v", err)
			}
			defer stg.Close()

			rsp, err := stg.Truncate(tc.ctx, &proto.TruncateRequest{Tables: []string{"tests"}})
			if err != nil {
				t.Fatalf("failed to truncate collection: %v", err)
			}
			if rsp == nil {
				t.Fatalf("truncate response is nil")
			}

			// TODO use a read to make sure this actually worked.
		})
	}
}

func TestStartTx(t *testing.T) {
	testCases := []storageTestCase{
		{context.Background(), "mongodb://mongo1:27017/coinbasepro"},
		{context.Background(), "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("tx should commit %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tc.ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			// Encode some JSON data to test with.
			data := map[string]interface{}{"test_string": "test", "id": "1"}
			b, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			// Insert some data.
			tx.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     b,
					DataType: int32(tools.UpsertDataJSON),
				})
				return err
			})

			if err := tx.Commit(); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}

			// TODO: check if the data was actually inserted

			// Truncate the test table
			truncateReq := new(proto.TruncateRequest)
			truncateReq.Tables = []string{"tests"}
			_, err = stg.Truncate(tc.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tc.ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			// Encode some JSON data to test with.
			data := map[string]interface{}{"test_string": "test", "id": "1"}
			b, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			// Insert some data.
			tx.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     b,
					DataType: int32(tools.UpsertDataJSON),
				})
				return err
			})

			if err := tx.Rollback(); err != nil {
				t.Fatalf("failed to rollback transaction: %v", err)
			}

			// TODO: check if the data was actually inserted

			// Truncate the test table
			truncateReq := new(proto.TruncateRequest)
			truncateReq.Tables = []string{"tests"}
			_, err = stg.Truncate(tc.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", tc.dns), func(t *testing.T) {
			stg, err := New(tc.ctx, tc.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tc.ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			tx.Send(func(_ context.Context, _ Storage) error {
				return fmt.Errorf("test error")
			})

			tx.Send(func(_ context.Context, _ Storage) error {
				return nil
			})

			if err := tx.Commit(); err == nil {
				t.Fatalf("expected error, got nil")
			}

			// TODO check if the data was actually not inserted
		})

	}
}
