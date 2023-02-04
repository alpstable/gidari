// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"context"
)

type UpsertWriter interface {
	// Upsert will use an UpsertRequest to upsert a new or existing
	// object into the storage backend.
	Write(context.Context, *UpsertRequest) error
}
