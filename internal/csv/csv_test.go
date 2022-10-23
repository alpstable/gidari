package csv

import (
	"context"
	"testing"

	"github.com/alpstable/gidari/internal/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("directory does not exist", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, err := New(ctx, "testdata-dne")
		if err != ErrNoDir {
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

func sameStringSlice(t *testing.T, x, y []string) bool {
	t.Helper()

	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
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

	for _, tcase := range []struct {
		name string
		reqs []*proto.UpsertRequest
		want [][]string
	}{
		{
			name: "empty",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(``),
				},
			},
		},
		{
			name: "single row",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`{"id":1,"name":"test"}`),
				},
			},
			want: [][]string{
				{"id", "name"},
				{"1.000000", "test"},
			},
		},
		{
			name: "single row object",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`{"id":1,"properties":{"name":"test"}}`),
				},
			},
			want: [][]string{
				{"id", "properties.name"},
				{"1.000000", "test"},
			},
		},
		{
			name: "smaller row before larger row",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"name":"test","age":10}]`),
				},
			},
			want: [][]string{
				{"id", "name", "age"},
				{"1.000000", "test"},
				{"2.000000", "test", "10.000000"},
			},
		},
		{
			name: "smaller row before larger row unordered",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
				},
			},
			want: [][]string{
				{"id", "name", "age"},
				{"1.000000", "test"},
				{"2.000000", "test", "10.000000"},
			},
		},
		{
			name: "empty leading row",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`[{},{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
				},
			},
			want: [][]string{
				{"id", "name", "age"},
				{},
				{"1.000000", "test"},
				{"2.000000", "test", "10.000000"},
			},
		},
		{
			name: "sparse rows with uneven columns",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data: []byte(`[{},{"id":1,"name":"test"},{"other":"test"},
{"age":10,"name":"test"}]`),
				},
			},
			want: [][]string{
				{"id", "name", "other", "age"},
				{},
				{"1.000000", "test"},
				{"", "", "test"},
				{"", "test", "", "10.000000"},
			},
		},
		{
			name: "multiple requests",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test",
					Data:  []byte(`{"id":1,"name":"test"}`),
				},
				{
					Table: "test",
					Data:  []byte(`{"id":1,"name":"test","age":10}`),
				},
				{
					Table: "test",
					Data:  []byte(`{"id":1}`),
				},
				{
					Table: "test",
					Data:  []byte(`{"other":"test"}`),
				},
			},
			want: [][]string{
				{"id", "name", "age", "other"},
				{"1.000000", "test"},
				{"1.000000", "test", "10.000000"},
				{"1.000000"},
				{"", "", "", "test"},
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
				go decodeUpsertRequest(req, state, rows)

				for row := range rows {
					channeledRows[req.Table] = append(channeledRows[req.Table], row.data)
				}
			}

			// If no rows were channeled, we're done.
			if len(channeledRows) == 0 {
				return
			}

			for table, headerRow := range state.headerRowByTable {
				if !sameStringSlice(t, headerRow.data, tcase.want[0]) {
					t.Errorf("unexpected header row: got %v, want %v", headerRow, tcase.want[0])
				}

				for idx, got := range channeledRows[table] {
					if !sameStringSlice(t, got, tcase.want[idx+1]) {
						t.Errorf("unexpected header row for table %q: got %v, want %v",
							table, got, tcase.want[idx+1])
					}
				}
			}
		})
	}
}
