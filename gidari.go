// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0\n
package gidari

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alpstable/gidari/internal/transport"
)

// Config is the configuration object used to make programatic Transport requests.
type Config struct {
	transport.Config
}
// TODO #265: Remove this routine
func NewConfig(ctx context.Context, file *os.File) (*Config, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	cfg, err := transport.NewConfig(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create new config: %w", err)
	}

	// Disable logger
	cfg.Logger.SetOutput(io.Discard)

	return &Config{*cfg}, nil
}

// TransportFile will construct the transport operation using a configuration YAML file.
func TransportFile(ctx context.Context, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	cfg, err := transport.NewConfig(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to create new config: %w", err)
	}

	// Disable logger
	cfg.Logger.SetOutput(io.Discard)

	return &Config{*cfg}, nil
	if err != nil {
		return fmt.Errorf("unable to create new config: %w", err)
	}

	return Transport(ctx, cfg)
}

// Transport will construct the transport operation using a "transport.Config" object.
func Transport(ctx context.Context, cfg *Config) error {
	if err := transport.Upsert(ctx, &cfg.Config); err != nil {
		return fmt.Errorf("unable to upsert the config: %w", err)
	}

	return nil
}
