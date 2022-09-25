package storage

import (
	"context"
)

// TxnChanFn is a function that will be sent to the transaction channel.
type TxnChanFn func(context.Context, Storage) error

// Txn is a wrapper for a mongo session that can be used to perform CRUD operations on a mongo DB instance.
type Txn struct {
	ch     chan TxnChanFn
	done   chan error
	commit chan bool
}

// Transactor is an interface that can be used to perform CRUD operations within the context of a database transaction.
type Transactor interface {
	Commit() error
	Rollback() error
	Send(TxnChanFn)
}

// Commit will commit the transaction.
func (txn *Txn) Commit() error {
	close(txn.ch)
	txn.commit <- true

	return <-txn.done
}

// Rollback will rollback the transaction.
func (txn *Txn) Rollback() error {
	close(txn.ch)
	txn.commit <- false

	return <-txn.done
}

// Send will send a function to the transaction channel.
func (txn *Txn) Send(fn TxnChanFn) {
	txn.ch <- fn
}
