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
	"net/http"
	"time"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/repository"
	"github.com/alpine-hodler/gidari/tools"
	"github.com/sirupsen/logrus"
)

type repoCloser func()

type repoJob struct {
	req   http.Request
	b     []byte
	table string
}

type repoConfig struct {
	repos      []repository.Generic
	closeRepos repoCloser
	jobs       chan *repoJob
	done       chan bool
	logger     *logrus.Logger
}

func newRepoConfig(ctx context.Context, cfg *Config, volume int) (*repoConfig, error) {
	repos, closeRepos, err := cfg.repos(ctx)
	if err != nil {
		return nil, err
	}

	return &repoConfig{
		repos:      repos,
		closeRepos: closeRepos,
		jobs:       make(chan *repoJob, volume*len(repos)),
		done:       make(chan bool, volume),
		logger:     cfg.Logger,
	}, nil
}

func repositoryWorker(_ context.Context, workerID int, cfg *repoConfig) {
	for job := range cfg.jobs {
		reqs := []*proto.UpsertRequest{
			{
				Table:    job.table,
				Data:     job.b,
				DataType: int32(tools.UpsertDataJSON),
			},
		}

		for _, req := range reqs {
			for _, repo := range cfg.repos {
				txfn := func(sctx context.Context, repo repository.Generic) error {
					start := time.Now()

					rsp, err := repo.Upsert(sctx, req)
					if err != nil {
						cfg.logger.Fatalf("error upserting data: %v", err)

						return fmt.Errorf("error upserting data: %w", err)
					}

					rt := repo.Type()

					msg := fmt.Sprintf("partial upsert completed: %s.%s", storage.Scheme(rt), req.Table)
					logInfo := tools.LogFormatter{
						WorkerID:      workerID,
						WorkerName:    "repository",
						Duration:      time.Since(start),
						Msg:           msg,
						UpsertedCount: rsp.UpsertedCount,
						MatchedCount:  rsp.MatchedCount,
					}

					cfg.logger.Infof(logInfo.String())

					return nil
				}
				// Put the data onto the transaction channel for storage.
				repo.Transact(txfn)
			}
		}

		cfg.done <- true
	}
}
