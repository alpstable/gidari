// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package repository_test

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"sync"
// 	"testing"

// 	"github.com/alpstable/gidari/internal/repository"
// 	"github.com/alpstable/gidari/proto"
// 	"github.com/alpstable/gidari/tools"
// )

// func truncateGenericRepo(ctx context.Context, t *testing.T, connectionString string, tables ...string) {
// 	t.Helper()

// 	t.Cleanup(func() {
// 		repo, err := repository.New(ctx, connectionString)
// 		if err != nil {
// 			t.Fatalf("failed to create repository: %v", err)
// 		}

// 		if _, err := repo.Truncate(ctx, &proto.TruncateRequest{Tables: tables}); err != nil {
// 			t.Fatalf("failed to truncate storage: %v", err)
// 		}
// 	})
// }

// // The Example configuration used in for examples is a MongoDB replica set and has the following:
// // - One database: "GidariExample".
// // - Three collections: "ExampleTable", "TxnExampleTable", "AnotherExampleTable".
// func TestExamples(t *testing.T) {
// 	t.Cleanup(tools.Quiet()) // This will suppress the `fmt` print statements in CI
// 	t.Parallel()

// 	mux := sync.Mutex{}

// 	for _, tcase := range []struct {
// 		name     string
// 		mongoURI string
// 		table    string
// 		ex       func()
// 	}{
// 		{"New", "mongodb://mongo1:27017/coll1", "table1", ExampleNew},
// 		{"NewTx", "mongodb://mongo1:27017/coll2", "table2", ExampleNewTx},
// 		{"Truncate", "mongodb://mongo1:27017/coll3", "table3", ExampleGenericService_Truncate},
// 		{"Upsert", "mongodb://mongo1:27017/coll4", "table4", ExampleGenericService_Upsert},
// 		{"ListTables", "mongodb://mongo1:27017/coll5", "table5", ExampleGenericService_ListTables},
// 	} {
// 		mux.Lock()

// 		tcase := tcase

// 		err := os.Setenv("MONGODB_URI", tcase.mongoURI)
// 		if err != nil {
// 			t.Fatalf("failed to set environment variable: %v", err)
// 		}

// 		// Register the example tests
// 		t.Run(tcase.name, func(t *testing.T) {
// 			t.Parallel()

// 			ctx := context.Background()
// 			truncateGenericRepo(ctx, t, tcase.mongoURI, tcase.table)

// 			tcase.ex()
// 		})

// 		mux.Unlock()
// 	}
// }

// func ExampleNew() {
// 	dsn := os.Getenv("MONGODB_URI")
// 	ctx := context.TODO()

// 	repo, err := repository.New(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(repo.Storage.IsNoSQL())
// 	// Output:
// 	// true
// }

// func ExampleNewTx() {
// 	dsn := os.Getenv("MONGODB_URI")
// 	ctx := context.TODO()

// 	table := "table2"

// 	txRepo, err := repository.NewTx(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req := &proto.UpsertRequest{
// 		Table:    table,
// 		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
// 		DataType: int32(tools.UpsertDataJSON),
// 	}

// 	_, err = txRepo.Upsert(ctx, req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err := txRepo.Commit(); err != nil {
// 		panic(err)
// 	}

// 	repo, err := repository.New(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tresp, err := repo.ListTables(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the table size in bits
// 	fmt.Println(tresp.TableSet[table].Size)

// 	// Output:
// 	// 67
// }

// func ExampleGenericService_Truncate() {
// 	ctx := context.Background()
// 	dns := os.Getenv("MONGODB_URI")

// 	table := "table3"

// 	repo, err := repository.New(ctx, dns)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req := &proto.TruncateRequest{
// 		Tables: []string{table},
// 	}

// 	rsp, err := repo.Truncate(ctx, req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(rsp.GetDeletedCount())

// 	// Output:
// 	// 0
// }

// func ExampleGenericService_Upsert() {
// 	ctx := context.Background()
// 	dns := os.Getenv("MONGODB_URI")

// 	table := "table4"

// 	repo, err := repository.New(ctx, dns)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req := &proto.UpsertRequest{
// 		Table:    table,
// 		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
// 		DataType: int32(tools.UpsertDataJSON),
// 	}

// 	_, err = repo.Upsert(ctx, req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tresp, err := repo.ListTables(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the table size in bits
// 	fmt.Println(tresp.TableSet[table].Size)

// 	// Output:
// 	// 67
// }

// func ExampleGenericService_ListTables() {
// 	var err error

// 	ctx := context.TODO()
// 	dsn := os.Getenv("MONGODB_URI")

// 	repo, err := repository.New(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rsp, err := repo.ListTables(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tables := []string{}
// 	for table := range rsp.TableSet {
// 		tables = append(tables, table)
// 	}

// 	fmt.Println(len(tables) > 1)

// 	// Output:
// 	// true
// }

// func ExampleGenericService_ListPrimaryKeys() {
// 	var err error

// 	ctx := context.TODO()
// 	dsn := os.Getenv("MONGODB_URI")

// 	table := "table6"

// 	repo, err := repository.New(ctx, dsn)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Upsert data so that we have a collection to query on.
// 	req := &proto.UpsertRequest{
// 		Table:    table,
// 		Data:     []byte(`[{"id": "7fd0abc0-e5ad-4cbb-8d54-f2b3f43364da"}]`),
// 		DataType: int32(tools.UpsertDataJSON),
// 	}

// 	_, err = repo.Upsert(ctx, req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rsp, err := repo.ListPrimaryKeys(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	tablePKs := rsp.GetPKSet()[table].GetList()
// 	fmt.Println(tablePKs[0])

// 	// Output:
// 	// _id
// }
