package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDNS(t *testing.T) {
	t.Run("mongodb", func(t *testing.T) {
		dns := "mongodb://mongo-coinbasepro:27017/coinbasepro"
		storagetype, err := parseDNS(dns)
		require.NoError(t, err)
		require.Equal(t, Mongo, storagetype)
	})

	t.Run("postgres", func(t *testing.T) {
		dns := "postgresql://postgres-coinbasepro:5432/coinbasepro?sslmode=disable"
		storagetype, err := parseDNS(dns)
		require.NoError(t, err)
		require.Equal(t, Postgres, storagetype)
	})
}
