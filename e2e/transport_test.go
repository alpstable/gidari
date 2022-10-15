// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/config"
	"github.com/sirupsen/logrus"
)

func TestUpsert(t *testing.T) {
	t.Parallel()

	// Iterate over the fixtures/upsert directory and run each configuration file.
	fixtureRoot := "testdata/upsert"

	fixtures, err := os.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatalf("error reading fixtures: %v", err)
	}

	for _, fixture := range fixtures {
		name := fixture.Name()
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			path := filepath.Join(fixtureRoot, name)

			file, err := os.Open(path)
			if err != nil {
				t.Fatalf("error opening fixture: %v", err)
			}

			cfg, err := config.New(ctx, file)
			if err != nil {
				t.Fatalf("error creating config: %v", err)
			}

			cfg.Logger = logrus.New()

			// Fill in the authentication details for the fixture.
			cfgAuth := cfg.Authentication
			if cfgAuth.APIKey != nil {
				// The "passhprase" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = config.Authentication{
					APIKey: &config.APIKey{
						Key:        os.Getenv(cfgAuth.APIKey.Key),
						Secret:     os.Getenv(cfgAuth.APIKey.Secret),
						Passphrase: os.Getenv(cfgAuth.APIKey.Passphrase),
					},
				}
			}

			if cfgAuth.Auth2 != nil {
				// The "bearer" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = config.Authentication{
					Auth2: &config.Auth2{
						Bearer: os.Getenv(cfgAuth.Auth2.Bearer),
					},
				}
			}

			// Upsert the fixture.
			if err := gidari.Transport(context.Background(), cfg); err != nil {
				t.Fatalf("error upserting: %v", err)
			}
		})
	}
}
