package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func TestMongoDBTxn(t *testing.T) {
	t.Parallel()

	t.Run("txns should reset with when the lifetime is reached", func(t *testing.T) {
		t.Parallel()

		const collection = "test-ceebf"
		const database = "ltest"
		const tolerance = 50_000

		ctx := context.Background()
		mdb, err := NewMongo(ctx, fmt.Sprintf("mongodb://mongo1:27017/%s", database))
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
			txn.Send(func(sctx context.Context, stg Storage) error {
				_, err := stg.Upsert(sctx, &proto.UpsertRequest{
					Table:    collection,
					Data:     bytes,
					DataType: int32(tools.UpsertDataJSON),
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
	})
}
