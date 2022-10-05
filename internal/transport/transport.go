// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package transport

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
)

// Truncate will truncate the defined tables in the configuration.
func Truncate(ctx context.Context, cfg *Config) error {
	// truncateRequest is a special request that will truncate the table before upserting data.
	truncateRequest := new(proto.TruncateRequest)

	if cfg.Truncate {
		for _, req := range cfg.Requests {
			// Add the table to the list of tables to truncate.
			if req.Truncate != nil && *req.Truncate {
				truncateRequest.Tables = append(truncateRequest.Tables, req.Table)
			}
		}
	} else {
		// checking for request-specific truncate
		for _, req := range cfg.Requests {
			if table := req.Table; req.Truncate != nil && *req.Truncate && table != "" {
				truncateRequest.Tables = append(truncateRequest.Tables, table)
			}
		}
	}

	return truncate(ctx, cfg, truncateRequest)
}

// Upsert will use the configuration file to upsert data from the
//
// For each DNS entry in the configuration file, a repository will be created and used to upsert data. For each
// repository, a transaction will be created and used to upsert data. The transaction will be committed at the end
// of the upsert operation. If the transaction fails, the transaction will be rolled back. Note that it is possible
// for some repository transactions to succeed and others to fail.
func Upsert(ctx context.Context, cfg *Config) error {
	start := time.Now()
	threads := runtime.NumCPU()

	if err := Truncate(ctx, cfg); err != nil {
		return err
	}

	flattenedRequests, err := cfg.flattenRequests(ctx)
	if err != nil {
		return err
	}

	repoConfig, err := newRepoConfig(ctx, cfg, len(flattenedRequests))
	if err != nil {
		return err
	}

	defer repoConfig.closeRepos()

	// Start the repository workers.
	for id := 1; id <= threads; id++ {
		go repositoryWorker(ctx, id, repoConfig)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "repository workers started"}.String())

	webWorkerJobs := make(chan *webJob, len(cfg.Requests))

	// Start the same number of web workers as the cores on the machine.
	for id := 1; id <= threads; id++ {
		go webWorker(ctx, id, webWorkerJobs)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "web workers started"}.String())

	// Enqueue the worker jobs
	for _, req := range flattenedRequests {
		webWorkerJobs <- newWebJob(cfg, req, repoConfig.jobs)
	}

	cfg.Logger.Info(tools.LogFormatter{Msg: "web worker jobs enqueued"}.String())

	// Wait for all of the data to flush.
	for a := 1; a <= len(flattenedRequests); a++ {
		<-repoConfig.done
	}

	// Commit the transactions and check for errors.
	for _, repo := range repoConfig.repos {
		if err := repo.Commit(); err != nil {
			return fmt.Errorf("unable to commit transaction: %w", err)
		}
	}

	logInfo := tools.LogFormatter{Duration: time.Since(start), Msg: "upsert completed"}
	cfg.Logger.Info(logInfo.String())

	return nil
}

func truncate(ctx context.Context, cfg *Config, truncateRequest *proto.TruncateRequest) error {
	start := time.Now()

	repos, closeRepos, err := cfg.repos(ctx)
	if err != nil {
		return err
	}

	defer closeRepos()

	for _, repo := range repos {
		start := time.Now()

		_, err := repo.Truncate(ctx, truncateRequest)
		if err != nil {
			return fmt.Errorf("unable to truncate tables: %w", err)
		}

		rt := repo.Type()
		tables := strings.Join(truncateRequest.Tables, ", ")
		msg := fmt.Sprintf("truncated tables on %q: %v", storage.Scheme(rt), tables)

		logInfo := tools.LogFormatter{
			Duration: time.Since(start),
			Msg:      msg,
		}
		cfg.Logger.Infof(logInfo.String())
	}

	logInfo := tools.LogFormatter{
		Duration: time.Since(start),
		Msg:      "truncate completed",
	}
	cfg.Logger.Info(logInfo.String())

	return nil
}
