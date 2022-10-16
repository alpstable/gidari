//go:build pginteg

// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0\n
package postgres

import (
	"context"
	"testing"

	"github.com/alpstable/gidari/internal/proto"
)

const defaultConnectionString = "postgresql://root:root@postgres1:5432/defaultdb?sslmode=disable"

func TestPostgres(t *testing.T) {
	const defaultTestTable = "tests1"
	const listTablesTable = "lttests1"
	const listPrimaryKeysTable = "pktests1"

	defaultData := map[string]interface{}{
		"test_string": "test",
		"id":          "1",
	}

	ctx := context.Background()

	pg, err := New(ctx, defaultConnectionString)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	proto.RunTest(context.Background(), t, pg, func(runner *proto.TestRunner) {
		runner.AddCloseDBCases(
			[]proto.TestCase{
				{
					Name: "close postgres",
				},
			}...,
		)

		runner.AddListPrimaryKeysCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listPrimaryKeysTable,
					ExpectedPrimaryKeys: map[string][]string{
						listPrimaryKeysTable: {"test_string"},
					},
				},
			}...,
		)

		runner.AddListTablesCases(
			[]proto.TestCase{
				{
					Name:  "single",
					Table: listTablesTable,
				},
			}...,
		)

		runner.AddUpsertTxnCases(
			[]proto.TestCase{
				{
					Name:               "commit",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 8192,
					Data:               defaultData,
				},
				{
					Name:               "rollback",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					Rollback:           true,
					Data:               defaultData,
				},
				{
					Name:               "rollback on error",
					Table:              defaultTestTable,
					ExpectedUpsertSize: 0,
					ForceError:         true,
					Data:               defaultData,
				},
			}...,
		)
	})
}
