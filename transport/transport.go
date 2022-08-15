package transport

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"github.com/alpine-hodler/driver/internal"
	"github.com/alpine-hodler/driver/proto"
	"github.com/alpine-hodler/driver/repository"
	"github.com/alpine-hodler/driver/web"
	"github.com/alpine-hodler/driver/web/auth"
	"github.com/sirupsen/logrus"
)

// APIKey is one method of HTTP(s) transport that requires a passphrase, key, and secret.
type APIKey struct {
	Passphrase string `yaml:"passphrase"`
	Key        string `yaml:"key"`
	Secret     string `yaml:"secret"`
}

// Authentication is the credential information to be used to construct an HTTP(s) transport for accessing the API.
type Authentication struct {
	APIKey *APIKey `yaml:"apiKey"`
}

// Request is the information needed to query the web API for data to transport.
type Request struct {
	// Method is the HTTP(s) method used to construct the http request to fetch data for storage.
	Method string `yaml:"method"`

	// Endpoint is the fragment of the URL that will be used to request data from the API. This value can include
	// query parameters.
	Endpoint string `yaml:"endpoint"`

	// RateLimitBurstCap represents the number of requests that can be made per second to the endpoint. The
	// value of this should come from the documentation in the underlying API.
	RateLimitBurstCap int `yaml:"ratelimit"`
}

// Config is the configuration used to query data from the web using HTTP requests and storing that data using
// the repositories defined by the "DNSList".
type Config struct {
	URL                     string         `yaml:"url"`
	Authentication          Authentication `yaml:"authentication"`
	DNSList                 []string       `yaml:"dnsList"`
	GlobalRateLimitBurstCap int            `yaml:"ratelimit"`
	Requests                []*Request     `yaml:"requests"`

	Logger   *logrus.Logger
	Truncate bool
}

// connect will attempt to connect to the web API client. Since there are multiple ways to build a transport given the
// authentication data, this method will exhuast every transport option in the "Authentication" struct.
func (cfg *Config) connect(ctx context.Context) (*web.Client, error) {
	if apiKey := cfg.Authentication.APIKey; apiKey != nil {
		return web.NewClient(ctx, auth.NewAPIKey().
			SetURL(cfg.URL).
			SetKey(apiKey.Key).
			SetPassphrase(apiKey.Passphrase).
			SetSecret(apiKey.Secret))
	}
	return nil, nil
}

func newFetchConfig(ctx context.Context, cfg *Config, req *Request, client *web.Client) (*web.FetchConfig, error) {
	u, err := url.JoinPath(cfg.URL, req.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error joining url %q to endpoint %q: %v", cfg.URL, req.Endpoint, err)
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}
	webcfg := &web.FetchConfig{
		Client: client,
		Method: req.Method,
		URL:    parsedURL,
	}

	// Use the request's local rate limit, if it is not non-zero then use the global rate limit.
	if req.RateLimitBurstCap != 0 {
		webcfg.RateLimitBurstCap = req.RateLimitBurstCap
	} else {
		webcfg.RateLimitBurstCap = cfg.GlobalRateLimitBurstCap
	}

	return webcfg, nil
}

type reposet map[uint8]repository.Generic

type repoJob struct {
	b   []byte
	url *url.URL
}

type repoConfig struct {
	set      reposet
	jobs     <-chan *repoJob
	done     chan bool
	logger   *logrus.Logger
	truncate bool
}

func repositoryWorker(ctx context.Context, id int, cfg *repoConfig) {
	for job := range cfg.jobs {
		// ? Should we put the repos in a worker and run them concurrently as well?
		for _, repo := range cfg.set {
			table := strings.TrimPrefix(job.url.EscapedPath(), "/")
			cfg.logger.Infof("{status: started, worker: repo, id: %v, table: %s}", id, table)

			if cfg.truncate {
				treq := new(proto.TruncateTablesRequest)
				treq.Tables = append(treq.Tables, table)
				if err := repo.TruncateTables(ctx, treq); err != nil {
					cfg.logger.Fatal(err)
				}

			}

			rsp := new(proto.CreateResponse)
			if err := repo.UpsertJSON(ctx, table, job.b, rsp); err != nil {
				cfg.logger.Fatal(err)
			}
			cfg.logger.Infof("{status: finished, worker: repo, id: %v, table: %s}", id, table)
		}
		cfg.done <- true
	}
}

type webWorkerJob struct {
	repoJobs    chan<- *repoJob
	client      *web.Client
	fetchConfig *web.FetchConfig
	logger      *logrus.Logger
}

func webWorker(ctx context.Context, id int, jobs <-chan *webWorkerJob) {
	for job := range jobs {
		bytes, err := web.Fetch(ctx, job.fetchConfig)
		if err != nil {
			job.logger.Fatal(err)
		}
		job.repoJobs <- &repoJob{b: bytes, url: job.fetchConfig.URL}
		job.logger.Infof("request completed: (id=%v) %s", id, job.fetchConfig.URL.String())
	}
}

// Upsert will use the configuration file to upsert data from the
func Upsert(ctx context.Context, cfg *Config) error {
	client, err := cfg.connect(ctx)
	if err != nil {
		return fmt.Errorf("unable to connect to client: %v", err)
	}
	cfg.Logger.Info("connection establed")

	// Make a set of repositories from the "DNSList"
	set := make(reposet)
	for _, dns := range cfg.DNSList {
		stg, err := internal.NewStorage(ctx, dns)
		if err != nil {
			return err
		}
		repo := repository.New(ctx, stg)
		set[repo.Type()] = repo
	}
	cfg.Logger.Info("repositories indexed")

	// ? how do we make this a limited buffer?
	repoJobCh := make(chan *repoJob)

	// Construct the repository worker configuration.
	repoWorkerCfg := &repoConfig{
		set:      set,
		logger:   cfg.Logger,
		done:     make(chan bool, len(cfg.Requests)),
		jobs:     repoJobCh,
		truncate: cfg.Truncate,
	}

	// For each repository defined in the DNSList, create a go routine to schedule upserts concurrently.
	for id := 1; id <= len(cfg.DNSList); id++ {
		go repositoryWorker(ctx, id, repoWorkerCfg)
	}
	cfg.Logger.Info("repository workers started")

	webWorkerJobs := make(chan *webWorkerJob, len(cfg.Requests))

	// Start the same number of web workers as the cores on the machine.
	for id := 1; id <= runtime.NumCPU(); id++ {
		go webWorker(ctx, id, webWorkerJobs)
	}
	cfg.Logger.Info("web workers started")

	// Enqueue the web worker jobs.
	for _, req := range cfg.Requests {
		fetchConfig, err := newFetchConfig(ctx, cfg, req, client)
		if err != nil {
			return err
		}

		// Construct the web worker configuration.
		webWorkerJobs <- &webWorkerJob{
			repoJobs:    repoJobCh,
			client:      client,
			fetchConfig: fetchConfig,
			logger:      cfg.Logger,
		}
	}
	cfg.Logger.Info("web worker jobs enqueued")

	// Wait for all of the data to flush.
	for a := 1; a <= len(cfg.Requests); a++ {
		<-repoWorkerCfg.done
	}
	cfg.Logger.Info("repository workers finished")

	return nil
}
