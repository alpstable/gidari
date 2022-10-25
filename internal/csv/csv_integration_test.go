package csv

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpstable/gidari/internal/proto"
)

func TestUpsert(t *testing.T) {
	ctx := context.Background()
	csv, err := New(ctx, "testdata")
	if err != nil {
		t.Fatal(err)
	}

	req := &proto.UpsertRequest{
		Table: "test",
		Data:  []byte(`{"id": 1, "name": "test"}`),
	}

	ch := make(chan *proto.UpsertRequest)
	_, errGroup := csv.Upsert(ctx, ch)

	ch <- req
	close(ch)

	if err := errGroup.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestComplexUpsert(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		reqs []*proto.UpsertRequest
		want [][]string
	}{
		{
			name: "single row",
			reqs: []*proto.UpsertRequest{
				{
					Table: "test_6b3fe527-4268-4b4d-8477-2da84df678c6",
					Data:  []byte(`{"id":1,"name":"test","info":{"age":10,"address":"test"}}`),
				},
			},
			want: [][]string{
				{"id", "name", "info.age", "info.address"},
				{"1.000000", "test", "10.000000", "test"},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			csv, err := New(ctx, "testdata")
			if err != nil {
				t.Fatal(err)
			}

			reqCh := make(chan *proto.UpsertRequest)
			_, errGroup := csv.Upsert(ctx, reqCh)

			tables := make(map[string]struct{})

			for _, req := range tcase.reqs {
				if _, ok := tables[req.Table]; !ok {
					t.Cleanup(func() {
						filename := filepath.Join("testdata", req.Table+".csv")
						if err := os.Remove(filename); err != nil {
							t.Fatal(err)
						}
					})
				}

				tables[req.Table] = struct{}{}
				reqCh <- req
			}
			close(reqCh)

			if err := errGroup.Wait(); err != nil {
				t.Fatalf("failed to upsert: %v", err)
			}

			for table := range tables {
				filename := filepath.Join("testdata", table+".csv")

				// Get the lines from the file
				file, err := os.Open(filename)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				linNum := 0

				for scanner.Scan() {
					line := scanner.Text()
					if !sameStringSlice(t, strings.Split(line, ","), tcase.want[linNum]) {
						t.Fatalf("got %v, want %v", strings.Split(line, ","), tcase.want[linNum])
					}

					linNum++
				}

			}
		})
	}
}
