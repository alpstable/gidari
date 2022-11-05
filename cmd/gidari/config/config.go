package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
	"gopkg.in/yaml.v2"
)

// New takes a path to a config file and returns a gidari.Config object to be used by the CLI to run a transport
// operation.
func New(ctx context.Context, path string) (*gidari.Config, error) {
	cfg, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	cfg.StorageOptions, err = addAllStorage(ctx, cfg.StorageOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to add storage: %w", err)
	}

	cfg.URL, err = url.Parse(cfg.RawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse URL: %w", err)
	}

	if err := addRequestData(ctx, cfg.RateLimitConfig, cfg.Requests); err != nil {
		return nil, fmt.Errorf("unable to add request data: %w", err)
	}

	return cfg, nil
}

// readFile reads the config yaml file from the given path and unmarshals the contents into a "config" struct.
func readFile(path string) (*gidari.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening config file  %s: %v", path, err)
	}

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file stat for reading: %w", err)
	}

	bytes := make([]byte, info.Size())

	_, err = file.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	cfg := new(gidari.Config)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal YAML: %w", err)
	}

	return cfg, nil
}

// addMongoStorage will use the connection string from the gidari.StorageOptions struct to create a new storage
// connection for MongoDB.
func addMongoStorage(ctx context.Context, storage *gidari.StorageOptions) error {
	// Create a new MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://mongo1:27017/defaultcoll")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("unable to connect to MongoDB: %w", err)
	}

	// Plug the client into a Gidari MongoDB Storage adapater.
	mdbStorage, err := gmongo.New(ctx, client)
	if err != nil {
		return fmt.Errorf("unable to create new MongoDB storage: %w", err)
	}

	// Add the storage to the gidari.StorageOptions struct.
	storage.Storage = mdbStorage

	return nil
}

// addStorage will use the connection string from the gidari.StorageOptions struct to create a new storage
// connection.
func addStorage(ctx context.Context, storage *gidari.StorageOptions) error {
	if storage == nil {
		return fmt.Errorf("storage options is nil")
	}

	if storage.ConnectionString == nil {
		return fmt.Errorf("connection string is nil")
	}

	// Check to see if the connection strings is prepended with "mongodb://"
	if strings.HasPrefix(*storage.ConnectionString, "mongodb://") {
		return addMongoStorage(ctx, storage)
	}

	return nil
}

// addAllStorage will use the connection string from the gidari.StorageOptions struct to create a new storage.
func addAllStorage(ctx context.Context, opts []gidari.StorageOptions) ([]gidari.StorageOptions, error) {
	newStorageOptions := make([]gidari.StorageOptions, 0, len(opts))

	for _, opt := range opts {
		newOpts := opt
		newOpts.Close = true

		if err := addStorage(ctx, &newOpts); err != nil {
			return nil, fmt.Errorf("unable to add storage: %w", err)
		}

		newStorageOptions = append(newStorageOptions, newOpts)
	}

	return newStorageOptions, nil
}

// addRequestData cleans up the request data for the requests on a configuration.
func addRequestData(ctx context.Context, cfg *gidari.RateLimitConfig, reqs []*gidari.Request) error {
	// create a rate limiter to pass to all "flattenedRequest". This has to be defined outside of the scope of
	// individual "flattenedRequest"s so that they all share the same rate limiter, even concurrent requests to
	// different endpoints could cause a rate limit error on a web API.
	rateLimiter := rate.NewLimiter(rate.Every(*cfg.Period), *cfg.Burst)

	// Update default request data.
	for _, req := range reqs {
		if req.Method == "" {
			req.Method = http.MethodGet
		}

		if req.Table == "" {
			endpointParts := strings.Split(req.Endpoint, "/")
			req.Table = endpointParts[len(endpointParts)-1]
		}

		req.RateLimiter = rateLimiter
	}

	return nil
}
