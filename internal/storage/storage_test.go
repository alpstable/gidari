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
	"sync"
	"testing"
	"time"

	"github.com/alpstable/gidari/proto"
	"github.com/alpstable/gidari/tools"
	"golang.org/x/sync/errgroup"
)

func truncateStorage(ctx context.Context, t *testing.T, stg Storage, tables ...string) {
	t.Helper()

	if _, err := stg.Truncate(ctx, &proto.TruncateRequest{Tables: tables}); err != nil {
		t.Fatalf("failed to truncate storage: %v", err)
	}
}

type testRunner struct {
	dns        string
	table      string
	data       map[string]interface{}
	rollback   bool
	forceError bool
	stg        Storage
}

// forceTxnError forces an error to occur in the transaction. It sends two requests to further test the reseliency of
// the "Send" method. Despite two requests being sent, only the first one (the one that fails) should propagate due to
// the error it returns.
func forceTxnError(t *testing.T, txn *Txn) {
	t.Helper()

	txn.Send(func(_ context.Context, _ Storage) error {
		return fmt.Errorf("test error")
	})

	txn.Send(func(_ context.Context, _ Storage) error {
		return nil
	})

	if err := txn.Commit(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// resolveTxn will either rollback or commit a transaction based on the value of the "rollback" field.
func resolveTxn(t *testing.T, txn *Txn, rollback bool) {
	t.Helper()

	if rollback {
		if err := txn.Rollback(); err != nil {
			t.Fatalf("failed to rollback transaction: %v", err)
		}
	} else {
		if err := txn.Commit(); err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}
	}
}

// upsertWithTx will try to simulate an upsert operation with a transaction. It will either commit or rollback the
// transaction based on the value of the "rollback" field. It can also force an error to occur in the transaction,
// depending on the testing requirements.
func (runner testRunner) upsertWithTx(ctx context.Context, t *testing.T, mtx *sync.Mutex) *proto.ListTablesResponse {
	t.Helper()

	// Lock the mutext between upsert calls
	if mtx != nil {
		mtx.Lock()
		defer mtx.Unlock()
	}

	stg := runner.stg
	if stg == nil {
		var err error
		stg, err = New(ctx, runner.dns)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		truncateStorage(ctx, t, stg, runner.table)
		t.Cleanup(func() {
			truncateStorage(ctx, t, stg, runner.table)
			stg.Close()

			t.Logf("connection closed: %q", runner.dns)
		})
	}

	txn, err := stg.StartTx(ctx)
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	if runner.forceError {
		forceTxnError(t, txn)

		return &proto.ListTablesResponse{}
	}

	// Encode some JSON data to test with.
	bytes, err := json.Marshal(runner.data)
	if err != nil {
		t.Fatalf("failed to marshal data: %v", err)
	}

	txn.Send(func(sctx context.Context, stg Storage) error {
		_, err := stg.Upsert(sctx, &proto.UpsertRequest{
			Table:    runner.table,
			Data:     bytes,
			DataType: int32(tools.UpsertDataJSON),
		})
		if err != nil {
			return fmt.Errorf("failed to upsert data: %w", err)
		}

		return nil
	})

	resolveTxn(t, txn, runner.rollback)

	// Check that the data was inserted.
	tableInfo, err := stg.ListTables(ctx)
	if err != nil {
		t.Fatalf("failed to list tables: %v", err)
	}

	return tableInfo
}

func TestTransactions(t *testing.T) {
	t.Parallel()

	// mtx for locking writes between tests
	mtx := &sync.Mutex{}

	// defaultTestTable will be the table used for all tests unless otherwise specified.
	const defaultTestTable = "tests1"

	// defaultMongoDBConnString is the default MongoDB connection string to use for tests.
	const defaultMongoDBConnString = "mongodb://mongo1:27017/defaultdb"

	// defaultPostgreSQLConnString is the default PostgreSQL connection string to use for tests.
	const defaultPostgreSQLConnString = "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"

	defaultData := map[string]interface{}{
		"test_string": "test",
		"id":          "1",
	}

	for _, tcase := range []struct {
		dns                string                 // dns is the connection string to use for the test
		name               string                 // name is the name of the test case
		expectedUpsertSize int64                  // expectedUpsertSize is in bits
		table              string                 // table is where to insert the data
		data               map[string]interface{} // data is the data to insert
		rollback           bool                   // rollback will rollback the transaction
		forceError         bool                   // forceError will force an error to occur
	}{
		{
			name:               "commit",
			dns:                defaultMongoDBConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 54,
			data:               defaultData,
		},
		{
			name:               "rollback",
			dns:                defaultMongoDBConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 0,
			rollback:           true,
			data:               defaultData,
		},
		{
			name:               "rollback on error",
			dns:                defaultMongoDBConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 0,
			forceError:         true,
			data:               defaultData,
		},
		{
			name:               "commit",
			dns:                defaultPostgreSQLConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 8192,
			data:               defaultData,
		},
		{
			name:               "rollback",
			dns:                defaultPostgreSQLConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 0,
			rollback:           true,
			data:               defaultData,
		},
		{
			name:               "rollback on error",
			dns:                defaultPostgreSQLConnString,
			table:              defaultTestTable,
			expectedUpsertSize: 0,
			forceError:         true,
			data:               defaultData,
		},
	} {
		tcase := tcase

		// Test all connection strings for each case.
		t.Run(fmt.Sprintf("%s %s", tcase.name, SchemeFromConnectionString(tcase.dns)), func(t *testing.T) {
			t.Parallel()

			tcase := tcase

			runner := testRunner{
				table:      tcase.table,
				data:       tcase.data,
				rollback:   tcase.rollback,
				forceError: tcase.forceError,
				dns:        tcase.dns,
			}

			size := runner.
				upsertWithTx(context.Background(), t, mtx).
				GetTableSet()[runner.table].
				GetSize()

			if size != tcase.expectedUpsertSize {
				t.Fatalf("expected upsert count to be %d, got %d",
					tcase.expectedUpsertSize, size)
			}
		})
	}
}

func TestConcurrentTransactions(t *testing.T) {
	t.Parallel()

	// threshold is the number of concurrent transactions to run.
	const threshold = 10

	// defaultMongoDBconnString is the default MongoDB connection string to use for tests.
	const defaultMongoDBConnString = "mongodb://mongo1:27017/defaultdb"

	// defaultPostgreSQLConnString is the default PostgreSQL connection string to use for tests.
	const defaultPostgreSQLConnString = "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"

	defaultData := map[string]interface{}{
		"test_string": "test",
		"id":          "1",
	}

	for _, tcase := range []struct {
		name               string                 // name is the name of the test case
		dns                string                 // dns is the connection string to use for the test
		data               map[string]interface{} // data is the data to insert
		expectedUpsertSize int64                  // expectedUpsertSize is in bits
		rollback           bool                   // rollback will rollback the transaction
		forceError         bool                   // forceError will force an error to occur
	}{
		{
			name:               "commit",
			dns:                defaultMongoDBConnString,
			expectedUpsertSize: 54,
			data:               defaultData,
		},
		{
			name:               "rollback",
			dns:                defaultMongoDBConnString,
			expectedUpsertSize: 0,
			rollback:           true,
			data:               defaultData,
		},
		{
			name:               "rollback on error",
			dns:                defaultMongoDBConnString,
			expectedUpsertSize: 0,
			forceError:         true,
			data:               defaultData,
		},
		{
			name:               "commit",
			dns:                defaultPostgreSQLConnString,
			expectedUpsertSize: 8192,
			data:               defaultData,
		},
		{
			name:               "rollback",
			dns:                defaultPostgreSQLConnString,
			expectedUpsertSize: 0,
			rollback:           true,
			data:               defaultData,
		},
		{
			name:               "rollback on error",
			dns:                defaultPostgreSQLConnString,
			expectedUpsertSize: 0,
			forceError:         true,
			data:               defaultData,
		},
	} {
		tcase := tcase

		// Each iteration awaits the completion of the previous iteration.
		errs, ctx := errgroup.WithContext(context.Background())

		t.Run(fmt.Sprintf("%s %s", tcase.name, SchemeFromConnectionString(tcase.dns)), func(t *testing.T) {
			fmt.Println("meep")
			t.Parallel()

			tcase := tcase

			// The storage case should be persisted throughout the subtest.
			stg, err := New(ctx, tcase.dns)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			t.Cleanup(func() {
				stg.Close()
			})

			// In theory, you should not start two clients that will commit competing transactions. The
			// real goal of this tests is to see if these resources can share the same context. So we pass
			// this lock to the test runner.
			runnerMtx := &sync.Mutex{}

			for itr := 0; itr < threshold; itr++ {
				itr := itr

				errs.Go(func() error {
					testTable := fmt.Sprintf("tests%d", itr%5+5)

					byts := []byte{0}
					if _, err := rand.Reader.Read(byts); err != nil {
						panic(err)
					}

					// Sleep for a random amount of time during the operation to help ensure that
					// the requests are occassionaly truly asynchronous.
					sleepDuration := time.Duration(byts[0]%255) * time.Millisecond
					time.Sleep(sleepDuration)

					runner := testRunner{
						stg:        stg,
						table:      testTable,
						data:       tcase.data,
						rollback:   tcase.rollback,
						forceError: tcase.forceError,
					}

					// Do not run these with a mutex as we want to test concurrent transactions.
					tableInfo := runner.upsertWithTx(ctx, t, runnerMtx)

					if tableInfo.GetTableSet()[testTable].GetSize() != tcase.expectedUpsertSize {
						return fmt.Errorf("failure to run %q for %q on table %q: "+
							"expected upsert count to be %d, got %d",
							tcase.name,
							SchemeFromConnectionString(tcase.dns),
							testTable,
							tcase.expectedUpsertSize,
							tableInfo.GetTableSet()[testTable].GetSize())
					}

					truncateStorage(ctx, t, stg, testTable)

					return nil
				})

			}
		})

		if err := errs.Wait(); err != nil {
			t.Fatal(err)
		}
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
