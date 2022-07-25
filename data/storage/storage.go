package storage

import (
	"context"

	"github.com/alpine-hodler/driver/data/option"
	"github.com/alpine-hodler/driver/internal/storage"
)

// NewMongo will return a new mongo client that can be used to perform CRUD operations on a mongo DB instance. This
// constructor uses a URI to make the client connection, and the URI is of the form
// Mongo://username:password@host:port
func NewMongo(ctx context.Context, uri string, opts ...func(*option.Database)) (*storage.Mongo, error) {
	return storage.NewMongo(ctx, uri, opts...)
}

// NewPostgres will return a new Postgres option for querying data through a Postgres DB.
func NewPostgres(ctx context.Context, connectionURL string, opts ...func(*option.Database)) (*storage.Postgres, error) {
	return storage.NewPostgres(ctx, connectionURL, opts...)
}
