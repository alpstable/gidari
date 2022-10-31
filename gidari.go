// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"context"
	"fmt"

	"github.com/alpstable/gidari/config"
	"github.com/alpstable/gidari/internal/transport"
)

// Transport will construct the transport operation using a "transport.Config" object.
func Transport(ctx context.Context, cfg *config.Config) error {
	if err := transport.Upsert(ctx, cfg); err != nil {
		return fmt.Errorf("unable to upsert the config: %w", err)
	}

	return nil
}
