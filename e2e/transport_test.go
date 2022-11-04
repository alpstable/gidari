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
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v2"
)

func newConfig(ctx context.Context, file *os.File) (*gidari.Config, error) {
	var cfg gidari.Config

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cfg.URL, err = url.Parse(cfg.RawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %w", err)
	}

	// create a rate limiter to pass to all "flattenedRequest". This has to be defined outside of the scope of
	// individual "flattenedRequest"s so that they all share the same rate limiter, even concurrent requests to
	// different endpoints could cause a rate limit error on a web API.
	rateLimiter := rate.NewLimiter(rate.Every(*cfg.RateLimitConfig.Period), *cfg.RateLimitConfig.Burst)

	// Update default request data.
	for _, req := range cfg.Requests {
		if req.Method == "" {
			req.Method = http.MethodGet
		}

		if req.Table == "" {
			endpointParts := strings.Split(req.Endpoint, "/")
			req.Table = endpointParts[len(endpointParts)-1]
		}

		req.RateLimiter = rateLimiter
	}

	// Add storage using the connection strings.
	for idx, stgOpts := range cfg.StorageOptions {
		// If the connection string has "mongodb:" as a prefix, then it is a MongoDB connection string.
		connStr := stgOpts.ConnectionString
		if connStr != nil && strings.HasPrefix(*connStr, "mongodb:") {
			// Create a MongoDB client using the official MongoDB Go Driver.
			clientOptions := options.Client().ApplyURI(*connStr)
			client, _ := mongo.Connect(ctx, clientOptions)

			// Ping the client
			if err := client.Ping(ctx, nil); err != nil {
				return nil, fmt.Errorf("unable to ping MongoDB client: %w", err)
			}

			// Create a MongoDB storage.
			mstg, err := gmongo.New(ctx, client)
			if err != nil {
				return nil, fmt.Errorf("unable to create MongoDB storage: %w", err)
			}

			cfg.StorageOptions[idx].Storage = mstg
		}
	}

	return &cfg, nil
}

func TestUpsert(t *testing.T) {
	t.Parallel()

	// Iterate over the fixtures/upsert directory and run each gidari.ration file.
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

			cfg, err := newConfig(ctx, file)
			if err != nil {
				t.Fatalf("error creating gidari. %v", err)
			}

			// Fill in the authentication details for the fixture.
			cfgAuth := cfg.Authentication
			if cfgAuth.APIKey != nil {
				// The "passhprase" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = gidari.Authentication{
					APIKey: &gidari.APIKey{
						Key:        os.Getenv(cfgAuth.APIKey.Key),
						Secret:     os.Getenv(cfgAuth.APIKey.Secret),
						Passphrase: os.Getenv(cfgAuth.APIKey.Passphrase),
					},
				}
			}

			if cfgAuth.Auth2 != nil {
				// The "bearer" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = gidari.Authentication{
					Auth2: &gidari.Auth2{
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
