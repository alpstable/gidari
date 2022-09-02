package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alpine-hodler/sherpa/internal/web/coinbasepro"
	"github.com/sirupsen/logrus"
)

// The goal of a repository encoder is to transform the data from a web request
// into a byte slice that can be passed to a repository upsert method.
//
// - MUST be able to register encoders upto table + API (e.g. candles + CoinbasePro)
// - MUST be mutable, user should be able to register encoders

// ErrRepositoryEncoderExists indicates that an encoder has already been registered for the given url and table.
var ErrRepositoryEncoderExists = fmt.Errorf("repository encoder already exists")

// RepositoryEncoderKey is a struct used to obtain a repository encoder from the registry.
type RepositoryEncoderKey struct {
	host string
}

// NewRepositoryEncoderKey will return a key by using the table and parsing the URL.
func NewRepositoryEncoderKey(u url.URL) RepositoryEncoderKey {
	return RepositoryEncoderKey{host: u.Host}
}

// NewDefaultRepositoryEncoderKey will return a key with not host and no table, which is implicitly the default case.
func NewDefaultRepositoryEncoderKey() RepositoryEncoderKey {
	u, _ := url.Parse("")
	return NewRepositoryEncoderKey(*u)
}

// RepositoryEncoder is used to transform data from a web request into a byte slice that can be passed to a repository
// upsert method.
type RepositoryEncoder interface {
	Encode(http.Request, *[]byte) (string, error)
}

// RepositoryEncoderRegistry is a map of registered repository encoders.
type RepositoryEncoderRegistry map[RepositoryEncoderKey]RepositoryEncoder

// Register will map a "RepositoryEncoderKey" created by the URL and table to the given encoder. If the encoder has
// already been regisered, this method will throw the "ErrRepositorEncoderExists" error.
func (rer RepositoryEncoderRegistry) Register(u url.URL, encoder RepositoryEncoder) error {
	key := NewRepositoryEncoderKey(u)
	if rer[key] != nil {
		return ErrRepositoryEncoderExists
	}
	rer[key] = encoder
	return nil
}

// Lookup will lookup the "RepositoryEncoder" using a URL and table.
func (rer RepositoryEncoderRegistry) Lookup(u url.URL) RepositoryEncoder {
	key := NewRepositoryEncoderKey(u)
	if encoder := rer[key]; encoder != nil {
		return encoder
	}
	return rer[NewDefaultRepositoryEncoderKey()]
}

// RepositoryEncoders is the registry of encoders used to transform web request data into a byte slice that can be
// passed to a repository upsert method.
var RepositoryEncoders = make(RepositoryEncoderRegistry)

func init() {
	// Register the default case
	u, _ := url.Parse("")
	if err := RepositoryEncoders.Register(*u, new(DefaultRepositoryEncoder)); err != nil {
		logrus.Fatalf("error registering default encoder: %v", err)
	}

	// Register the CoinbasePro Sandbox Candles case
	u, _ = url.Parse("https://api-public.sandbox.exchange.coinbase.com/candles")
	if err := RepositoryEncoders.Register(*u, new(CBPSandboxCandlesEncoder)); err != nil {
		logrus.Fatalf("error registering Coinbase Pro Sandbox Candles encoder: %v", err)
	}
}

// DefaultRepositoryEncoder is the encoder used when no other encoder can be found for the registry. It will assume
// that the data from the web request is already correctly formatted for upserting to data storage.
type DefaultRepositoryEncoder struct{}

func (dre *DefaultRepositoryEncoder) Encode(req http.Request, _ *[]byte) (string, error) {
	endpoint := strings.TrimPrefix(req.URL.EscapedPath(), "/")
	endpointParts := strings.Split(endpoint, "/")

	table := endpointParts[len(endpointParts)-1]
	return table, nil
}

type CBPSandboxCandlesEncoder struct{}

func (ccre *CBPSandboxCandlesEncoder) Encode(req http.Request, dst *[]byte) (string, error) {
	endpoint := strings.TrimPrefix(req.URL.EscapedPath(), "/")
	endpointParts := strings.Split(endpoint, "/")

	table := endpointParts[len(endpointParts)-1]

	switch table {
	case "candles":
		granularity := req.URL.Query()["granularity"][0]
		switch granularity {
		case "60":
			table = "candle_minutes"
		}

		productID := endpointParts[1]
		var candles coinbasepro.Candles
		if err := json.Unmarshal(*dst, &candles); err != nil {
			return "", err
		}
		for _, candle := range candles {
			candle.ProductID = productID
		}

		var err error
		*dst, err = json.Marshal(candles)
		if err != nil {
			return "", err
		}
	}

	return table, nil
}
