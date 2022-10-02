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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
	"golang.org/x/sync/errgroup"
)

func truncateStorage(ctx context.Context, t *testing.T, stg Storage, tables ...string) {
	t.Helper()

	if _, err := stg.Truncate(ctx, &proto.TruncateRequest{Tables: tables}); err != nil {
		t.Fatalf("failed to truncate storage: %v", err)
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct{ dns string }{
		{"mongodb://mongo1:27017/db1"},
		{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
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

	for _, tcase := range []struct{ dns string }{
		{"mongodb://mongo1:27017/db2"},
		{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
		dns := tcase.dns
		ctx := context.Background()

		stg, err := New(ctx, dns)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		const (
			testTable2 = "tests2"
			testTable3 = "tests3"
			testTable4 = "tests4"
		)

		truncateStorage(ctx, t, stg, testTable2, testTable3, testTable4)
		t.Cleanup(func() {
			truncateStorage(ctx, t, stg, testTable2, testTable3, testTable4)
			stg.Close()

			t.Logf("connection closed: %q", dns)
		})

		t.Run(fmt.Sprintf("tx should commit %s", dns), func(t *testing.T) {
			t.Parallel()

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
					Table:    testTable2,
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

			if tableInfo.GetTableSet()[testTable2].GetSize() == 0 {
				t.Fatalf("expected data to be inserted for %q", testTable2)
			}
		})
		t.Run(fmt.Sprintf("tx should rollback %s", dns), func(t *testing.T) {
			t.Parallel()

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
					Table:    testTable3,
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

			// Check that the data was inserted.
			tableInfo, err := stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if tableInfo.GetTableSet()[testTable3].GetSize() != 0 {
				t.Fatalf("expected data to be rolled back")
			}
		})
		t.Run(fmt.Sprintf("tx should rollback on error %s", dns), func(t *testing.T) {
			t.Parallel()

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

			if tableInfo.GetTableSet()[testTable4].GetSize() != 0 {
				t.Fatalf("expected data to be rolled back")
			}
		})
		t.Run("tx should commit in parallel with the same context", func(t *testing.T) {
			t.Parallel()

			// Run this test 3 times.
			for j := 0; j < 3; j++ {
				errs, ctx := errgroup.WithContext(context.Background())
				for i := 0; i < 10; i++ {
					errs.Go(func() error {
						// Get a random number between 5 and 10.
						byts := []byte{0}
						if _, err := rand.Reader.Read(byts); err != nil {
							panic(err)
						}
						randNum := byts[0]%5 + 5

						testTable := fmt.Sprintf("tests%d", randNum)

						// Wait for a random amount of milliseconds.

						sleepDuration := time.Duration(byts[0]) * time.Millisecond
						time.Sleep(sleepDuration)

						txn, err := stg.StartTx(ctx)
						if err != nil {
							return fmt.Errorf("failed to start transaction: %w", err)
						}

						// Encode some JSON data to test with.
						data := map[string]interface{}{"test_string": "test", "id": "1"}
						bytes, err := json.Marshal(data)
						if err != nil {
							return fmt.Errorf("failed to marshal data: %w", err)
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
							return fmt.Errorf("failed to commit transaction: %w", err)
						}

						// Check that the data was inserted.
						tableInfo, err := stg.ListTables(ctx)
						if err != nil {
							return fmt.Errorf("failed to list tables: %w", err)
						}

						if tableInfo.GetTableSet()[testTable].GetSize() == 0 {
							return fmt.Errorf("expected data to be inserted for %q", testTable2)
						}

						return nil
					})
				}

				if err := errs.Wait(); err != nil {
					t.Fatalf("failed to commit transaction: %v", err)
				}
			}
		})
	}
}

func TestListTables(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct{ dns string }{
		{"mongodb://mongo1:27017/db4"},
		{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
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

			truncateStorage(ctx, t, stg, testTable)

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

	for _, tcase := range []struct {
		dns           string
		expectedPKSet map[string][]string
	}{
		{
			"mongodb://mongo1:27017/db3",
			map[string][]string{
				"tests5": {"_id"},
			},
		},
		{
			"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable",
			map[string][]string{
				"tests5":         {"id"},
				"candle_minutes": {"product_id", "unix"},
			},
		},
	} {
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

			// Insert some data to initialize NoSQL tables.
			_, err = stg.Upsert(ctx, &proto.UpsertRequest{
				Table:    "tests5",
				Data:     []byte(`{"id": 1, "test_string":"test"}`),
				DataType: int32(tools.UpsertDataJSON),
			})
			if err != nil {
				t.Fatalf("failed to upsert data: %v", err)
			}

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
