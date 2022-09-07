package storage

import (
	"context"
)

// Tx is a wrapper for a mongo session that can be used to perform CRUD operations on a mongo DB instance.
type Tx struct {
	Ch chan func(context.Context) error

	done   chan error
	commit chan bool
}

// Commit will commit the transaction.
func (tx Tx) Commit() error {
	close(tx.Ch)
	tx.commit <- true
	return <-tx.done
}

// Rollback will rollback the transaction.
func (tx Tx) Rollback() error {
	close(tx.Ch)
	tx.commit <- false
	return <-tx.done
}
