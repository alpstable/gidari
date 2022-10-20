// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/alpstable/gidari/internal/proto"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultPartitionSize               = 1000
	defaultGarbageCollectionRetryLimit = 10
)

var (
	ErrTransactionNotFound   = fmt.Errorf("transaction not found")
	ErrNoTables              = fmt.Errorf("no tables found")
	ErrUnsupportedDataType   = fmt.Errorf("unsupported data type")
	ErrFailedToMarshalJSON   = fmt.Errorf("failed to marshal json")
	ErrFailedToUnmarshalJSON = fmt.Errorf("failed to unmarshal json")
)

// postgresTxType is a type alias for the postgres transaction type.
type postgresTxType uint8

const (
	basicPostgressTxID postgresTxType = iota
)

// sqlPrepareContextFn can be used to prepare a statement and return the result.
type sqlPrepareContextFn func(context.Context, string) (*sql.Stmt, error)

type pgmeta struct {
	// cols are the columns for a specific table.
	cols map[string][]string

	// pks are the primary keys for a specific table.
	pks map[string][]string

	// bytes are the size in bytes for a specific table.
	bytes map[string]int64
}

func (meta *pgmeta) isPK(table, name string) bool {
	for _, pk := range meta.pks[table] {
		if pk == name {
			return true
		}
	}

	return false
}

// fortmatPlaceholders will return a string of placeholders that require a number next to the placeholder string
// that iteratively increases by the number of arguments passed to the query. For example, if a string has numCols=3
// and numRows=2, this function will return "(?1,?2,?3),(?4,?5,?6)".
func formatPlaceholders(numCols int, numRows int, symbol string) string {
	if numCols == 0 || numRows == 0 {
		return "()"
	}

	if symbol == "" {
		symbol = "?"
	}

	var strBldr strings.Builder

	for pos := 0; pos < numRows*numCols; pos++ {
		if pos%numCols == 0 {
			strBldr.WriteString("(")
		}

		strBldr.WriteString(symbol)
		strBldr.WriteString(strconv.Itoa(pos + 1))

		if pos%numCols == numCols-1 {
			strBldr.WriteString(")")

			if pos != numRows*numCols-1 {
				strBldr.WriteString(",")
			}
		} else {
			strBldr.WriteString(",")
		}
	}

	return strBldr.String()
}

// flattenPartition will take a slice of structures, extract data from their fields, and append it to a slice.
// This will "flatten" the data to be used in conjunctino with placeholders in a SQL query.
func flattenPartition(columns []string, partition []*structpb.Struct) []interface{} {
	var args []interface{}

	for _, record := range partition {
		hash := record.AsMap()
		for _, column := range columns {
			args = append(args, hash[column])
		}
	}

	return args
}

// exclusionConstraints will return a string of non-primary key columns to "exclude" if they are not changed in the
// context of a Postgres insert. That is, if a column is not changed, it will not be updated. All columns beside primary
// keys must be included in the "excluded" clause.
func (meta *pgmeta) exclusionConstraints(table string) []string {
	var constraints []string

	for _, column := range meta.cols[table] {
		if !meta.isPK(table, column) {
			constraints = append(constraints, fmt.Sprintf("\"%s\" = EXCLUDED.\"%s\"", column, column))
		}
	}

	return constraints
}

// upsertStatement will return a postgres upsert statement for the meta object.
func (meta *pgmeta) upsertStmt(ctx context.Context, table string, pcf sqlPrepareContextFn, vol int) (*sql.Stmt, error) {
	query := fmt.Sprintf(`INSERT INTO %s(%s) VALUES %s ON CONFLICT (%s) DO UPDATE SET %s`, table,
		strings.Join(meta.cols[table], ","),
		formatPlaceholders(len(meta.cols[table]), vol, "$"),
		strings.Join(meta.pks[table], ","),
		strings.Join(meta.exclusionConstraints(table), ","))

	stmt, err := pcf(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %w", err)
	}

	return stmt, nil
}

