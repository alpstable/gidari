package main

import (
	"context"
	_ "embed" // Embed external data.
	"log"
	"os"

	"github.com/alpine-hodler/gidari/internal/transport"
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
		Example:                "gidari -c config.yaml -v",
		BashCompletionFunction: bashCompletion,
		Deprecated:             "",
		Version:                version.Gidari,

		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()

			// Register supported encoders.
			err := transport.RegisterEncoders(transport.RegisterDefaultEncoder,
				transport.RegisterCBPEncoder)
			if err != nil {
				log.Fatalf("error registering encoders: %v", err)
			}

			bytes, err := os.ReadFile(configFilepath)
			if err != nil {
				log.Fatalf("error reading config file  %s: %v", configFilepath, err)
			}

			cfg, err := transport.NewConfig(bytes)
			if err != nil {
				log.Fatalf("error creating config: %v", err)
			}

			cfg.Logger = logrus.New()

			// If the user has not set the verbose flag, only log fatals.
			if !verbose {
				cfg.Logger.SetLevel(logrus.FatalLevel)
			}

			if err := transport.Upsert(ctx, cfg); err != nil {
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
