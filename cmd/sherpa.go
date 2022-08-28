package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/alpine-hodler/sherpa/internal/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func main() {
	var configFilepath string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "sherpa",
		Short: "sherpa is an ETL executable for storing web data for analysis",
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()

			bytes, err := ioutil.ReadFile(configFilepath)
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

	cmd.MarkFlagRequired("config")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
