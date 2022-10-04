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

// The Example configuration used in for examples is a MongoDB replica set and has the following:
// - One database: "GidariExample".
// - Three collections: "ExampleTable", "TxnExampleTable", "AnotherExampleTable".
func TestExamples(t *testing.T) {
	t.Cleanup(tools.Quiet()) // This will suppress the `fmt` print statements in CI
	t.Parallel()

	for _, tcase := range []struct{ mongoURI string }{
		{"mongodb://mongo1:27017/repositoryExamples"},
	} {
		err := os.Setenv("DB_CONN_STRING", tcase.mongoURI)
		if err != nil {
			t.Fatalf("failed to set environment variable: %v", err)
		}

		// Register the example tests

		t.Run("Example_New",
			func(t *testing.T) {
				t.Parallel()
				ExampleNew()
			})

		t.Run("Example_NewTx",
			func(t *testing.T) {
				t.Parallel()
				ExampleNewTx()
			})

		t.Run("ExampleGenericService_Truncate",
			func(t *testing.T) {
				t.Parallel()
				ExampleGenericService_Truncate()
			})

		t.Run("ExampleGenericService_Upsert",
			func(t *testing.T) {
				t.Parallel()
				ExampleGenericService_Upsert()
			})

		t.Run("ExampleGenericService_ListTables",
			func(t *testing.T) {
				t.Parallel()
				ExampleGenericService_ListTables()
			})

		t.Run("ExampleGenericService_ListPrimaryKeys",
			func(t *testing.T) {
				t.Parallel()
				ExampleGenericService_ListPrimaryKeys()
			})
	}
}

func ExampleNew() {
	dsn := os.Getenv("DB_CONN_STRING")
	ctx := context.TODO()

	repo, err := repository.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println(repo.Storage.IsNoSQL())
	// Output:
	// true
}

func ExampleNewTx() {
	dsn := os.Getenv("DB_CONN_STRING")
	ctx := context.TODO()

	txRepo, err := repository.NewTx(ctx, dsn)
	if err != nil {
		panic(err)
	}

	req := &proto.UpsertRequest{
		Table:    "TxnExampleTable",
		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
		DataType: int32(tools.UpsertDataJSON),
	}

	rsp, err := txRepo.Upsert(ctx, req)
	if err != nil {
		panic(err)
	}

	if err := txRepo.Commit(); err != nil {
		panic(err)
	}

	fmt.Println(rsp.GetUpsertedCount())

	// Output:
	// Not Deterministic
}

func ExampleGenericService_Truncate() {
	ctx := context.Background()
	dns := os.Getenv("DB_CONN_STRING")

	repo, err := repository.New(ctx, dns)
	if err != nil {
		panic(err)
	}

	req := &proto.TruncateRequest{
		Tables: []string{"ExampleTable"},
	}

	rsp, err := repo.Truncate(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println(rsp.DeletedCount)

	// Output:
	// Not Deterministic
}

func ExampleGenericService_Upsert() {
	ctx := context.Background()
	dns := os.Getenv("DB_CONN_STRING")

	repo, err := repository.New(ctx, dns)
	if err != nil {
		panic(err)
	}

	req := &proto.UpsertRequest{
		Table:    "ExampleTable",
		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
		DataType: int32(tools.UpsertDataJSON),
	}

	rsp, err := repo.Upsert(ctx, req)
	if err != nil {
		panic(err)
	}

	tally := rsp.GetMatchedCount()
	fmt.Println(tally)

	// Output:
	// Not Deterministic
}

func ExampleGenericService_ListTables() {
	var err error

	ctx := context.TODO()
	dsn := os.Getenv("DB_CONN_STRING")

	repo, err := repository.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	rsp, err := repo.ListTables(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(rsp.TableSet))

	// Output:
	// Not Deterministic
}

func ExampleGenericService_ListPrimaryKeys() {
	var err error

	ctx := context.TODO()
	dsn := os.Getenv("DB_CONN_STRING")

	repo, err := repository.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	rsp, err := repo.ListPrimaryKeys(ctx)
	if err != nil {
		panic(err)
	}

	totalPKeys := 0
	for _, keys := range rsp.PKSet {
		totalPKeys += len(keys.GetList())
	}

	fmt.Println(totalPKeys)

	// Output:
	// Not Deterministic
}
