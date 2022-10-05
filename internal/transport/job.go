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
	"io"
	"strings"
	"time"

	"github.com/alpine-hodler/gidari/internal/web"
	"github.com/alpine-hodler/gidari/tools"
	"github.com/sirupsen/logrus"
)

type webJob struct {
	*flattenedRequest
	repoJobs chan<- *repoJob
	logger   *logrus.Logger
}

func newWebJob(cfg *Config, req *flattenedRequest, repoJobs chan<- *repoJob) *webJob {
	return &webJob{
		flattenedRequest: req,
		repoJobs:         repoJobs,
		logger:           cfg.Logger,
	}
}

func webWorker(ctx context.Context, workerID int, jobs <-chan *webJob) {
	for job := range jobs {
		start := time.Now()

		rsp, err := web.Fetch(ctx, job.fetchConfig)
		if err != nil {
			job.logger.Fatal(err)
		}

		bytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			job.logger.Fatal(err)
		}

		job.repoJobs <- &repoJob{b: bytes, req: *rsp.Request, table: job.table}

		// strings.Replace is used to ensure no line endings are present in the user input.
		escapedPath := strings.ReplaceAll(rsp.Request.URL.Path, "\n", "")
		escapedPath = strings.ReplaceAll(escapedPath, "\r", "")

		escapedHost := strings.ReplaceAll(rsp.Request.URL.Host, "\n", "")
		escapedHost = strings.ReplaceAll(escapedHost, "\r", "")

		logInfo := tools.LogFormatter{
			WorkerID:   workerID,
			WorkerName: "web",
			Duration:   time.Since(start),
			Host:       escapedHost,
			Msg:        fmt.Sprintf("web request completed: %s", escapedPath),
		}
		job.logger.Infof(logInfo.String())
	}
}
