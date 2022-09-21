package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // postgres driver
)

const (
	postgresPartitionSize = 1000
)

// postgresTxType is a type alias for the postgres transaction type.
type postgresTxType uint8

const (
	basicPostgressTxID postgresTxType = iota
)

type pgmeta struct {
	table   string
	pk      []string
	columns []string
}

func (meta *pgmeta) addColumn(pk string) {
	meta.columns = append(meta.columns, pk)
}

func (meta *pgmeta) addPK(pk string) {
	meta.pk = append(meta.pk, pk)
}

func (meta *pgmeta) isPK(name string) bool {
	for _, pk := range meta.pk {
		if pk == name {
			return true
		}
	}
	return false
}

// getMeta will get postgres database metadata for processing generalized functionality, such as upsert.
func (pg *Postgres) getMeta(ctx context.Context, table string) (*pgmeta, error) {
	if len(pg.meta) == 0 {
		columns, err := pg.ListColumns(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting postgres metadata: %v", err)
		}

		pg.meta = make(map[string]*pgmeta)
		for _, record := range columns.Records {
			table := record.AsMap()["table_name"].(string)

			// Initialize the table pgmeta if it does not exist.
			if pg.meta[table] == nil {
				pg.meta[table] = &pgmeta{table: table}
			}

			// Add PK and general column data to the pgmeta table object.
			meta := pg.meta[table]
			columnName := record.AsMap()["column_name"].(string)
			if record.AsMap()["primary_key"].(float64) == 1.0 {
				meta.addPK(columnName)
			}
			meta.addColumn(columnName)
		}
	}

	meta := pg.meta[table]
	if meta == nil {
		return nil, fmt.Errorf("table doesn't exist %q", table)
	}
	return meta, nil
}

// exec executes a query that requires no input, passing the resulting rows into a user-defined teardown
// function.
func (pg *Postgres) exec(ctx context.Context, query []byte, teardown func(*sql.Rows) error) error {
	stmt, err := pg.DB.PrepareContext(ctx, string(query))
	if err != nil {
		return fmt.Errorf("unable to prepare statement: %v", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to query: %v", err)
	}
	defer rows.Close()
	return teardown(rows)
}

// Close will close the underlying database / transaction.
func (pg *Postgres) Close() {
	if pg.DB != nil {
		pg.DB.Close()
	}
}

// ListColumns will set a complete list of available columns per table on the response.
func (pg *Postgres) ListColumns(ctx context.Context) (*proto.ListColumnsResponse, error) {
	stmt, err := pg.DB.PrepareContext(ctx, string(pgColumns))
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to query: %v", err)
	}
	defer rows.Close()

	var rsp proto.ListColumnsResponse
	err = tools.AssignStructs(rows, &rsp.Records)
	if err != nil {
		return nil, fmt.Errorf("unable to assign structs: %v", err)
	}
	return &rsp, nil
}

// ListTables will set a complete list of available tables on the response.
func (pg *Postgres) ListTables(ctx context.Context) (*proto.ListTablesResponse, error) {
	stmt, err := pg.DB.PrepareContext(ctx, string(pgTables))
	if err != nil {
		return nil, fmt.Errorf("unable to prepare statement: %v", err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to query: %v", err)
	}
	defer rows.Close()

	var rsp proto.ListTablesResponse
	err = tools.AssignStructs(rows, &rsp.Records)
	if err != nil {
		return nil, fmt.Errorf("unable to assign structs: %v", err)
	}
	return &rsp, nil
}

// Read read will attempt to assign a reader buidler based on the request, assinging the resuling rows to the response
// in-memory.
func (pg *Postgres) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	// bldr, err := query.GetReadBuilder(query.ReadBuilderType(req.ReaderBuilder[0]))
	// if err != nil {
	// 	return err
	// }

	// query, err := bldr.ReaderQuery(query.PostgresStorage)
	// if err != nil {
	// 	return err
	// }
	// stmt, err := pg.PrepareContext(ctx, string(query))
	// if err != nil {
	// 	return err
	// }

	// args, err := bldr.ReaderArgs(req)
	// if err != nil {
	// 	return err
	// }
	// rows, err := stmt.QueryContext(ctx, args...)
	// if err != nil {
	// 	return err
	// }

	// if err := tools.AssignStructs(rows, &rsp.Records); err != nil {
	// 	return err
	// }
	return nil
}

