package storage

import (
	"context"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err := m.Client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

// StartTx will start a mongodb session where all data from write methods can be rolled back.
func (m *Mongo) StartTx(ctx context.Context) (Tx, error) {
	// Construct a transaction.
	tx := &tx{
		make(chan TXChanFn),
		make(chan error, 1),
		make(chan bool, 1),
	}

	// Create a go routine that creates a session and listens for writes.
	go func() {
		tx.done <- m.Client.UseSession(ctx, func(sctx mongo.SessionContext) error {
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
				err = fn(sctx, m)
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
				if err := sctx.AbortTransaction(sctx); err == nil {
					return err
				}
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

// Truncate will delete all records in a collection.
func (m *Mongo) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	// If there are no collections to truncate, return.
	if len(req.Tables) == 0 {
		return &proto.TruncateResponse{}, nil
	}

	cs, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return nil, err
	}

	for _, collection := range req.GetTables() {
		coll := m.Client.Database(cs.Database).Collection(collection)
		_, err = coll.DeleteMany(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
	}
	return &proto.TruncateResponse{}, nil
}

// Upsert will insert or update a record in a collection.
func (m *Mongo) Upsert(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {
	records, err := tools.DecodeUpsertRecords(req)
	if err != nil {
		return nil, err
	}

	// If there are no records to upsert, return.
	if len(records) == 0 {
		return &proto.UpsertResponse{}, nil
	}

	models := []mongo.WriteModel{}
	for _, record := range records {
		doc := bson.D{}
		if err := tools.AssingRecordBSONDocument(record, &doc); err != nil {
			return nil, err
		}
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(doc).
			SetUpdate(bson.D{primitive.E{Key: "$set", Value: doc}}).
			SetUpsert(true))
	}

	cs, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return nil, err
	}

	coll := m.Client.Database(cs.Database).Collection(req.Table)
	bwr, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return nil, err
	}

	rsp := &proto.UpsertResponse{
		MatchedCount:  bwr.MatchedCount,
		UpsertedCount: bwr.UpsertedCount,
	}
	return rsp, nil
}
