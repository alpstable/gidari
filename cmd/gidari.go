package main

import (
	"context"
	"log"
	"os"

	"github.com/alpine-hodler/gidari/internal/transport"
	"github.com/alpine-hodler/gidari/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	_ "embed" // Embed external data.
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
		Example:                "gidari -c config.yaml -v",
		BashCompletionFunction: bashCompletion,
		Deprecated:             "",
		Version:                version.Gidari,

		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()

			bytes, err := os.ReadFile(configFilepath)
			if err != nil {
				log.Fatalf("error reading config file  %s: %v", configFilepath, err)
			}

			var cfg transport.Config
			if err := yaml.Unmarshal(bytes, &cfg); err != nil {
				log.Fatalf("error unmarshaling data: %v", err)
			}

			cfg.Logger = logrus.New()

			// If the user has not set the verbose flag, only log fatals.
			if !verbose {
				cfg.Logger.SetLevel(logrus.FatalLevel)
			}

			if err := transport.Upsert(ctx, &cfg); err != nil {
				log.Fatalf("error upserting data: %v", err)
			}
		},
	}

	cmd.Flags().StringVar(&configFilepath, "config", "c", "path to configuration")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "print log data as the binary executes")

	if err := cmd.MarkFlagRequired("config"); err != nil {
		logrus.Fatalf("error marking flag as required: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
