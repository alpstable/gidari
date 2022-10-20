// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"context"
	"encoding/json"
	"fmt"
	reflect "reflect"
	sync "sync"
	"testing"
)

var ErrTest = fmt.Errorf("test error")

// truncateStorage will truncate all tables in the storage.
func truncateStorage(ctx context.Context, t *testing.T, stg Storage) {
	t.Helper()

	if _, err := stg.Truncate(ctx, &TruncateRequest{}); err != nil {
		t.Fatalf("failed to truncate storage: %v", err)
	}
}

// truncateTables will truncate the specified tables in the storage.
func truncateTables(ctx context.Context, t *testing.T, stg Storage, tables ...string) {
	t.Helper()

	if _, err := stg.Truncate(ctx, &TruncateRequest{
		Tables: tables,
	}); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

// RunTest will run all of the test cases in the "TestRunner", and will truncate the storage in cleanup.
func RunTest(ctx context.Context, t *testing.T, stg Storage, runnerCB func(*TestRunner)) {
	t.Helper()

	runner := newTestRunner(stg)

	t.Cleanup(func() {
		// Truncate the storage.
		truncateStorage(ctx, t, stg)

		// Close the connection.
		stg.Close()
	})

	runnerCB(runner)
	runner.Run(ctx, t)
}

// TestCase is a test case for the "TestRunner".
type TestCase struct {
	Name                string                 // name is the name of the test case
	ExpectedIsNoSQL     bool                   // expectedIsNoSQL is a bool
	ExpectedUpsertSize  int64                  // expectedUpsertSize is in bits
	ExpectedPrimaryKeys map[string][]string    // expectedPrimaryKeys is a map of table name to primary keys
	Table               string                 // table is where to insert the data
	Data                map[string]interface{} // data is the data to insert
	Rollback            bool                   // rollback will rollback the transaction
	ForceError          bool                   // forceError will force an error to occur
	BinaryColumn        string                 // binaryColumn is the column to insert the binary data into
	PrimaryKeyMap       map[string]string      // primaryKeyMap is a map of data columns to primary key columns
}

// TestRunner is the storage test runner.
type TestRunner struct {
	closeDBCases         []TestCase
	isNoSQLCases         []TestCase
	listPrimaryKeysCases []TestCase
	listTablesCases      []TestCase
	upsertTxnCases       []TestCase
	upsertBinaryCases    []TestCase
	pingCases            []TestCase
	Mutex                *sync.Mutex
	Storage              Storage
}

// newTestRunner will create a new "TestRunner".
func newTestRunner(stg Storage) *TestRunner {
	return &TestRunner{
		Mutex:   &sync.Mutex{},
		Storage: stg,
	}
}

// Run will run all the "TestRunner" tests.
func (runner TestRunner) Run(ctx context.Context, t *testing.T) {
	t.Helper()

	runner.closeDB(ctx, t)
	runner.isNoSQL(ctx, t)
	runner.listTables(ctx, t)
	runner.listPrimaryKeys(ctx, t)
	runner.upsertTxn(ctx, t)
	runner.upsertBinary(ctx, t)
	runner.ping(ctx, t)
}

// AddCloseDBCases will add test cases to the "closeDB" test.
func (runner *TestRunner) AddCloseDBCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.closeDBCases = append(runner.closeDBCases, cases...)
}

func (runner *TestRunner) AddIsNoSQLCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.isNoSQLCases = append(runner.isNoSQLCases, cases...)
}

// AddListPrimaryKeysCases will add test cases to the "listPrimaryKeys" test.
func (runner *TestRunner) AddListPrimaryKeysCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.listPrimaryKeysCases = append(runner.listPrimaryKeysCases, cases...)
}

// AddListTablesCases will add test cases to the "listTables" test.
func (runner *TestRunner) AddListTablesCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.listTablesCases = append(runner.listTablesCases, cases...)
}

// AddUpsertTxnCases will add test cases to the "upsertTxn" test.
func (runner *TestRunner) AddUpsertTxnCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.upsertTxnCases = append(runner.upsertTxnCases, cases...)
}

// AddUpsertBinaryCases will add test cases to the "upsertBinary" test.
func (runner *TestRunner) AddUpsertBinaryCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.upsertBinaryCases = append(runner.upsertBinaryCases, cases...)
}

// AddPingCases will add test cases to the "ping" test.
func (runner *TestRunner) AddPingCases(cases ...TestCase) {
	runner.Mutex.Lock()
	defer runner.Mutex.Unlock()

	runner.pingCases = append(runner.pingCases, cases...)
}

// forceTxnError forces an error to occur in the transaction. It sends two requests to further test the reseliency of
// the "Send" method. Despite two requests being sent, only the first one (the one that fails) should propagate due to
// the error it returns.
func forceTxnError(t *testing.T, txn *Txn) {
	t.Helper()

	txn.Send(func(_ context.Context, _ Storage) error {
		return ErrTest
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

// closeDB will test the "Close" storage method.
func (runner TestRunner) closeDB(_ context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.closeDBCases {
		name := fmt.Sprintf("%s close db", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

		})
	}
}

// isNoSQL will test the "IsNoSQL" storage method.
func (runner TestRunner) isNoSQL(_ context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.isNoSQLCases {
		name := fmt.Sprintf("%s is no sql db", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

			got := runner.Storage.IsNoSQL()
			if got != tcase.ExpectedIsNoSQL {
				t.Fatalf("expected IsNoSQL to be: %v", tcase.ExpectedIsNoSQL)
			}
		})
	}
}

// listTables will test the "ListTables" storage method.
func (runner TestRunner) listTables(_ context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.listTablesCases {
		name := fmt.Sprintf("%s list tables", tcase.Name)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

			ctx := context.Background()

			// Upsert some data to a random table
			_, err := runner.Storage.Upsert(ctx, &UpsertRequest{
				Table: tcase.Table,
				Data:  []byte(`{"test_string": "test", "id": "1"}`),
			})
			if err != nil {
				t.Fatalf("failed to upsert data: %v", err)
			}

			// Get the table data.
			rsp, err := runner.Storage.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if len(rsp.GetTableSet()) == 0 {
				t.Fatalf("expected tables, got none")
			}

			if rsp.GetTableSet()[tcase.Table].Size == 0 {
				t.Fatalf("expected table size to be greater than zero")
			}

			truncateTables(ctx, t, runner.Storage, tcase.Table)

			// Get the table data.
			rsp, err = runner.Storage.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if rsp.GetTableSet()[tcase.Table].Size != 0 {
				t.Fatalf("expected table size to be zero")
			}
		})
	}
}

