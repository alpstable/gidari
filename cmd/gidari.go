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
	"os"

	"github.com/alpine-hodler/gidari"
	"github.com/alpine-hodler/gidari/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed bash-completion.sh
var bashCompletion string

func main() {
	// configFilepath is the path to the configuration file.
	var configFilepath string

	// verbose is a flag that enables verbose logging.
	var verbose bool

	cmd := &cobra.Command{
		Long: "Gidari is a tool for querying web APIs and persisting resultant data onto local storage\n" +
			"using a configuration file.",

		Use:                    "gidari",
		Short:                  "Persisted data from the web to your database",
		Example:                "gidari --config config.yaml",
		BashCompletionFunction: bashCompletion,
		Deprecated:             "",
		Version:                version.Gidari,

		Run: func(_ *cobra.Command, args []string) { run(configFilepath, verbose, args) },
	}

	cmd.Flags().StringVar(&configFilepath, "config", "c", "path to configuration")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "print log data as the binary executes")

	if err := cmd.MarkFlagRequired("config"); err != nil {
		logrus.Fatalf("error marking flag as required: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(configFilepath string, verboseLogging bool, _ []string) {
	file, err := os.Open(configFilepath)
	if err != nil {
		log.Fatalf("error opening config file  %s: %v", configFilepath, err)
	}

	cfg, err := gidari.NewConfig(context.Background(), file)
	if err != nil {
		log.Fatalf("error creating new config: %v", err)
	}

	if verboseLogging {
		cfg.Logger.SetOutput(os.Stdout)
		cfg.Logger.SetLevel(logrus.InfoLevel)
	}

	err = gidari.Transport(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to transport data: %v", err)
	}
}