// garbageCollect will garbage collect the database. This will return disk space to the OS by running `VACUUM FULL`.
// For more information, see: https://www.postgresql.org/docs/current/sql-vacuum.html
func (pg *Postgres) garbageCollect(ctx context.Context, retryCount uint8, tables ...string) error {
	query := fmt.Sprintf(string(pgGarbageCollect), strings.Join(tables, ","))

	stmt, err := pg.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}

	// Execute the garbage collection query.
	if _, err := stmt.ExecContext(ctx); err != nil {
		// If the garbage collection fails due to a deadlock, we will retry the operation. We should not
		// retry more than a deterministic number of times, defined by "pgGCRetryLimit".
		var pqErr *pq.Error
		if retryCount <= defaultGarbageCollectionRetryLimit && errors.As(err, &pqErr) && pqErr.Code == "40P01" {
			return pg.garbageCollect(ctx, retryCount+1, tables...)
		}

		return fmt.Errorf("unable to execute statement: %w", err)
	}

	return nil
}

// loadMeta will load the postgres metadata for the database. If the data has already been loaded, this method will do
// nothing.
func (pg *Postgres) loadMeta(ctx context.Context, garbageC bool) error {
	pg.metaMutex.Lock()
	defer pg.metaMutex.Unlock()

	if garbageC {
		// Need to garbage collected the database before making this query.
		if err := pg.garbageCollect(ctx, 0); err != nil {
			return fmt.Errorf("unable to garbage collect: %w", err)
		}
	}

	stmt, err := pg.DB.PrepareContext(ctx, string(pgColumns))
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %w", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	pg.meta.cols = make(map[string][]string)
	pg.meta.pks = make(map[string][]string)
	pg.meta.bytes = make(map[string]int64)

	for rows.Next() {
		var (
			table      string
			column     string
			primaryKey bool
			bytes      int64
		)

		if err := rows.Scan(&column, &table, &primaryKey, &bytes); err != nil {
			return fmt.Errorf("unable to scan row: %w", err)
		}

		if primaryKey {
			pg.meta.pks[table] = append(pg.meta.pks[table], column)
		}

		pg.meta.cols[table] = append(pg.meta.cols[table], column)
		pg.meta.bytes[table] = bytes
	}

	return nil
}

// Close will close the underlying database / transaction.
func (pg *Postgres) Close() {
	if pg.DB != nil {
		pg.DB.Close()
	}
}

// ListColumns will set a complete list of available columns per table on the response.
func (pg *Postgres) ListColumns(ctx context.Context) (*proto.ListColumnsResponse, error) {
	if err := pg.loadMeta(ctx, false); err != nil {
		return nil, fmt.Errorf("unable to load postgres metadata: %w", err)
	}

	var rsp proto.ListColumnsResponse
	for table, columns := range pg.meta.cols {
		rsp.ColSet[table].List = append(rsp.ColSet[table].List, columns...)
	}

	return &rsp, nil
}

// ListPrimaryKeys will list all primary keys for all of the tables in the database defined by the DNS used to create
// the postgres instance.
func (pg *Postgres) ListPrimaryKeys(ctx context.Context) (*proto.ListPrimaryKeysResponse, error) {
	if err := pg.loadMeta(ctx, false); err != nil {
		return nil, fmt.Errorf("unable to load postgres metadata: %w", err)
	}

	rsp := &proto.ListPrimaryKeysResponse{PKSet: make(map[string]*proto.PrimaryKeys)}

	for table, pks := range pg.meta.pks {
		if rsp.PKSet[table] == nil {
			rsp.PKSet[table] = &proto.PrimaryKeys{}
		}

		rsp.PKSet[table].List = append(rsp.PKSet[table].List, pks...)
	}

	return rsp, nil
}

// ListTables will set a complete list of available tables on the response.
func (pg *Postgres) ListTables(ctx context.Context) (*proto.ListTablesResponse, error) {
	// Since tables have a "size" associated with them, we need to garbage collect the database before we can
	// get a complete list of tables.
	if err := pg.loadMeta(ctx, true); err != nil {
		return nil, fmt.Errorf("unable to load postgres metadata: %w", err)
	}

	rsp := &proto.ListTablesResponse{TableSet: make(map[string]*proto.Table)}

	for table := range pg.meta.cols {
		rsp.TableSet[table] = &proto.Table{Size: pg.meta.bytes[table]}
	}

	return rsp, nil
}

// Truncate will truncate a table.
func (pg *Postgres) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	// If the table is not specified, return an error.
	if len(req.Tables) == 0 {
		return &proto.TruncateResponse{}, nil
	}

	tables := req.GetTables()
	if len(tables) == 0 {
		return nil, ErrNoTables
	}

	stmt, err := pg.DB.PrepareContext(ctx, fmt.Sprintf(string(pgTruncatedTables), strings.Join(tables, ",")))
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %w", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	return &proto.TruncateResponse{}, nil
}

