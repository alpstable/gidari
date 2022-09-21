package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alpine-hodler/gidari/internal/web/coinbasepro"
	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
)

// ErrRepositoryEncoderExists indicates that an encoder has already been registered for the given url and table.
var ErrRepositoryEncoderExists = fmt.Errorf("repository encoder already exists")

// RepositoryEncoderKey is a struct used to obtain a repository encoder from the registry.
type RepositoryEncoderKey struct {
	host string
}

// NewRepositoryEncoderKey will return a key by using the table and parsing the URL.
func NewRepositoryEncoderKey(u *url.URL) RepositoryEncoderKey {
	return RepositoryEncoderKey{host: u.Host}
}

// NewDefaultRepositoryEncoderKey will return a key for the default encoder. This key will have an empty host.
func NewDefaultRepositoryEncoderKey() RepositoryEncoderKey {
	u, _ := url.Parse("")
	return NewRepositoryEncoderKey(u)
}

// RepositoryEncoder is an interface that defines how to transform data from a web API request into a byte slice that
// can be passed to repository upsert methods.
type RepositoryEncoder interface {
	// Encode will transform the data from a web request into a byte slice that can be passed to repository upsert
	// methods.
	Encode(http.Request, []byte) (*proto.UpsertRequest, error)
}

// RepositoryEncoderRegistry is a map of registered repository encoders.
type RepositoryEncoderRegistry map[RepositoryEncoderKey]RepositoryEncoder

// Register will map a "RepositoryEncoderKey" created by the URL and table to the given encoder. If the encoder has
// already been regisered, this method will throw the "ErrRepositorEncoderExists" error.
func (rer RepositoryEncoderRegistry) Register(u *url.URL, encoder RepositoryEncoder) error {
	key := NewRepositoryEncoderKey(u)
	if rer[key] != nil {
		return ErrRepositoryEncoderExists
	}
	rer[key] = encoder
	return nil
}

// Lookup will lookup the "RepositoryEncoder" using a URL and table.
func (rer RepositoryEncoderRegistry) Lookup(u *url.URL) RepositoryEncoder {
	key := NewRepositoryEncoderKey(u)
	if encoder := rer[key]; encoder != nil {
		return encoder
	}
	return rer[NewDefaultRepositoryEncoderKey()]
}

// RepositoryEncoders is the registry of encoders used to transform web request data into a byte slice that can be
// passed to a repository upsert method. The reason for making RepositoryEncoders a global variable is to (1) avoid
// needing to pass it around to every function that needs to access the data, (2) allow custom registration of encoders
// in the init function, and (3) allow for the possibility of having multiple registries.
var RepositoryEncoders = make(RepositoryEncoderRegistry)

// RegisterCustomEncoders will register custom encoder specific to the default project.
func RegisterCustomEncoders() error {
	// Register the default case
	uri, _ := url.Parse("")
	if err := RepositoryEncoders.Register(uri, new(DefaultRepositoryEncoder)); err != nil {
		return err
	}

	// Register the CoinbasePro Sandbox Candles case
	uri, _ = url.Parse("https://api-public.sandbox.exchange.coinbase.com/candles")
	if err := RepositoryEncoders.Register(uri, new(CBPSandboxEncoder)); err != nil {
		return err
	}
	return nil
}

// DefaultRepositoryEncoder is the encoder used when no other encoder can be found for the registry. It will assume
// that the data from the web request is already correctly formatted for upserting to data storage.
type DefaultRepositoryEncoder struct{}

// Encode will transform the data from arbitrary web API requests into a byte slice that can be passed to repository
// upsert methods.
func (dre *DefaultRepositoryEncoder) Encode(req http.Request, bytes []byte) (*proto.UpsertRequest, error) {
	table, err := tools.ParseDBTableFromURL(req)
	if err != nil {
		return nil, fmt.Errorf("error getting table from request: %v", err)
	}
	return &proto.UpsertRequest{
		Table:    table,
		Data:     bytes,
		DataType: int32(tools.UpsertDataJSON),
	}, nil
}

// CBPSandboxEncoder is the encoder used to transform data from Coinbase Pro Sandbox web requests into bytes that
// can be processed by repository upsert methods.
type CBPSandboxEncoder struct{}

// Encode will transform the data from Coinbase Pro Sandbox web requests into a byte slice that can be passed to
// repository.
func (ccre *CBPSandboxEncoder) Encode(req http.Request, bytes []byte) (*proto.UpsertRequest, error) {
	table, err := tools.ParseDBTableFromURL(req)
	if err != nil {
		return nil, fmt.Errorf("error getting table from request: %v", err)
	}

	const candleMinutesGranularity = "60"

	switch table {
	case "candles":
		granularity := req.URL.Query()["granularity"][0]
		if granularity == candleMinutesGranularity {
			table = "candle_minutes"
		}

		productID := tools.SplitURLPath(req)[1]

		// initialize the slice of candles.
		var candles coinbasepro.Candles

		if err := json.Unmarshal(bytes, &candles); err != nil {
			return nil, fmt.Errorf("error unmarshaling candles: %v", err)
		}
		for _, candle := range candles {
			candle.ProductID = productID
		}

		var err error
		updatedBytes, err := json.Marshal(candles)
		if err != nil {
			return nil, fmt.Errorf("error marshaling candles: %v", err)
		}
		return &proto.UpsertRequest{
			Table:    table,
			Data:     updatedBytes,
			DataType: int32(tools.UpsertDataJSON),
		}, nil
	default:
		u, _ := url.Parse("")
		return RepositoryEncoders.Lookup(u).Encode(req, bytes)
	}
}
