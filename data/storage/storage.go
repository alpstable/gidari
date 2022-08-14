package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/driver/data/option"
	"github.com/alpine-hodler/driver/internal/storage"
	"github.com/alpine-hodler/driver/tools"
)

// Type represents the type of storage device.
type Type uint8

const (
	Invalid Type = iota
	Mongo
	Postgres
)

// newMongo will return a new mongo client that can be used to perform CRUD operations on a mongo DB instance. This
// constructor uses a URI to make the client connection, and the URI is of the form
// Mongo://username:password@host:port
func newMongo(ctx context.Context, uri string, opts ...func(*option.Database)) (*storage.Mongo, error) {
	return storage.NewMongo(ctx, uri, opts...)
}

// newPostgres will return a new Postgres option for querying data through a Postgres DB.
func newPostgres(ctx context.Context, dns string, opts ...func(*option.Database)) (*storage.Postgres, error) {
	return storage.NewPostgres(ctx, dns, opts...)
}

// parseDNS will try to determine the storage type from a DNS string.
func parseDNS(dns string) (Type, error) {
	if strings.HasPrefix(dns, "postgresql://") {
		return Postgres, nil
	}
	if strings.HasPrefix(dns, "mongodb://") {
		return Mongo, nil
	}
	return Invalid, fmt.Errorf("invalid DNS: %v", dns)
}

// New will try to interpret which storage device to connect to using a DNS, returing a generic storage object that
// can be asserted when needed.
func New(ctx context.Context, dns string, opts ...func(*option.Database)) (tools.GenericStorage, error) {
	st, err := parseDNS(dns)
	if err != nil {
		return nil, err
	}
	switch st {
	case Mongo:
		return newMongo(ctx, dns, opts...)
	case Postgres:
		return newPostgres(ctx, dns, opts...)
	}
	return nil, fmt.Errorf("invalid DNS: %v", dns)
}
