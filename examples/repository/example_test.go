package repository_test

import (
	"context"
	"fmt"
	"os"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/repository"
	"github.com/alpine-hodler/gidari/tools"
	_ "github.com/joho/godotenv/autoload"
)

func ExampleNew() {
	// cluster_uri looks something like 'mongodb+srv://<username>:<password>@<clustername><host>/<Collection>'
	dsn := os.Getenv("MONGO_URI")
	ctx := context.TODO()

	repo, err := repository.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	fmt.Println(repo.Storage.IsNoSQL())
	// Output:
	// true
}

// TODO
// Didn't really know how to make a meaningful example for NewTx(), but heres some boilerplate to get statrted

// func ExampleNewTx() {
// 	dsn := os.Getenv("MONGO_URI")
// 	ctx := context.TODO()
//
// 	txRepo, err := repository.NewTx(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	txInit, err := txRepo.StartTx(ctx)
//
// 	if condition {
// 		txInit.Rollback()
// 	}
// 	txInit.Commit()
// 	// Output:
// 	//TODO
// }

func ExampleTruncate() {
	ctx := context.Background()
	dns := os.Getenv("MONGO_URI")

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
	// 0
}

func ExampleUpsert() {
	ctx := context.Background()
	dns := os.Getenv("MONGO_URI")

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
	// 0
}

func ExampleListTables() {
	var err error

	ctx := context.TODO()
	dsn := os.Getenv("MONGO_URI")

	repo, err := repository.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	rsp, err := repo.ListTables(ctx)
	if err != nil {
		panic(err)
	}

	for table := range rsp.TableSet {
		fmt.Println(table)
	}
	// Output:
	// AnotherExampleTable
	// ExampleTable
}

func ExampleListPrimaryKeys() {
	var err error

	ctx := context.TODO()
	dsn := os.Getenv("MONGO_URI")

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
	// 2
}
