package storage

import (
	"fmt"
	"sync"
	"testing"
)

func TestPGMeta(t *testing.T) {
	t.Parallel()

	mockPgMeta := pgmeta{
		cols: map[string][]string{
			"table1": []string{"0", "jason", "big ben"},
			"table2": []string{"1", "john", "bakers street"},
			"table3": []string{"2", "harry", "leicester square"},
		},
		pks: map[string][]string{
			"table1": []string{"id"},
			"table2": []string{"id", "name"},
			"table3": []string{"id", "name", "address"},
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
		t.Run(fmt.Sprintf("tableName=%s,column=%s,isPK=%v", test.tableName, test.column, test.isPk), func(t *testing.T) {
			t.Parallel()

			actual := pdb.meta.isPK(test.tableName, test.column)
			if actual != test.isPk {
				t.Fatalf("expected: %v, got: %v", test.isPk, actual)
			}
		})

	}

}
