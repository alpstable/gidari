package storage

import (
	"context"
)

// TXChanFn is a function that will be sent to the transaction channel.
type TXChanFn func(context.Context, Storage) error

// tx is a wrapper for a mongo session that can be used to perform CRUD operations on a mongo DB instance.
type tx struct {
	ch     chan TXChanFn
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

// Send will send a function to the transaction channel.
func (tx *tx) Send(fn TXChanFn) {
	tx.ch <- fn
}
