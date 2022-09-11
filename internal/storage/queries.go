package storage

import (
	_ "embed"
)

//go:embed queries/pg_columns.sql
var pgColumns []byte

//go:embed queries/pg_tables.sql
var pgTables []byte

//go:embed queries/pg_truncate_tables.sql
var pgTruncatedTables []byte
