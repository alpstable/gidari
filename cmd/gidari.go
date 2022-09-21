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
	var configFilepath string
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

			// ctx := context.Background()

			// apiKey := cfg.Authentication.APIKey
			// client, err := web.NewClient(ctx, transport.NewAPIKey().
			// 	SetKey(apiKey.Key).
			// 	SetPassphrase(apiKey.Passphrase).
			// 	SetSecret(apiKey.Secret).
			// 	SetURL(cfg.URL))
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// for _, dns := range cfg.DNSList {
			// 	stg, err := storage.New(ctx, dns)
			// 	if err != nil {
			// 		log.Fatalf("error connecting to DNS %q: %v", dns, err)
			// 	}

			// 	repo := repository.NewCoinbasePro(ctx, stg)

			// 	for _, request := range requests {
			// 		rsplit := strings.Split(request, " ")
			// 		method := rsplit[0]
			// 		endpoint := rsplit[1]

			// 		u, err := url.JoinPath(cfg.URL, endpoint)
			// 		if err != nil {
			// 			log.Fatalf("error joining url %q to endpoint %q: %v", cfg.URL, endpoint, err)
			// 		}

			// 		parsedURL, err := url.Parse(u)
			// 		if err != nil {
			// 			log.Fatalf("error parsing URL: %v", err)
			// 		}

			// 		cfg := &web.FetchConfig{
			// 			Client: client,
			// 			Method: method,
			// 			URL:    parsedURL,
			// 		}

			// 		body, err := web.Fetch(ctx, cfg)
			// 		if err != nil {
			// 			log.Fatalf("error fetching accounts: %v", err)
			// 		}
			// 		defer body.Close()

			// 		table := strings.TrimPrefix(parsedURL.EscapedPath(), "/")

			// 		rsp := new(proto.CreateResponse)
			// 		if err = repo.UpsertJSON(ctx, table, body, rsp); err != nil {
			// 			log.Fatalf("error upserting data: %v", err)
			// 		}
			// 	}
			//}

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
