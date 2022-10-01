// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

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

func truncateStorage(ctx context.Context, t *testing.T, stg Storage, tables ...string) {
	t.Helper()

	if _, err := stg.Truncate(ctx, &proto.TruncateRequest{Tables: tables}); err != nil {
		t.Fatalf("failed to truncate storage: %v", err)
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	for _, tcase := range testCases {
		dns := tcase.dns
		t.Run(fmt.Sprintf("empty case: %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			testTable := "tests1"

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create storage: %v", err)
			}

			truncateStorage(ctx, t, stg, testTable)
			t.Cleanup(func() {
				truncateStorage(ctx, t, stg)
				stg.Close()
			})
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
			testTable := "tests2"

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			truncateStorage(ctx, t, stg, testTable)
			t.Cleanup(func() {
				truncateStorage(ctx, t, stg, testTable)
				stg.Close()
			})

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
					Table:    testTable,
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

			// Check that the data was inserted.
			tableInfo, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if tableInfo.GetTableSet()[testTable].GetSize() == 0 {
				t.Fatalf("expected data to be inserted")
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			testTable := "tests3"

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			truncateStorage(ctx, t, stg, testTable)
			t.Cleanup(func() {
				truncateStorage(ctx, t, stg, testTable)
				stg.Close()
			})

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
					Table:    testTable,
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

			time.Sleep(1 * time.Second)

			// Check that the data was inserted.
			tableInfo, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if tableInfo.GetTableSet()[testTable].GetSize() != 0 {
				t.Fatalf("expected data to be rolled back")
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			testTable := "tests4"

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			truncateStorage(ctx, t, stg, testTable)
			t.Cleanup(func() {
				truncateStorage(ctx, t, stg, testTable)
				stg.Close()
			})

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

			// Check that the data was inserted.
			tableInfo, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if tableInfo.GetTableSet()[testTable].GetSize() != 0 {
				t.Fatalf("expected data to be rolled back")
			}
		})
	}
}

func TestListTables(t *testing.T) {
	t.Parallel()

	for _, tcase := range testCases {
		dns := tcase.dns
		t.Run(fmt.Sprintf("get accounts size: %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			testTable := "accounts"

			stg, err := New(ctx, dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			defer stg.Close()

			truncateStorage(ctx, t, stg)

			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}

			// Upsert some data to a random table
			_, err = stg.Upsert(ctx, &proto.UpsertRequest{
				Table: "accounts",
				Data: []byte(`{
"id": "1",
"available": 1,
"balance": 1,
"hold": 0,
"currency": "A",
"profile_id": "1",
"trading_enabled": true
}`),
				DataType: int32(tools.UpsertDataJSON),
			})
			if err != nil {
				t.Fatalf("failed to upsert data: %v", err)
			}

			// Get the table data.
			rsp, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if len(rsp.GetTableSet()) == 0 {
				t.Fatalf("expected tables, got none")
			}

			if rsp.GetTableSet()[testTable].Size == 0 {
				t.Fatalf("expected table size to be greater than zero")
			}

			// Truncate the test table.
			_, err = stg.Truncate(ctx, &proto.TruncateRequest{
				Tables: []string{testTable},
			})
			if err != nil {
				t.Fatalf("failed to truncate table: %v", err)
			}

			// Get the table data.
			rsp, err = stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if rsp.GetTableSet()[testTable].Size != 0 {
				t.Fatalf("expected table size to be zero")
			}
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

			t.Cleanup(func() {
				stg.Close()
			})

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
