package storage

import (
	"context"

	"github.com/alpine-hodler/gidari/pkg/proto"
	"github.com/alpine-hodler/gidari/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Mongo is a wrapper for *mongo.Client, use to perform CRUD operations on a mongo DB instance.
type Mongo struct {
	*mongo.Client
	dns string
}

// NewMongo will return a new mongo client that can be used to perform CRUD operations on a mongo DB instance. This
// constructor uses a URI to make the client connection, and the URI is of the form
// Mongo://username:password@host:port
func NewMongo(ctx context.Context, uri string) (*Mongo, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	mdb := new(Mongo)
	mdb.Client = client
	mdb.dns = uri

	return mdb, nil
}

// Type returns the type of storage.
func (m *Mongo) Type() uint8 {
	return MongoType
}

// Close will close the mongo client.
func (m *Mongo) Close() {
	m.Close()
}

// StartTx will start a mongodb session where all data from write methods can be rolled back.
func (m *Mongo) StartTx(ctx context.Context) (Tx, error) {
	// Construct a transaction.
	tx := &tx{
		make(chan func(context.Context) error),
		make(chan error, 1),
		make(chan bool, 1),
	}

	// Create a go routine that creates a session and listens for writes.
	go func() {
		tx.done <- m.UseSession(ctx, func(sctx mongo.SessionContext) error {
			// Start the transaction, if there is an error break the go routine.
			err := sctx.StartTransaction()
			if err != nil {
				return err
			}

			// listen for writes.
			for fn := range tx.ch {
				// If an error has registered, do nothing.
				if err != nil {
					continue
				}
				err = fn(sctx)
			}

			if err != nil {
				return err
			}

			// Await the decision to commit or rollback.
			switch {
			case <-tx.commit:
				if err := sctx.CommitTransaction(sctx); err == nil {
					return err
				}
			default:
				sctx.AbortTransaction(sctx)
			}
			return nil
		})
	}()
	return tx, nil
}

func (m *Mongo) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	// bldr, err := query.GetReadBuilder(query.ReadBuilderType(req.ReaderBuilder[0]))
	// if err != nil {
	// 	return err
	// }

	// args, err := bldr.ReaderArgs(req)
	// if err != nil {
	// 	return err
	// }
	// filterbytes, err := bldr.ReaderQuery(query.MongoStorage, args...)
	// if err != nil {
	// 	return err
	// }

	// var outputBuffer bytes.Buffer
	// outputBuffer.Write(filterbytes)

	// q := query.Mongo{}
	// if err = gob.NewDecoder(&outputBuffer).Decode(&q); err != nil {
	// 	return err
	// }

	// cs, err := connstring.ParseAndValidate(m.dns)
	// if err != nil {
	// 	return nil
	// }

	// coll := m.Database(cs.Database).Collection(q.Collection)
	// cursor, err := coll.Find(ctx, q.D)
	// if err != nil {
	// 	return err
	// }

	// for cursor.Next(ctx) {
	// 	m := make(map[string]interface{})
	// 	err := cursor.Decode(&m)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	delete(m, "_id")
	// 	record, err := structpb.NewStruct(m)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	rsp.Records = append(rsp.Records, record)
	// }
	return nil
}

func (m *Mongo) TruncateTables(context.Context, *proto.TruncateTablesRequest) error { return nil }

// UpsertCoinbaseProCandles60 will upsert candles to the 60-granularity Mongo DB collection for a given productID.
func (m *Mongo) Upsert(ctx context.Context, req *proto.UpsertRequest, rsp *proto.UpsertResponse) error {
	models := []mongo.WriteModel{}
	for _, record := range req.Records {
		doc := bson.D{}
		if err := tools.AssingRecordBSONDocument(record, &doc); err != nil {
			return err
		}
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(doc).
			SetUpdate(bson.D{{"$set", doc}}).
			SetUpsert(true))
	}

	cs, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return err
	}

	coll := m.Database(cs.Database).Collection(req.Table)
	bwr, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	rsp.UpsertedCount = int32(bwr.UpsertedCount)
	rsp.MatchedCount = int32(bwr.MatchedCount)
	return nil
}
