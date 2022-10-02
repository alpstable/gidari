// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"golang.org/x/sync/errgroup"
)

const (
	mdbLifetime              = 60 * time.Second
	mdbTransactionRetryLimit = 3
	mdbWriteConflicErrCode   = 112
)

// Mongo is a wrapper for *mongo.Client, use to perform CRUD operations on a mongo DB instance.
type Mongo struct {
	*mongo.Client
	dns        string
	lifetime   time.Duration
	writeMutex sync.Mutex
}

// NewMongo will return a new mongo client that can be used to perform CRUD operations on a mongo DB instance. This
// constructor uses a URI to make the client connection, and the URI is of the form
// Mongo://username:password@host:port
func NewMongo(ctx context.Context, uri string) (*Mongo, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo: %w", err)
	}

	mdb := new(Mongo)
	mdb.Client = client
	mdb.dns = uri
	mdb.lifetime = mdbLifetime
	mdb.writeMutex = sync.Mutex{}

	return mdb, nil
}

// IsNoSQL returns "true" indicating that the "MongoDB" database is NoSQL.
func (m *Mongo) IsNoSQL() bool { return true }

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

// commitTransactionWithRetry will commit the transaction on the context, and retry the transaction if the commit
// fails due to a transient error.
func (m *Mongo) commitTransactionWithRetry(ctx mongo.SessionContext, retryCount int) error {
	if err := ctx.CommitTransaction(ctx); err != nil {
		// Check if the transaction error is a "mongo.ServerError".
		var mdbErr mongo.ServerError
		if retryCount <= mdbTransactionRetryLimit && errors.As(err, &mdbErr) &&
			mdbErr.HasErrorCode(mdbWriteConflicErrCode) {
			// Check if the server error is a "WriteConflict", if so then retry the transaction.
			return m.commitTransactionWithRetry(ctx, retryCount+1)
		}
	}

	return nil
}

// ReceiveWrites will listen for writes to the transaction and commit them to the database every time the lifetime
// limit is reached, or when the transaction is committed through the commit channel.
func (m *Mongo) receiveWrites(sctx mongo.SessionContext, txn *Txn) *errgroup.Group {
	lifetimeTicker := time.NewTicker(m.lifetime)
	errs, _ := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		var err error

		// Receive write requests.
		for opr := range txn.ch {
			select {
			case <-lifetimeTicker.C:
				// If the transaction has exceeded the lifetime, commit the transaction and start a new
				// one.
				if err := m.commitTransactionWithRetry(sctx, 0); err != nil {
					panic(fmt.Errorf("commit transaction: %w", err))
				}

				// Start a new transaction on the context.
				if err := sctx.StartTransaction(); err != nil {
					panic(fmt.Errorf("error starting transaction: %w", err))
				}
			default:
				if err != nil {
					continue
				}

				err = opr(sctx, m)
			}
		}

		if err != nil {
			return fmt.Errorf("error in transaction: %w", err)
		}

		return nil
	})

	return errs
}

// startSession will create a session and listen for writes, committing and reseting the transaction every 60 seconds
// to avoid lifetime limit errors.
func (m *Mongo) startSession(ctx context.Context, txn *Txn) {
	txn.done <- m.Client.UseSession(ctx, func(sctx mongo.SessionContext) error {
		// Start the transaction, if there is an error break the go routine.
		err := sctx.StartTransaction()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		errs := m.receiveWrites(sctx, txn)
		err = errs.Wait()
		if err != nil {
			return fmt.Errorf("error in transaction: %w", err)
		}

		// Await the decision to commit or rollback.
		switch {
		case <-txn.commit:
			if err := m.commitTransactionWithRetry(sctx, 0); err != nil {
				return fmt.Errorf("commit transaction: %w", err)
			}
		default:
			if err := sctx.AbortTransaction(sctx); err != nil {
				return ErrTransactionAborted
			}
		}

		return nil
	})
}

