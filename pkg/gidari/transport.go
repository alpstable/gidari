package gidari

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alpstable/gidari/internal/transport"
	"github.com/alpstable/gidari/pkg/config"
)

// Transport will construct the transport operation using a "transport.Config" object.
func Transport(ctx context.Context, cfg *config.Config) error {
	if err := transport.Upsert(ctx, cfg); err != nil {
		return fmt.Errorf("unable to upsert the config: %w", err)
	}

	return nil
}

// TransportFile will construct the transport operation using a configuration YAML file.
func TransportFile(ctx context.Context, file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	cfg, err := config.New(bytes)
	if err != nil {
		return fmt.Errorf("unable to create new config: %w", err)
	}

	// Disable logger
	cfg.Logger.SetOutput(io.Discard)

	if err != nil {
		return fmt.Errorf("unable to create new config: %w", err)
	}

	return Transport(ctx, cfg)
}
