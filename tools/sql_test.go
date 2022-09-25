package tools

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestSqlIterativePlaceholders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		numCols  int
		numRows  int
		symbol   string
		expected string
	}{
		{
			name:     "simple case",
			numCols:  3,
			numRows:  2,
			symbol:   "$",
			expected: "($1,$2,$3),($4,$5,$6)",
		},
		{
			name:     "no placeholders",
			numCols:  0,
			numRows:  0,
			symbol:   "",
			expected: "()",
		},
		{
			name:     "one placeholder",
			numCols:  1,
			numRows:  1,
			symbol:   "",
			expected: "(?1)",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := SQLIterativePlaceholders(test.numCols, test.numRows, test.symbol)
			if actual != test.expected {
				t.Errorf("SqlIterativePlaceholders(%d, %d, %q) = %q; want %q", test.numCols,
					test.numRows, test.symbol, actual, test.expected)
			}
		})
	}
}

func TestSqlFlattenPartition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		columns   []string
		partition []*structpb.Struct
		expected  []interface{}
	}{
		{
			name:    "simple case",
			columns: []string{"a", "b", "c"},
			partition: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"a": {Kind: &structpb.Value_StringValue{StringValue: "1"}},
						"b": {Kind: &structpb.Value_StringValue{StringValue: "2"}},
						"c": {Kind: &structpb.Value_StringValue{StringValue: "3"}},
					},
				},
			},
			expected: []interface{}{"1", "2", "3"},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := SQLFlattenPartition(test.columns, test.partition)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("SqlFlattenPartition(%q, %q) = %q; want %q", test.columns, test.partition,
					actual, test.expected)
			}
		})
	}
}