// getPrepareContextFn will return a function that can prepare an upsert statement for a given table.
func (pg *Postgres) getPrepareContextFn(ctx context.Context) (sqlPrepareContextFn, error) {
	// First check to see if a transaction has been assigned to the context. If it has, use the transaction.
	// Otherwise, use the database.
	txID, ok := ctx.Value(basicPostgressTxID).(string)
	if ok {
		if tx, ok := pg.activeTx.Load(txID); ok {
			tx, ok := tx.(*sql.Tx)
			if !ok {
				return nil, ErrTransactionNotFound
			}

			return tx.PrepareContext, nil
		}
	}

	return pg.DB.PrepareContext, nil
}

func (pg *Postgres) upsert(ctx context.Context, table string, records []*structpb.Struct) error {
	prepareContextFn, err := pg.getPrepareContextFn(ctx)
	if err != nil {
		return fmt.Errorf("unable to get preparer: %w", err)
	}

	if err := pg.loadMeta(ctx, false); err != nil {
		return fmt.Errorf("unable to load postgres metadata: %w", err)
	}

	// Upsert 1000 records at a time, the maximum number of records that can be inserted in a single statement on a
	// postgres database.
	for _, partition := range proto.PartitionStructs(defaultPartitionSize, records) {
		stmt, err := pg.meta.upsertStmt(ctx, table, prepareContextFn, len(partition))
		if err != nil {
			return fmt.Errorf("unable to prepare statement: %w", err)
		}

		// Execute upsert.
		arguments := flattenPartition(pg.meta.cols[table], partition)
		if _, err := stmt.ExecContext(ctx, arguments...); err != nil {
			return fmt.Errorf("unable to execute upsert: %w", err)
		}
	}

	return nil
}

// Upsert will insert the records on the request if they do not exist in the database. On conflict, it will use the
// PK on the request record to update the data in the database. An upsert request will update the entire table
// for a given record, include fields that have not been set directly.
func (pg *Postgres) Upsert(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {
	pg.writeMutex.Lock()
	defer pg.writeMutex.Unlock()

	records, err := proto.DecodeUpsertRequest(req)
	if err != nil {
		return nil, fmt.Errorf("unable to decode records: %w", err)
	}

	// Do nothing if there are no records.
	if len(records) == 0 {
		return &proto.UpsertResponse{}, nil
	}

	table := req.GetTable()
	if err := pg.upsert(ctx, table, records); err != nil {
		return nil, fmt.Errorf("unable to upsert: %w", err)
	}

	return &proto.UpsertResponse{}, nil
}

// Postgres is a wrapper around the sql.DB object.
type Postgres struct {
	*sql.DB

	// meta hold metdata about the database.
	meta *pgmeta

	metaMutex  sync.Mutex
	writeMutex sync.Mutex

	// activeTx are the transactions that are currently active on this connection. When a user calls "StartTx" on
	// a Postgres intance, a transaction is created and added to this map. Afterward, if the user calls a write
	// method (e.g. Insert, Update, Delete, Upsert), the transaction will be used to execute the query. In order for
	// the write method to know which transaction to use, a context with the transaction ID must be passed into
	// the method. The transaction ID is added to the context in the "StartTx" method. The transaction ID is
	// removed from the context in the "CommitTx" and "RollbackTx" methods.
	activeTx sync.Map
}

// New will return a new Postgres option for querying data through a Postgres DB.
func New(ctx context.Context, connectionURL string) (*Postgres, error) {
	postgres := new(Postgres)

	var err error

	postgres.DB, err = sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %w", err)
	}

	postgres.setMaxOpenConns()
	postgres.meta = new(pgmeta)
	postgres.metaMutex = sync.Mutex{}
	postgres.writeMutex = sync.Mutex{}
	postgres.activeTx = sync.Map{}

	return postgres, nil
}

// IsNoSQL returns "false" to indicate that "Postgres" is not a NoSQL database.
func (pg *Postgres) IsNoSQL() bool { return false }

// Type implements the storage interface.
func (pg *Postgres) Type() uint8 { return proto.PostgresType }

