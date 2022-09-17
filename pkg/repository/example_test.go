package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/alpine-hodler/gidari/pkg/proto"
	"github.com/alpine-hodler/gidari/pkg/repository"
)

// testCase is the configuration for the example test.
type testCase struct {
	databaseURL string
}

func TestExamples(t *testing.T) {
	cases := []testCase{
		{databaseURL: "mongodb://mongo:27017/coinbasepro"},
	}

	for _, tc := range cases {
		err := os.Setenv("DATABASE_URL", tc.databaseURL)
		if err != nil {
			t.Fatalf("failed to set environment variable: %v", err)
		}

		t.Run(fmt.Sprintf("ExampleGenericService_UpsertRawJSON databaseURL=%s", tc.databaseURL),
			func(t *testing.T) {
				ExampleGenericService_UpsertRawJSON()
			})
	}

}

func ExampleGenericService_UpsertRawJSON() {
	ctx := context.Background()
	dns := os.Getenv("DATABASE_URL")

	repo, err := repository.New(ctx, dns)
	if err != nil {
		panic(err)
	}

	raw := &repository.Raw{
		Table: "accounts",
		Data:  []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
	}

	rsp := new(proto.UpsertResponse)
	err = repo.UpsertRawJSON(ctx, raw, rsp)
	if err != nil {
		panic(err)
	}
}
