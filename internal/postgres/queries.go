package postgres

import (
	_ "embed" // Embed external data.
)

//go:embed queries/columns.sql
var pgColumns []byte

//go:embed queries/truncate_tables.sql
var pgTruncatedTables []byte

//go:embed queries/garbage_collect.sql
var pgGarbageCollect []byte
