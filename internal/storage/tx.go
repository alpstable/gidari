package storage

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Tx is a wrapper for a mongo session that can be used to perform CRUD operations on a mongo DB instance.
type Tx struct {
	Errs *errgroup.Group
	Ch   chan func(context.Context) error

	done     chan bool
	commit   chan bool
	rollback chan bool
}

// Commit will commit the transaction.
func (tx Tx) Commit() {
	close(tx.Ch)
	tx.commit <- true
	tx.rollback <- false
	<-tx.done
}

func (tx Tx) Rollback() {
	close(tx.Ch)
	tx.commit <- false
	tx.rollback <- true
	<-tx.done
}
