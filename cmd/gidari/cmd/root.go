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
	"fmt"
	"os"
	"strings"

	"github.com/alpine-hodler/gidari"
	"github.com/alpine-hodler/gidari/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagConfig  = "config"
	flagVerbose = "verbose"
)

var (
	configFilepath string // configFilepath is the path to the configuration file.
	verbose        bool   // verbose is a flag that enables verbose logging.
)

var rootCMD = &cobra.Command{
	Long: "Gidari is a tool for querying web APIs and persisting resultant data onto local storage\n" +
		"using a configuration file.",

	Use:     "gidari",
	Short:   "Persist data from the web to your database",
	Example: "gidari --config config.yaml",
	Version: version.Gidari,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		v, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(v, ".yaml") && !strings.HasSuffix(v, ".yml") {
			return fmt.Errorf("configuration must be a YAML document")
		}

		return nil
	},
	RunE: runE,
}

func Execute() error {
	rootCMD.PersistentFlags().StringVarP(&configFilepath, flagConfig, "c", "", "path to configuration")
	rootCMD.PersistentFlags().BoolVar(&verbose, flagVerbose, false, "print log data as the binary executes")

	_ = rootCMD.MarkPersistentFlagRequired(flagConfig)

	return rootCMD.Execute()
}

func runE(cmd *cobra.Command, _ []string) error {
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