// listPrimaryKeys will test the "ListPrimaryKeys" storage method.
func (runner TestRunner) listPrimaryKeys(ctx context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.listPrimaryKeysCases {
		name := fmt.Sprintf("%s list primary keys", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

			// Insert some data to initialize NoSQL tables.
			_, err := runner.Storage.Upsert(ctx, &UpsertRequest{
				Table: tcase.Table,
				Data:  []byte(`{"test_string":"test", "test_int":1}`),
			})
			if err != nil {
				t.Fatalf("failed to upsert data: %v", err)
			}
			defer truncateTables(ctx, t, runner.Storage, tcase.Table)

			pks, err := runner.Storage.ListPrimaryKeys(ctx)
			if err != nil {
				t.Fatalf("failed to list primary keys: %v", err)
			}

			if len(pks.GetPKSet()) == 0 {
				t.Fatalf("expected primary keys, got none")
			}

			successCount := 0
			for table, pk := range pks.GetPKSet() {
				if len(tcase.ExpectedPrimaryKeys[table]) == 0 {
					continue
				}
				if reflect.DeepEqual(pk.List, tcase.ExpectedPrimaryKeys[table]) {
					successCount++
				}
			}

			if successCount != len(tcase.ExpectedPrimaryKeys) {
				t.Fatalf("expected primary keys not found")
			}
		})
	}
}

// UpsertTxn will try to simulate an upsert operation with a transaction. It will either commit or rollback the
// transaction based on the value of the "rollback" field. It can also force an error to occur in the transaction,
// depending on the testing requirements.
func (runner TestRunner) upsertTxn(ctx context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.upsertTxnCases {
		name := fmt.Sprintf("%s upsert txn", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

			if runner.Storage == nil {
				t.Fatalf("storage is nil")
			}

			txn, err := runner.Storage.StartTx(ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			if tcase.ForceError {
				forceTxnError(t, txn)

				return
			}

			// Encode some JSON data to test with.
			bytes, err := json.Marshal(tcase.Data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			txn.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &UpsertRequest{
					Table: tcase.Table,
					Data:  bytes,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert data: %w", err)
				}

				return nil
			})

			resolveTxn(t, txn, tcase.Rollback)

			// Check that the data was inserted.
			tableInfo, err := runner.Storage.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			size := tableInfo.GetTableSet()[tcase.Table].GetSize()
			if size != tcase.ExpectedUpsertSize {
				t.Fatalf("expected upsert count to be %d, got %d",
					tcase.ExpectedUpsertSize, size)
			}

			truncateTables(ctx, t, runner.Storage, tcase.Table)
		})
	}
}

// upsertBinary will test the "UpsertBinary" storage method.
func (runner TestRunner) upsertBinary(ctx context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.upsertBinaryCases {
		name := fmt.Sprintf("%s upsert binary", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runner.Mutex.Lock()
			defer runner.Mutex.Unlock()

			if runner.Storage == nil {
				t.Fatalf("storage is nil")
			}

			txn, err := runner.Storage.StartTx(ctx)
			if err != nil {
				t.Fatalf("failed to start transaction: %v", err)
			}

			if tcase.ForceError {
				forceTxnError(t, txn)

				return
			}

			// Encode some JSON data to test with.
			bytes, err := json.Marshal(tcase.Data)
			if err != nil {
				t.Fatalf("failed to marshal data: %v", err)
			}

			txn.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.UpsertBinary(sctx, &UpsertBinaryRequest{
					Table:         tcase.Table,
					Data:          bytes,
					BinaryColumn:  tcase.BinaryColumn,
					PrimaryKeyMap: tcase.PrimaryKeyMap,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert data: %w", err)
				}

				return nil
			})

			resolveTxn(t, txn, tcase.Rollback)

			// Check that the data was inserted.
			tableInfo, err := runner.Storage.ListTables(ctx)
			if err != nil {
				t.Fatalf("failed to list tables: %v", err)
			}

			if size := tableInfo.GetTableSet()[tcase.Table].GetSize(); size != tcase.ExpectedUpsertSize {
				t.Fatalf("expected upsert count to be %d, got %d", tcase.ExpectedUpsertSize, size)
			}

			truncateTables(ctx, t, runner.Storage, tcase.Table)
		})
	}
}

func (runner TestRunner) ping(_ context.Context, t *testing.T) {
	t.Helper()

	for _, tcase := range runner.pingCases {
		name := fmt.Sprintf("%s ping db", tcase.Name)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := runner.Storage.Ping(); err != nil {
				t.Errorf("An error was returned: %v", err)
			}
		})
	}
}
