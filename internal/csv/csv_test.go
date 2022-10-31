// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package csv

import (
	"context"
	"errors"
	"testing"

	"github.com/alpstable/gidari/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("directory does not exist", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, err := New(ctx, "testdata-dne")
		if !errors.Is(err, ErrNoDir) {
			t.Fatalf("expected ErrNoDir, got %v", err)
		}
	})

	t.Run("directory exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, err := New(ctx, "testdata")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func sameStringSlice(t *testing.T, arr1, arr2 []string) bool {
	t.Helper()

	if len(arr1) != len(arr2) {
		return false
	}

	// create a map of string -> int
	diff := make(map[string]int, len(arr1))
	for _, val := range arr1 {
		// 0 value for int is 0, so just increment a counter for the string
		diff[val]++
	}

	for _, val := range arr2 {
		// If the string _y is not in diff bail out early
		if _, ok := diff[val]; !ok {
			return false
		}

		diff[val]--
		if diff[val] == 0 {
			delete(diff, val)
		}
	}

	return len(diff) == 0
}

func TestFlattenStruct(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name   string
		table  string
		object *structpb.Struct
		want   map[string]string
	}{
		{
			name:  "empty",
			table: "test",
			object: func() *structpb.Struct {
				spb, _ := structpb.NewStruct(map[string]interface{}{})

				return spb
			}(),
			want: map[string]string{},
		},
		{
			name:  "one field",
			table: "test",
			object: func() *structpb.Struct {
				spb, _ := structpb.NewStruct(map[string]interface{}{
					"foo": "bar",
				})

				return spb
			}(),
			want: map[string]string{
				"foo": "bar",
			},
		},
		{
			name:  "one field as struct",
			table: "test",
			object: func() *structpb.Struct {
				spb, _ := structpb.NewStruct(map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "baz",
					},
				})

				return spb
			}(),
			want: map[string]string{
				"foo.bar": "baz",
			},
		},
		{
			name:  "many fields",
			table: "test",
			object: func() *structpb.Struct {
				spb, _ := structpb.NewStruct(map[string]interface{}{
					"foo": "bar",
					"baz": "qux",
					"quux": map[string]interface{}{
						"corge":  "grault",
						"garply": true,
						"waldo":  nil,
					},
					"garply": 1,
				})

				return spb
			}(),
			want: map[string]string{
				"foo":         "bar",
				"baz":         "qux",
				"quux.corge":  "grault",
				"quux.garply": "true",
				"quux.waldo":  "",
				"garply":      "1.000000",
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := flattenStruct(tcase.object)
			if err != nil {
				t.Fatal(err)
			}

			if len(got) != len(tcase.want) {
				t.Fatalf("expected %d fields, got %d", len(tcase.want), len(got))
			}

			for k, v := range tcase.want {
				if got[k] != v {
					t.Fatalf("expected %s, got %s", v, got[k])
				}
			}
		})
	}
}

func TestDecodeUpsertRequest(t *testing.T) {
	t.Parallel()

	testTable := &proto.Table{Name: "test"}

	for _, tcase := range []struct {
		name string
		reqs []*proto.UpsertRequest
		want map[string][]string
	}{
		{
			name: "empty",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(``),
				},
			},
			want: map[string][]string{},
		},
		{
			name: "single row",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`{"id":1,"name":"test"}`),
				},
			},
			want: map[string][]string{
				"id":   {"1.000000"},
				"name": {"test"},
			},
		},
		{
			name: "single row object",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`{"id":1,"properties":{"name":"test"}}`),
				},
			},
			want: map[string][]string{
				"id":              {"1.000000"},
				"properties.name": {"test"},
			},
		},
		{
			name: "smaller row before larger row",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"name":"test","age":10}]`),
				},
			},
			want: map[string][]string{
				"id":   {"1.000000", "2.000000"},
				"name": {"test", "test"},
				"age":  {"", "10.000000"},
			},
		},
		{
			name: "smaller row before larger row unordered",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
				},
			},
			want: map[string][]string{
				"id":   {"1.000000", "2.000000"},
				"name": {"test", "test"},
				"age":  {"", "10.000000"},
			},
		},
		{
			name: "empty leading row",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`[{},{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
				},
			},
			want: map[string][]string{
				"id":   {"", "1.000000", "2.000000"},
				"name": {"", "test", "test"},
				"age":  {"", "", "10.000000"},
			},
		},
		{
			name: "sparse rows with uneven columns",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data: []byte(`[{},{"id":1,"name":"test"},{"other":"test"},
{"age":10,"name":"test"}]`),
				},
			},
			want: map[string][]string{
				"id":    {"", "1.000000", "", ""},
				"name":  {"", "test", "", "test"},
				"age":   {"", "", "", "10.000000"},
				"other": {"", "", "test", ""},
			},
		},
		{
			name: "multiple requests",
			reqs: []*proto.UpsertRequest{
				{
					Table: testTable,
					Data:  []byte(`{"id":1,"name":"test"}`),
				},
				{
					Table: testTable,
					Data:  []byte(`{"id":1,"name":"test","age":10}`),
				},
				{
					Table: testTable,
					Data:  []byte(`{"id":1}`),
				},
				{
					Table: testTable,
					Data:  []byte(`{"other":"test"}`),
				},
			},
			want: map[string][]string{
				"id":    {"1.000000", "1.000000", "1.000000", ""},
				"name":  {"test", "test", "", ""},
				"age":   {"", "10.000000", "", ""},
				"other": {"", "", "", "test"},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			state := newWriteState()
			channeledRows := make(map[string][][]string, 0)

			for _, req := range tcase.reqs {
				rows := make(chan *row)

				errs, _ := errgroup.WithContext(context.Background())
				errs.Go(func() error {
					return decodeUpsertRequest(req, state, rows)
				})

				for row := range rows {
					channeledRows[req.Table.Name] = append(channeledRows[req.Table.Name], row.data)
				}

				if err := errs.Wait(); err != nil {
					t.Fatal(err)
				}
			}

			// If no rows were channeled, we're done.
			if len(channeledRows) == 0 {
				return
			}

			for table, headerRow := range state.headerRowByTable {
				for idx, got := range channeledRows[table] {
					want := []string{}
					for _, headerName := range headerRow.data {
						if idx < len(tcase.want[headerName]) {
							want = append(want, tcase.want[headerName][idx])
						}
					}

					// Add trailing empty columns to "got" if it's shorter than "want".
					for i := len(got); i < len(want); i++ {
						got = append(got, "")
					}

					if !sameStringSlice(t, got, want) {
						t.Errorf("unexpected header row for table %q: got %v, want %v",
							table, got, want)
					}
				}
			}
		})
	}
}