// pgMaxConnectionsUpperLimit will return the most ideal upper limit for the maximum number of connections for a
// Postgres DB. https://tinyurl.com/57kyjtwd
func (pg *Postgres) setMaxOpenConns() {
	// num_cores is the number of cores available
	numCores := runtime.NumCPU()

	// parallel_io_limit is the number of concurrent I/O requests your storage subsystem can handle.
	parallelIOLimit := 115
	// if pg.opts != nil && pg.opts.ParallelIOLimit != nil {
	// 	parallelIOLimit = *pg.opts.ParallelIOLimit
	// } else {
	// At provision, Databases for PostgreSQL sets the maximum number of connections to your PostgreSQL database to 115.
	// 15 connections are reserved for the superuser to maintain the state and integrity of your database, and 100
	// connections are available for you and your applications. https://tinyurl.com/3yyu6eaf
	// }

	// session_busy_ratio is the fraction of time that the connection is active executing a statement in the database.
	// If your workload consists of big analytical queries, session_busy_ratio can be up to 1.
	sessionBusyRatio := 1.0
	// if pg.opts != nil && pg.opts.SessionBusyRatio != nil {
	// 	sessionBusyRatio = *pg.opts.SessionBusyRatio
	// } else {
	// These queries for this db are typically expected to be 1-1 with upserting and finding records definied in the
	// alpstable/web API. That is, they should be extremely simple and devoid of business logic, and so the default
	// value for this ratio is 1.
	// }

	avgParallelism := 1.0
	// if pg.opts != nil && pg.opts.AvgParallelism != nil {
	// 	avgParallelism = *pg.opts.AvgParallelism
	// } else {
	// These queries for this db are typically expected to be 1-1 with upserting and finding records definied in the
	// alpstable/web API. That is, they should be extremely simple and devoid of business logic, and so the default
	// value for this average is one. We expect that the average number of backend processes working on a SINGLE query
	// to be 1.
	// }
	n := (math.Max(float64(numCores), float64(parallelIOLimit))) / (sessionBusyRatio * avgParallelism)
	pg.DB.SetMaxOpenConns(int(n))
}

// StartTx will start a transaction on the Postgres connection. The transaction ID is returned and should be used
// to commit or rollback the transaction.
func (pg *Postgres) StartTx(ctx context.Context) (*proto.Txn, error) {
	// Construct a gidari storage transaction.
	txn := &proto.Txn{
		FunctionCh: make(chan proto.TxnChanFn),
		DoneCh:     make(chan error, 1),
		CommitCh:   make(chan bool, 1),
	}

	// Instantiate a new transaction on the Postgres connection and store it in the activeTx map.
	txnID := uuid.New().String()

	pgtx, err := pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return txn, fmt.Errorf("failed to start transaction: %w", err)
	}

	pg.activeTx.Store(txnID, pgtx)

	// Create a copy of the parent context with a transaction ID.
	pgCtx := context.WithValue(ctx, basicPostgressTxID, txnID)

	go func() {
		defer func() {
			// Remove the transaction from the activeTx map.
			pg.activeTx.Delete(txnID)
		}()

		for fn := range txn.FunctionCh {
			if err != nil {
				continue
			}

			err = fn(pgCtx, pg)
		}

		if err != nil {
			txn.DoneCh <- err

			return
		}

		if <-txn.CommitCh {
			txn.DoneCh <- pgtx.Commit()
		} else {
			txn.DoneCh <- pgtx.Rollback()
		}
	}()

	return txn, nil
}

// UpsertBinary will upsert binary data into the PostgresDB for a "property bag"-like table, storing the data in a
// binary-type column.
func (pg *Postgres) UpsertBinary(ctx context.Context,
	req *proto.UpsertBinaryRequest,
) (*proto.UpsertBinaryResponse, error) {
	pg.writeMutex.Lock()
	defer pg.writeMutex.Unlock()

	records, err := proto.DecodeUpsertBinaryRequest(req)
	if err != nil {
		return nil, fmt.Errorf("unable to decode records: %w", err)
	}

	// Do nothing if there are no records.
	if len(records) == 0 {
		return &proto.UpsertBinaryResponse{}, nil
	}

	table := req.GetTable()
	if err := pg.upsert(ctx, table, records); err != nil {
		return nil, fmt.Errorf("unable to upsert: %w", err)
	}

	return &proto.UpsertBinaryResponse{}, nil
}

// Ping will return an error if the connection to the DB is lost
func (pg *Postgres) Ping() error {
	if err := pg.Ping(); err != nil {
		return fmt.Errorf("Connection lost, error: %v", err)
	}
	return nil
}
