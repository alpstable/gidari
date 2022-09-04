package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/sherpa/internal/storage"
)

// S is an implementaiton of a generic storage, any storage device should be implementable on S.
type S storage.Storage

// New will attempt to return a generic storage object given a DNS.
func New(ctx context.Context, dns string) (S, error) {
	mongoTypeStr := storage.DNSRoot(storage.MongoType)
	if strings.Contains(dns, mongoTypeStr) {
		return storage.NewMongo(ctx, dns)
	}

	postgresqlTypeStr := storage.DNSRoot(storage.PostgresType)
	if strings.Contains(dns, postgresqlTypeStr) {
		return storage.NewPostgres(ctx, dns)
	}

	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}
