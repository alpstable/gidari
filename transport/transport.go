package transport

import (
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
	APIKey APIKey `yaml:"apiKey"`
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
	URL            string         `yaml:"url"`
	Authentication Authentication `yaml:"authentication"`
	DNSList        []string       `yaml:"dnsList"`
	Requests       []Request      `yaml:"requests"`

	Logger *logrus.Logger
}

// Upsert will use the configuration file to upsert data from the
func Upsert(cfg *Config) error {
	cfg.Logger.Info("meep moop morp")
	return nil
}