// StartTx will start a mongodb session where all data from write methods can be rolled back.
//
// MongoDB best practice is to "abort any multi-document transactions that runs for more than 60 seconds". The resulting
// error for exceeding this time constraint is "TransactionExceededLifetimeLimitSeconds". To maintain agnostism at the
// repository layer, we implement the logic to handle these transactions errors in the storage layer. Therefore, every
// 60 seconds, the transacting data will be committed commit the transaction and start a new one.
func (m *Mongo) StartTx(ctx context.Context) (*Txn, error) {
	// Construct a transaction.
	txn := &Txn{
		make(chan TxnChanFn),
		make(chan error, 1),
		make(chan bool, 1),
	}

	// Create a go routine that creates a session and listens for writes.
	go m.startSession(ctx, txn)

	return txn, nil
}

// Truncate will delete all records in a collection.
func (m *Mongo) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	// If there are no collections to truncate, return.
	if len(req.Tables) == 0 {
		return &proto.TruncateResponse{}, nil
	}

	connString, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connstring: %w", err)
	}

	for _, collection := range req.GetTables() {
		coll := m.Client.Database(connString.Database).Collection(collection)

		_, err = coll.DeleteMany(ctx, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("error truncating collection %s: %w", collection, err)
		}
	}

	return &proto.TruncateResponse{}, nil
}

// Upsert will insert or update a record in a collection.
func (m *Mongo) Upsert(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()

	records, err := tools.DecodeUpsertRecords(req)
	if err != nil {
		return nil, fmt.Errorf("failed to decode records: %w", err)
	}

	// If there are no records to upsert, return.
	if len(records) == 0 {
		return &proto.UpsertResponse{}, nil
	}

	models := []mongo.WriteModel{}

	for _, record := range records {
		doc := bson.D{}
		if err := tools.AssingRecordBSONDocument(record, &doc); err != nil {
			return nil, fmt.Errorf("failed to assign record to bson document: %w", err)
		}

		models = append(models, mongo.NewUpdateOneModel().SetFilter(doc).
			SetUpdate(bson.D{primitive.E{Key: "$set", Value: doc}}).
			SetUpsert(true))
	}

	cs, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	coll := m.Client.Database(cs.Database).Collection(req.Table)

	bwr, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return nil, fmt.Errorf("bulk write error: %w", err)
	}

	return &proto.UpsertResponse{MatchedCount: bwr.MatchedCount, UpsertedCount: bwr.UpsertedCount}, nil
}

// ListPrimaryKeys will return a "proto.ListPrimaryKeysResponse" containing a list of primary keys data for all tables
// in a database. MongoDB does not have a concept of primary keys, so we will return the "_id" field as the primary key
// for all collections in the database associated with the underlying connection string.
func (m *Mongo) ListPrimaryKeys(ctx context.Context) (*proto.ListPrimaryKeysResponse, error) {
	collections, err := m.ListTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing collections: %w", err)
	}

	rsp := &proto.ListPrimaryKeysResponse{PKSet: make(map[string]*proto.PrimaryKeys)}

	for collection := range collections.GetTableSet() {
		if rsp.PKSet[collection] == nil {
			rsp.PKSet[collection] = &proto.PrimaryKeys{}
		}

		rsp.PKSet[collection].List = append(rsp.PKSet[collection].List, "_id")
	}

	return rsp, nil
}

// ListTables will return a list of all tables in the MongoDB database.
func (m *Mongo) ListTables(ctx context.Context) (*proto.ListTablesResponse, error) {
	connString, err := connstring.ParseAndValidate(m.dns)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	collections, err := m.Client.Database(connString.Database).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	rsp := &proto.ListTablesResponse{TableSet: make(map[string]*proto.Table)}

	for _, collection := range collections {
		// Need to get the size of the collection
		result, err := m.Client.Database(connString.Database).RunCommand(ctx, bson.D{
			primitive.E{Key: "collStats", Value: collection},
		}).DecodeBytes()
		if err != nil {
			return nil, fmt.Errorf("failed to get collection stats: %w", err)
		}

		rawValue, err := result.LookupErr("size")
		if err != nil {
			return nil, fmt.Errorf("failed to get collection size: %w", err)
		}

		rsp.TableSet[collection] = &proto.Table{Size: int64(rawValue.Int32())}
	}

	return rsp, nil
}
