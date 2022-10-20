//go:build mdbinteg

// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alpstable/gidari/internal/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const defaultConnectionString = "mongodb://mongo1:27017/defaultcoll"

func TestMongo(t *testing.T) {
	t.Parallel()

	const defaultTestTable = "tests1"
	const listTablesTable = "lttests1"
	const listPrimaryKeysTable = "pktests1"

	defaultData := map[string]interface{}{
		"test_string": "test",
		"id":          "1",
	}

	ctx := context.Background()

	mongo, err := New(ctx, defaultConnectionString)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	proto.RunTest(context.Background(), t, mongo, func(runner *proto.TestRunner) {
		runner.AddCloseDBCases(
			[]proto.TestCase{
				{
					Name: "close mongo",
				},
			}...,
		)

		runner.AddPingDBCases(
			[]proto.TestCase{
				{
					Name: "check mongo connection"
				}
			}...,
		)

		runner.AddListPrimaryKeysCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listPrimaryKeysTable,
					ExpectedPrimaryKeys: map[string][]string{
						listPrimaryKeysTable: {"_id"},
					},
				},
			}...,
		)

		runner.AddListTablesCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listTablesTable,
				},
			}...,
		)

		runner.AddUpsertTxnCases(
			[]proto.TestCase{
				{
					Name:               "commit",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 54,
					Data:               defaultData,
				},
				{
					Name:               "rollback",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					Rollback:           true,
					Data:               defaultData,
				},
				{
					Name:               "rollback on error",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					ForceError:         true,
					Data:               defaultData,
				},
			}...,
		)
	})
}

func TestMongoDBTxn(t *testing.T) {
	t.Parallel()

	t.Run("txns should reset with when the lifetime is reached", func(t *testing.T) {
		t.Parallel()

		const collection = "test-ceebf"
		const database = "ltest"
		const tolerance = 5_000

		ctx := context.Background()
		mdb, err := New(ctx, fmt.Sprintf("mongodb://mongo1:27017/%s", database))
		if err != nil {
			t.Fatalf("failed to create mongo client: %v", err)
		}

		// Change the lifetime to 1 second to avoid long test times.
		mdb.lifetime = 1 * time.Second

		// Start a transaction.
		txn, err := mdb.StartTx(ctx)
		if err != nil {
			t.Fatalf("failed to start txn: %v", err)
		}

		// Create some data that we will encode into bytes to insert into the db in bulk.
		data := map[string]interface{}{"test_string": "test"}
		bytes, err := json.Marshal(data)
		if err != nil {
			t.Fatalf("failed to marshal data: %v", err)
		}

		// Add an index to the collection.
		indexView := mdb.Client.Database(database).Collection(collection).Indexes()
		_, err = indexView.CreateOne(context.Background(), mongo.IndexModel{
			Keys: bsonx.Doc{{Key: "test_string", Value: bsonx.Int32(1)}},
		})
		if err != nil {
			t.Fatalf("failed to create index: %v", err)
		}

		for i := 0; i < tolerance; i++ {
			if i%10_000 == 0 && i != 0 {
				t.Logf("inserted %d documents", i)
			}

			// Insert some data.
			txn.Send(func(sctx context.Context, stg proto.Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table: collection,
					Data:  bytes,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert data: %w", err)
				}

				return nil
			})
		}

		if err := txn.Commit(); err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}

		// Truncate the collection.
		treq := &proto.TruncateRequest{Tables: []string{collection}}
		if _, err := mdb.Truncate(ctx, treq); err != nil {
			t.Fatalf("failed to truncate collection: %v", err)
		}
	})
}