// Truncate will truncate a table.
func (pg *Postgres) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	// If the table is not specified, return an error.
	if len(req.Tables) == 0 {
		return &proto.TruncateResponse{}, nil
	}

	tables := req.GetTables()
	if len(tables) == 0 {
		return nil, fmt.Errorf("no tables specified")
	}
	query := fmt.Sprintf(string(pgTruncatedTables), strings.Join(tables, ","))
	return &proto.TruncateResponse{}, pg.exec(ctx, []byte(query), func(r *sql.Rows) error { return nil })
}

// Upsert will insert the records on the request if they do not exist in the database. On conflict, it will use the
// PK on the request record to update the data in the database. An upsert request will update the entire table
// for a given record, include fields that have not been set directly.
func (pg *Postgres) Upsert(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {
	records, err := tools.DecodeUpsertRecords(req)
	if err != nil {
		return nil, fmt.Errorf("unable to decode records: %v", err)
	}

	// Do nothing if there are no records.
	if len(records) == 0 {
		return &proto.UpsertResponse{}, nil
	}

	meta, err := pg.getMeta(ctx, req.Table)
	if err != nil {
		return nil, fmt.Errorf("table %q does not exist", req.Table)
	}

	exclusTemplate := "\"%s\" = EXCLUDED.\"%s\""

	exclusions := []string{}
	columns := []string{}
	for _, name := range meta.columns {
		if !meta.isPK(name) {
			exclusions = append(exclusions, fmt.Sprintf(exclusTemplate, name, name))
		}
		columns = append(columns, name)
	}

	pkstr := strings.Join(meta.pk, ",")
	exstr := strings.Join(exclusions, ",")
	clstr := strings.Join(columns, ",")

	upsertQuery := `INSERT INTO %s(%s) VALUES %s ON CONFLICT (%s) DO UPDATE SET %s`
	upsertTemplate := fmt.Sprintf(upsertQuery, meta.table, clstr, "%s", pkstr, exstr)

	// Get record-specific metadata from a sample record.
	sample := records[0]
	recordSize := len(sample.Fields)

	// Upsert 1000 records at a time.
	for _, partition := range tools.PartitionStructs(postgresPartitionSize, records) {
		volume := len(partition)
		placeholders := make([]string, 0, volume)
		arguments := make([]interface{}, 0, recordSize*volume)

		// Prepare data to populate the prepared statement.
		for idx, record := range partition {
			ph := []string{}
			for i := 1; i <= len(meta.columns); i++ {
				ph = append(ph, fmt.Sprintf("$%d", recordSize*idx+i))
			}
			recordph := fmt.Sprintf("(%s)", strings.Join(ph, ","))
			placeholders = append(placeholders, recordph)

			mrecord := record.AsMap()
			for _, col := range meta.columns {
				arguments = append(arguments, mrecord[col])
			}
		}

		// Create the prepared statement for execution.
		query := fmt.Sprintf(upsertTemplate, strings.Join(placeholders, ","))

		var stmt *sql.Stmt

		// If the context has a transaction ID set get it and then lookup the transaction from the activxTx map.
		if txID, ok := ctx.Value(basicPostgressTxID).(string); ok {
			if probablyTX, ok := pg.activeTx.Load(txID); ok {
				tx := probablyTX.(*sql.Tx)
				stmt, err = tx.PrepareContext(ctx, query)
				if err != nil {
					return nil, fmt.Errorf("unable to prepare statement: %v", err)
				}
			}
		} else {
			stmt, err = pg.DB.PrepareContext(ctx, query)
			if err != nil {
				return nil, fmt.Errorf("unable to prepare statement: %v", err)
			}
		}

		// Execute upsert.
		if _, err := stmt.ExecContext(ctx, arguments...); err != nil {
			return nil, fmt.Errorf("unable to execute upsert: %v", err)
		}
	}
	return &proto.UpsertResponse{}, nil
}

// Postgres is a wrapper around the sql.DB object.
type Postgres struct {
	*sql.DB
	meta map[string]*pgmeta

	// activeTx are the transactions that are currently active on this connection. When a user calls "StartTx" on
	// a Postgres intance, a transaction is created and added to this map. Afterward, if the user calls a write
	// method (e.g. Insert, Update, Delete, Upsert), the transaction will be used to execute the query. In order for
	// the write method to know which transaction to use, a context with the transaction ID must be passed into
	// the method. The transaction ID is added to the context in the "StartTx" method. The transaction ID is
	// removed from the context in the "CommitTx" and "RollbackTx" methods.
	activeTx sync.Map
}

// NewPostgres will return a new Postgres option for querying data through a Postgres DB.
func NewPostgres(ctx context.Context, connectionURL string) (*Postgres, error) {
	postgres := new(Postgres)

	var err error
	postgres.DB, err = sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to postgres: %v", err)
	}

	postgres.setMaxOpenConns()
	postgres.meta = make(map[string]*pgmeta)
	postgres.activeTx = sync.Map{}

	return postgres, nil
}

