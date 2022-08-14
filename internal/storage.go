package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/alpine-hodler/driver/internal/storage"
	"github.com/alpine-hodler/driver/tools"
)

// NewStorage will attempt to return a generic storage object given a DNS.
func NewStorage(ctx context.Context, dns string) (tools.GenericStorage, error) {
	if strings.Contains(dns, "mongodb") {
		return storage.NewMongo(ctx, dns)
	}
	if strings.Contains(dns, "postgresql") {
		return storage.NewPostgres(ctx, dns)
	}
	return nil, fmt.Errorf("databse for dns %q is not supported", dns)
}
