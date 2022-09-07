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

// ExecTx executes a function within a database transaction.
func (m *Mongo) ExecTx(ctx context.Context, fn func(context.Context, tools.GenericStorage) (bool, error)) error {
	return m.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// start the transactions
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		ok, err := fn(sessionContext, m)
		if err != nil {
			return err
		}
		if !ok {
			// rollback the transactions so the test db remains clean.
			if err := sessionContext.AbortTransaction(sessionContext); err != nil {
				return err
			}
			return nil
		}

		sessionContext.EndSession(ctx)
		return nil
	})
}

// StartTx will start a mongodb session where all data from write methods can be rolled back.
func (m *Mongo) StartTx(ctx context.Context) Tx {
	tx := &Tx{
		make(chan func(context.Context) error),
		make(chan error, 1),
		make(chan bool, 1),
		make(chan bool, 1),
	}

	var globalSessionContext mongo.SessionContext

	go func() {
	DecisionLoop:
		for {
			switch {
			case <-tx.rollback:
				globalSessionContext.AbortTransaction(globalSessionContext)
				break DecisionLoop
			case <-tx.commit:
				if err := globalSessionContext.CommitTransaction(globalSessionContext); err != nil {
					tx.rollback <- true
					tx.done <- err
				}
				break DecisionLoop
			}
		}
		tx.done <- nil
	}()

	// Create a go routine that creates a session and listens for writes.
	go func() {
		m.UseSession(ctx, func(sctx mongo.SessionContext) error {
			globalSessionContext = sctx

			// start the transaction, if there is an error break the go routine.
			err := sctx.StartTransaction()
			if err != nil {
				tx.done <- err
				return err
			}

			// listen for writes.
			for fn := range tx.Ch {
				if err != nil {
					continue
				}

				// execute the write function. If there is an error, rollback the transaction.
				if err = fn(sctx); err != nil {
					tx.rollback <- true
					tx.done <- err
				}
			}
			return nil
		})
	}()
	return *tx
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
func (m *Mongo) Upsert(ctx context.Context, req *proto.UpsertRequest, rsp *proto.CreateResponse) error {
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
	_, err = coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	return nil
}
