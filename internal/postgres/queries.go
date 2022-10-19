// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
