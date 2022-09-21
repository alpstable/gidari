package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/repository"
	"github.com/alpine-hodler/gidari/tools"
)

// testCase is the configuration for the example test.
type testCase struct {
	databaseURL string
}

func TestExamples(t *testing.T) {
	t.Cleanup(tools.Quiet())

	t.Parallel()

	cases := []testCase{
		{databaseURL: "mongodb://mongo1:27017/coinbasepro"},
	}

	for _, tcase := range cases {
		err := os.Setenv("DATABASE_URL", tcase.databaseURL)
		if err != nil {
			t.Fatalf("failed to set environment variable: %v", err)
		}

		t.Run(fmt.Sprintf("ExampleGenericService_UpsertRawJSON databaseURL=%s", tcase.databaseURL),
			func(t *testing.T) {
				t.Parallel()
				ExampleGenericService_Upsert()
			})
	}
}

func ExampleGenericService_Upsert() {
	ctx := context.Background()
	dns := os.Getenv("DATABASE_URL")

	repo, err := repository.New(ctx, dns)
	if err != nil {
		panic(err)
	}

	req := &proto.UpsertRequest{
		Table:    "accounts",
		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
		DataType: int32(tools.UpsertDataJSON),
	}

	rsp, err := repo.Upsert(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("upserted %d rows\n", rsp.UpsertedCount)
}
