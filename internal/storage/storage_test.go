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
	t.Parallel()
	for _, tcase := range []storageTestCase{
		{context.Background(), "mongodb://mongo1:27017/coinbasepro"},
		{context.Background(), "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
		t.Run(fmt.Sprintf("empty case: %s", tcase.dns), func(t *testing.T) {
			t.Parallel()
			s, err := New(tcase.ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer s.Close()

			if _, err := s.Truncate(tcase.ctx, &proto.TruncateRequest{}); err != nil {
				t.Fatalf("failed to truncate storage: %v", err)
			}
		})
		t.Run(tcase.dns, func(t *testing.T) {
			t.Parallel()
			stg, err := New(tcase.ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create new storage service: %v", err)
			}
			defer stg.Close()

			rsp, err := stg.Truncate(tcase.ctx, &proto.TruncateRequest{Tables: []string{"tests"}})
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
	t.Parallel()
	for _, tcase := range []storageTestCase{
		{context.Background(), "mongodb://mongo1:27017/coinbasepro"},
		{context.Background(), "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
		t.Run(fmt.Sprintf("tx should commit %s", tcase.dns), func(t *testing.T) {
			t.Parallel()
			stg, err := New(tcase.ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tcase.ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			// Encode some JSON data to test with.
			data := map[string]interface{}{"test_string": "test", "id": "1"}
			bytes, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			// Insert some data.
			tx.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     bytes,
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
			_, err = stg.Truncate(tcase.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", tcase.dns), func(t *testing.T) {
			t.Parallel()
			stg, err := New(tcase.ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tcase.ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			// Encode some JSON data to test with.
			data := map[string]interface{}{"test_string": "test", "id": "1"}
			dataBytes, err := json.Marshal(data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			// Insert some data.
			tx.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     dataBytes,
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
			_, err = stg.Truncate(tcase.ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", tcase.dns), func(t *testing.T) {
			t.Parallel()
			stg, err := New(tcase.ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			tx, err := stg.StartTx(tcase.ctx)
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
