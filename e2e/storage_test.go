// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0\n
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/alpstable/gidari/internal/proto"
	"github.com/alpstable/gidari/internal/repository"
)

func truncateStorage(ctx context.Context, t *testing.T, stg proto.Storage, tables ...string) {
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
	stg        proto.Storage
}

// forceTxnError forces an error to occur in the transaction. It sends two requests to further test the reseliency of
// the "Send" method. Despite two requests being sent, only the first one (the one that fails) should propagate due to
// the error it returns.
func forceTxnError(t *testing.T, txn *proto.Txn) {
	t.Helper()

	txn.Send(func(_ context.Context, _ proto.Storage) error {
		return fmt.Errorf("test error")
	})

	txn.Send(func(_ context.Context, _ proto.Storage) error {
		return nil
	})

	if err := txn.Commit(); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// resolveTxn will either rollback or commit a transaction based on the value of the "rollback" field.
func resolveTxn(t *testing.T, txn *proto.Txn, rollback bool) {
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

func newStg(ctx context.Context, t *testing.T, dns string) *proto.StorageService {
	t.Helper()

	stg, err := repository.NewStorage(ctx, dns)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	return stg
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
		stg = newStg(ctx, t, runner.dns)

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

	txn.Send(func(sctx context.Context, stg proto.Storage) error {
		_, err := stg.Upsert(sctx, &proto.UpsertRequest{
			Table: runner.table,
			Data:  bytes,
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

//nolint:tparallel
func TestTransactions(t *testing.T) {
	t.Parallel()

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

	// Running these tests in parallel will inevitably lead to race conditions.
	//
	//nolint:paralleltest
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
		//{
		//	name:               "commit",
		//	dns:                defaultPostgreSQLConnString,
		//	table:              defaultTestTable,
		//	expectedUpsertSize: 8192,
		//	data:               defaultData,
		//},
		//{
		//	name:               "rollback",
		//	dns:                defaultPostgreSQLConnString,
		//	table:              defaultTestTable,
		//	expectedUpsertSize: 0,
		//	rollback:           true,
		//	data:               defaultData,
		//},
		//{
		//	name:               "rollback on error",
		//	dns:                defaultPostgreSQLConnString,
		//	table:              defaultTestTable,
		//	expectedUpsertSize: 0,
		//	forceError:         true,
		//	data:               defaultData,
		//},
	} {
		// Test all connection strings for each case.
		t.Run(fmt.Sprintf("%s %s", tcase.name, proto.SchemeFromConnectionString(tcase.dns)), func(t *testing.T) {
			tcase := tcase

			runner := testRunner{
				table:      tcase.table,
				data:       tcase.data,
				rollback:   tcase.rollback,
				forceError: tcase.forceError,
				dns:        tcase.dns,
			}

			size := runner.
				upsertWithTx(context.Background(), t, nil).
				GetTableSet()[runner.table].
				GetSize()

			if size != tcase.expectedUpsertSize {
				t.Fatalf("expected upsert count to be %d, got %d",
					tcase.expectedUpsertSize, size)
			}
		})
	}
}

func TestListTables(t *testing.T) {
	t.Parallel()

	// defaultTestTable is the default test table to use for tests.
	const defaultTestTable = "lttests1"

	for _, tcase := range []struct{ dns string }{
		{"mongodb://mongo1:27017/db4"},
		//{"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"},
	} {
		dns := tcase.dns
		t.Run(fmt.Sprintf("get size: %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			stg := newStg(ctx, t, dns)
			defer stg.Close()

			truncateStorage(ctx, t, stg)

			// Upsert some data to a random table
			_, err := stg.Upsert(ctx, &proto.UpsertRequest{
				Table: defaultTestTable,
				Data:  []byte(`{"test_string": "test", "id": "1"}`),
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

			if rsp.GetTableSet()[defaultTestTable].Size == 0 {
				t.Fatalf("expected table size to be greater than zero")
			}

			truncateStorage(ctx, t, stg, defaultTestTable)

			// Get the table data.
			rsp, err = stg.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if rsp.GetTableSet()[defaultTestTable].Size != 0 {
				t.Fatalf("expected table size to be zero")
			}
		})
	}
}

func TestListPrimaryKeys(t *testing.T) {
	t.Parallel()

	// defaultTestTable is the default table name used for testing.
	const defaultTestTable = "pktests1"

	for _, tcase := range []struct {
		dns           string
		expectedPKSet map[string][]string
	}{
		{
			"mongodb://mongo1:27017/db3",
			map[string][]string{
				defaultTestTable: {"_id"},
			},
		},
		//{
		//	"postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable",
		//	map[string][]string{
		//		defaultTestTable: {"test_string"},
		//	},
		//},
	} {
		dns := tcase.dns
		expectedPKSet := tcase.expectedPKSet

		t.Run(fmt.Sprintf("list tables %s", dns), func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			stg := newStg(ctx, t, dns)
			t.Cleanup(func() {
				stg.Close()
			})

			// Insert some data to initialize NoSQL tables.
			_, err := stg.Upsert(ctx, &proto.UpsertRequest{
				Table: defaultTestTable,
				Data:  []byte(`{"test_string":"test", "test_int":1}`),
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
