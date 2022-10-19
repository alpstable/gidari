package csv

import (
	"context"
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

	_, err = csv.Upsert(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
}
