package storage

import (
	"context"
)

// tx is a wrapper for a mongo session that can be used to perform CRUD operations on a mongo DB instance.
type tx struct {
	ch     chan func(context.Context) error
	done   chan error
	commit chan bool
}

// Commit will commit the transaction.
func (tx *tx) Commit() error {
	close(tx.ch)
	tx.commit <- true
	return <-tx.done
}

// Rollback will rollback the transaction.
func (tx *tx) Rollback() error {
	close(tx.ch)
	tx.commit <- false
	return <-tx.done
}

// Transact will send a function to the transaction channel.
func (tx *tx) Transact(fn func(context.Context) error) {
	tx.ch <- fn
}