// Type implements the storage interface.
func (pg *Postgres) Type() uint8 { return PostgresType }

// pgMaxConnectionsUpperLimit will return the most ideal upper limit for the maximum number of connections for a
// Postgres DB. https://tinyurl.com/57kyjtwd
func (pg *Postgres) setMaxOpenConns() {

	// num_cores is the number of cores available
	numCores := runtime.NumCPU()

	// parallel_io_limit is the number of concurrent I/O requests your storage subsystem can handle.
	var parallelIOLimit = 115
	// if pg.opts != nil && pg.opts.ParallelIOLimit != nil {
	// 	parallelIOLimit = *pg.opts.ParallelIOLimit
	// } else {
	// At provision, Databases for PostgreSQL sets the maximum number of connections to your PostgreSQL database to 115.
	// 15 connections are reserved for the superuser to maintain the state and integrity of your database, and 100
	// connections are available for you and your applications. https://tinyurl.com/3yyu6eaf
	// }

	// session_busy_ratio is the fraction of time that the connection is active executing a statement in the database.
	// If your workload consists of big analytical queries, session_busy_ratio can be up to 1.
	var sessionBusyRatio = 1.0
	// if pg.opts != nil && pg.opts.SessionBusyRatio != nil {
	// 	sessionBusyRatio = *pg.opts.SessionBusyRatio
	// } else {
	// These queries for this db are typically expected to be 1-1 with upserting and finding records definied in the
	// alpine-hodler/web API. That is, they should be extremely simple and devoid of business logic, and so the default
	// value for this ratio is 1.
	// }

	var avgParallelism = 1.0
	// if pg.opts != nil && pg.opts.AvgParallelism != nil {
	// 	avgParallelism = *pg.opts.AvgParallelism
	// } else {
	// These queries for this db are typically expected to be 1-1 with upserting and finding records definied in the
	// alpine-hodler/web API. That is, they should be extremely simple and devoid of business logic, and so the default
	// value for this average is one. We expect that the average number of backend processes working on a SINGLE query
	// to be 1.
	// }
	n := (math.Max(float64(numCores), float64(parallelIOLimit))) / (sessionBusyRatio * avgParallelism)
	pg.DB.SetMaxOpenConns(int(n))
}

// StartTx will start a transaction on the Postgres connection. The transaction ID is returned and should be used
// to commit or rollback the transaction.
func (pg *Postgres) StartTx(ctx context.Context) (Tx, error) {
	// Construct a gidari storage transaction.
	txn := &tx{
		make(chan TXChanFn),
		make(chan error, 1),
		make(chan bool, 1),
	}

	// Instantiate a new transaction on the Postgres connection and store it in the activeTx map.
	txnID := uuid.New().String()
	pgtx, err := pg.DB.BeginTx(ctx, nil)
	if err != nil {
		return txn, fmt.Errorf("failed to start transaction: %w", err)
	}
	pg.activeTx.Store(txnID, pgtx)

	// Add the transaction ID to the context.
	pgCtx := context.WithValue(ctx, basicPostgressTxID, txnID)

	go func() {
		defer func() {
			// Remove the transaction from the activeTx map.
			pg.activeTx.Delete(txnID)
		}()

		for fn := range txn.ch {
			if err != nil {
				continue
			}
			err = fn(pgCtx, pg)
		}
		if err != nil {
			txn.done <- err
			return
		}

		// Await the decision to commit or rollback.
		select {
		case <-txn.commit:
			txn.done <- pgtx.Commit()
		default:
			txn.done <- pgtx.Rollback()
		}
	}()
	return txn, nil
}
