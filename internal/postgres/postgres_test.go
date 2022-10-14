package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
)

func TestPGMeta(t *testing.T) {
	t.Parallel()

	mockPgMeta := pgmeta{
		cols: map[string][]string{
			"table1": {"0", "jason", "big ben"},
			"table2": {"1", "john", "bakers street"},
			"table3": {"2", "harry", "leicester square"},
		},
		pks: map[string][]string{
			"table1": {"id"},
			"table2": {"id", "name"},
			"table3": {"id", "name", "address"},
		},
		bytes: map[string]int64{
			"table1": 1234,
			"table2": 2345,
			"table3": 10000,
		},
	}

	pdb := &Postgres{}
	pdb.meta = &mockPgMeta
	pdb.metaMutex = sync.Mutex{}

	// testing out isPk method
	tests := []struct {
		tableName string
		column    string
		isPk      bool
	}{
		// table 1 scenarios
		{
			tableName: "table1",
			column:    "id",
			isPk:      true,
		},
		{
			tableName: "table1",
			column:    "name",
			isPk:      false,
		},
		{
			tableName: "table1",
			column:    "address",
			isPk:      false,
		},
		// table 2 scenarios
		{
			tableName: "table2",
			column:    "id",
			isPk:      true,
		},
		{
			tableName: "table2",
			column:    "name",
			isPk:      true,
		},
		{
			tableName: "table2",
			column:    "address",
			isPk:      false,
		},
		// table 3 scenarios
		{
			tableName: "table3",
			column:    "id",
			isPk:      true,
		},
		{
			tableName: "table3",
			column:    "name",
			isPk:      true,
		},
		{
			tableName: "table3",
			column:    "address",
			isPk:      true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("tableName=%s,column=%s,isPK=%v", test.tableName, test.column, test.isPk), func(t *testing.T) {
			t.Parallel()
			actual := pdb.meta.isPK(test.tableName, test.column)
			if actual != test.isPk {
				t.Fatalf("expected: %v, got: %v", test.isPk, actual)
			}
		})
	}

	// testing out upsertStmt
	ctx := context.Background()

	upsertTests := []struct {
		tableName   string
		expectedSQL string
	}{
		{
			tableName:   "table1",
			expectedSQL: `INSERT INTO table1(0,jason,big ben) VALUES ($1,$2,$3) ON CONFLICT (id) DO UPDATE SET "0" = EXCLUDED."0","jason" = EXCLUDED."jason","big ben" = EXCLUDED."big ben"`,
		},
		{
			tableName:   "table2",
			expectedSQL: `INSERT INTO table2(1,john,bakers street) VALUES ($1,$2,$3) ON CONFLICT (id,name) DO UPDATE SET "1" = EXCLUDED."1","john" = EXCLUDED."john","bakers street" = EXCLUDED."bakers street"`,
		},
		{
			tableName:   "table3",
			expectedSQL: `INSERT INTO table3(2,harry,leicester square) VALUES ($1,$2,$3) ON CONFLICT (id,name,address) DO UPDATE SET "2" = EXCLUDED."2","harry" = EXCLUDED."harry","leicester square" = EXCLUDED."leicester square"`,
		},
	}

	for _, test := range upsertTests {
		test := test
		t.Run(fmt.Sprintf("upsertStmt for %q", test.tableName), func(t *testing.T) {
			t.Parallel()

			expectedSQL := test.expectedSQL

			mockPCF := func(_ context.Context, actualQuery string) (*sql.Stmt, error) {
				if actualQuery != expectedSQL {
					return nil, fmt.Errorf("expected and actual query not same")
				}

				return &sql.Stmt{}, nil
			}

			_, err := pdb.meta.upsertStmt(ctx, test.tableName, mockPCF, 1)
			if err != nil {
				t.Fatalf("failed to create upsert statement: %v", err)
			}
		})
	}
}
