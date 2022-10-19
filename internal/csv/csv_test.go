package csv

import (
	"context"
	"reflect"
	"testing"

	"github.com/alpstable/gidari/internal/proto"
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

func sameInterfaceSlice(x, y []interface{}) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[interface{}]int, len(x))
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

func TestDecodeUpsertRequest(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		req  *proto.UpsertRequest
		want map[string][]interface{}
	}{
		{
			name: "empty",
			req: &proto.UpsertRequest{
				Table: "test",
				Data:  []byte(``),
			},
			want: map[string][]interface{}{},
		},
		{
			name: "single row",
			req: &proto.UpsertRequest{

				Table: "test",
				Data:  []byte(`{"id": 1, "name": "test"}`),
			},
			want: map[string][]interface{}{
				"id":   {1.0},
				"name": {"test"},
			},
		},
		{
			name: "smaller row before larger row",
			req: &proto.UpsertRequest{

				Table: "test",
				Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"name":"test","age":10}]`),
			},
			want: map[string][]interface{}{
				"id":   {1.0, 2.0},
				"name": {"test", "test"},
				"age":  {nil, 10.0},
			},
		},
		{
			name: "smaller row before larger row unordered",
			req: &proto.UpsertRequest{

				Table: "test",
				Data:  []byte(`[{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
			},
			want: map[string][]interface{}{
				"id":   {1.0, 2.0},
				"name": {"test", "test"},
				"age":  {nil, 10.0},
			},
		},
		{
			name: "empty leading row",
			req: &proto.UpsertRequest{

				Table: "test",
				Data:  []byte(`[{},{"id":1,"name":"test"},{"id":2,"age":10,"name":"test"}]`),
			},
			want: map[string][]interface{}{
				"id":   {nil, 1.0, 2.0},
				"name": {nil, "test", "test"},
				"age":  {nil, nil, 10.0},
			},
		},
		{
			name: "sparse rows with uneven columns",
			req: &proto.UpsertRequest{

				Table: "test",
				Data:  []byte(`[{},{"id":1,"name":"test"},{"other":"test"},{"age":10,"name":"test"}]`),
			},
			want: map[string][]interface{}{
				"id":    {nil, 1.0, nil, nil},
				"name":  {nil, "test", nil, "test"},
				"other": {nil, nil, "test", nil},
				"age":   {nil, nil, nil, 10.0},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			rows := make(chan []interface{})
			go decodeUpsertRequest(tcase.req, rows)

			channeledRows := make([][]interface{}, 0)
			for row := range rows {
				channeledRows = append(channeledRows, row)
			}

			// If no rows were channeled, we're done.
			if len(channeledRows) == 0 {
				return
			}

			// The last row is the header row. We need this to compare.
			for idx, headerName := range channeledRows[len(channeledRows)-1] {
				column := make([]interface{}, 0)

				for _, row := range channeledRows[:len(channeledRows)-1] {
					if len(row) <= idx {
						column = append(column, nil)

						continue
					}

					column = append(column, row[idx])
				}

				if !reflect.DeepEqual(column, tcase.want[headerName.(string)]) {
					t.Fatalf("expected %v, got %v", tcase.want[headerName.(string)], column)
				}
			}
		})
	}
}
