// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package main

import (
	"context"
	_ "embed" // Embed external data.
	"log"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/cmd/gidari/config"
	"github.com/alpstable/gidari/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	// configFilepath is the path to the configuration file.
	var configFilepath string

	// verbose is a flag that enables verbose logging.
	var verbose bool

	cmd := &cobra.Command{
		Long:       "Gidari uses a configuration file for querying web APIs and persisting resultant data.",
		Use:        "gidari",
		Short:      "Persisted data from the web to your database",
		Example:    "gidari --config config.yaml",
		Deprecated: "",
		Version:    version.Gidari,

		Run: func(_ *cobra.Command, args []string) { run(configFilepath, verbose, args) },
	}

	cmd.Flags().StringVar(&configFilepath, "config", "", "path to configuration")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "print log data as the binary executes")

	if err := cmd.MarkFlagRequired("config"); err != nil {
		logrus.Fatalf("error marking flag as required: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(configFilepath string, verboseLogging bool, _ []string) {
	ctx := context.Background()

	cfg, err := config.New(ctx, configFilepath)
	if err != nil {
		logrus.Fatalf("error creating configuration: %v", err)
	}

	if err := gidari.Transport(ctx, cfg); err != nil {
		logrus.Fatalf("error running transport: %v", err)
	}
}
