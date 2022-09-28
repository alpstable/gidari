package storage

import (
	_ "embed" // Embed external data.
)

//go:embed queries/pg_columns.sql
var pgColumns []byte

//go:embed queries/pg_truncate_tables.sql
var pgTruncatedTables []byte
