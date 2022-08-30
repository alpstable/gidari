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
	if strings.Contains(dns, "mongodb") {
		return storage.NewMongo(ctx, dns)
	}
	if strings.Contains(dns, "postgresql") {
		return storage.NewPostgres(ctx, dns)
	}
	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}
