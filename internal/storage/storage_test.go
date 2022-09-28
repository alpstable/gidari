package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
)

// storageTestCase is a test case for generic storage operations.
type storageTestCase struct {
	dns string
}

type listPKStorageTestCase struct {
	storageTestCase

	// Map a table to a list of primary keys we expect to find on that table.
	expectedPKSet map[string][]string
}

var testCases = []storageTestCase{
	{"mongodb://mongo1:27017/cbp-stg"},
	{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
}

var listPKTestCases = []listPKStorageTestCase{
	{
		storageTestCase{"mongodb://mongo1:27017/cbp-stg"},
		map[string][]string{
			"accounts": {"_id"},
		},
	},
	{
		storageTestCase{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
		map[string][]string{
			"accounts":       {"id"},
			"candle_minutes": {"product_id", "unix"},
		},
	},
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	for _, tcase := range testCases {
		dns := tcase.dns
		t.Run(fmt.Sprintf("empty case: %s", dns), func(t *testing.T) {
			ctx := context.Background()
			t.Parallel()
			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}
			defer stg.Close()

			if _, err := stg.Truncate(ctx, &proto.TruncateRequest{}); err != nil {
				t.Fatalf("failed to truncate storage: %v", err)
			}
		})
		t.Run(dns, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create new storage service: %v", err)
			}
			defer stg.Close()

			rsp, err := stg.Truncate(ctx, &proto.TruncateRequest{Tables: []string{"tests"}})
			if err != nil {
				t.Fatalf("failed to truncate collection: %v", err)
			}
			if rsp == nil {
				t.Fatalf("truncate response is nil")
			}
		})
	}
}

func TestStartTx(t *testing.T) {
	t.Parallel()

	for _, tcase := range testCases {
		dns := tcase.dns
		t.Run(fmt.Sprintf("tx should commit %s", dns), func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			txn, err := stg.StartTx(ctx)
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
			txn.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     bytes,
					DataType: int32(tools.UpsertDataJSON),
				})
				if err != nil {
					return fmt.Errorf("failed to upsert data: %w", err)
				}

				return nil
			})

			if err := txn.Commit(); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}

			// Truncate the test table
			truncateReq := new(proto.TruncateRequest)
			truncateReq.Tables = []string{"tests"}
			_, err = stg.Truncate(ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", dns), func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			txn, err := stg.StartTx(ctx)
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
			txn.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    "tests",
					Data:     dataBytes,
					DataType: int32(tools.UpsertDataJSON),
				})
				if err != nil {
					return fmt.Errorf("failed to insert data: %w", err)
				}

				return nil
			})

			if err := txn.Rollback(); err != nil {
				t.Fatalf("failed to rollback transaction: %v", err)
			}

			// Truncate the test table
			truncateReq := new(proto.TruncateRequest)
			truncateReq.Tables = []string{"tests"}
			_, err = stg.Truncate(ctx, truncateReq)
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", dns), func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			txn, err := stg.StartTx(ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			txn.Send(func(_ context.Context, _ Storage) error {
				return fmt.Errorf("test error")
			})

			txn.Send(func(_ context.Context, _ Storage) error {
				return nil
			})

			if err := txn.Commit(); err == nil {
				t.Fatalf("expected error, got nil")
			}
		})
	}
}

func TestListTables(t *testing.T) {
	t.Parallel()

	for _, tcase := range testCases {
		dns := tcase.dns
		t.Run(fmt.Sprintf("list tables %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			// If the type is mongodb, then we need to seed some collections.
			if stg.Type() == MongoType {
				// Upsert some data to a random table
				_, err := stg.Upsert(ctx, &proto.UpsertRequest{
					Table:    "accounts",
					Data:     []byte(`{"test": "test"}`),
					DataType: int32(tools.UpsertDataJSON),
				})
				if err != nil {
					t.Fatalf("failed to upsert data: %v", err)
				}
			}

			rsp, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if len(rsp.GetTableSet()) == 0 {
				t.Fatalf("expected tables, got none")
			}

			// Make sure that one of the records has the table "accounts".
			for table := range rsp.GetTableSet() {
				if table == "accounts" {
					return
				}
			}

			t.Fatalf("expected to find table accounts, got none")
		})
	}
}

func TestListPrimaryKeys(t *testing.T) {
	t.Parallel()

	for _, tcase := range listPKTestCases {
		dns := tcase.dns
		expectedPKSet := tcase.expectedPKSet

		t.Run(fmt.Sprintf("list tables %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			pks, err := stg.ListPrimaryKeys(ctx)
			if err != nil {
				t.Fatalf("failed to list primary keys: %v", err)
			}

			if len(pks.GetPKSet()) == 0 {
				t.Fatalf("expected primary keys, got none")
			}

			successCount := 0
			for table, pk := range pks.GetPKSet() {
				if len(expectedPKSet[table]) == 0 {
					continue
				}
				if reflect.DeepEqual(pk.List, expectedPKSet[table]) {
					successCount++
				}
			}

			if successCount != len(expectedPKSet) {
				t.Fatalf("expected primary keys not found")
			}
		})
	}
}
