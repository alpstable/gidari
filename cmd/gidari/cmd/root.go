// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package cmd

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagConfig  = "config"
	flagVerbose = "verbose"
)

func RootCommand() *cobra.Command {
	var (
		configFilepath string // configFilepath is the path to the configuration file.
		verbose        bool   // verbose is a flag that enables verbose logging.
	)

	rootCMD := &cobra.Command{
		Long: "Gidari is a tool for querying web APIs and persisting resultant data onto local storage\n" +
			"using a configuration file.",

		Use:     "gidari",
		Short:   "Persist data from the web to your database",
		Example: "gidari --config config.yaml",
		Version: version.Gidari,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			value, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				//nolint:wrapcheck // need not wrap the error
				return err
			}

			if !strings.HasSuffix(value, ".yaml") && !strings.HasSuffix(value, ".yml") {
				//nolint:goerr113 // don't have static error
				return errors.New("configuration must be a YAML document")
			}

			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error { return runE(configFilepath, verbose) },
	}

	rootCMD.PersistentFlags().StringVarP(&configFilepath, flagConfig, "c", "", "path to configuration")
	rootCMD.PersistentFlags().BoolVar(&verbose, flagVerbose, false, "print log data as the binary executes")

	_ = rootCMD.MarkPersistentFlagRequired(flagConfig)

	return rootCMD
}

//nolint:wrapcheck // need not wrap the error
func runE(configFilepath string, verbose bool) error {
	file, err := os.Open(configFilepath)
	if err != nil {
		return err
	}

	cfg, err := gidari.NewConfig(context.Background(), file)
	if err != nil {
		return err
	}

	if verbose {
		cfg.Logger.SetOutput(os.Stdout)
		cfg.Logger.SetLevel(logrus.InfoLevel)
	}

	err = gidari.Transport(context.Background(), cfg)
	if err != nil {
		return err
	}

	return nil
}
