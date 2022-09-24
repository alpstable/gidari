package transport

import (
	"fmt"
	"net/http"
	"net/url"

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
	// methods. This method returns a slice of upserts, in the event that the data to encode will needs to be split
	// into multiple table upserts, e.g. relational data.
	Encode(http.Request, []byte) ([]*proto.UpsertRequest, error)
}

// RepositoryEncoderRegistry is a map of registered repository encoders.
type RepositoryEncoderRegistry map[RepositoryEncoderKey]RepositoryEncoder

// RegisterEncoders will register all listed encoders for the deployment.
func RegisterEncoders(encoders ...func() error) error {
	for _, encoder := range encoders {
		if err := encoder(); err != nil {
			return err
		}
	}

	return nil
}

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

// RegisterDefaultEncoder will register custom encoder specific to the default project.
func RegisterDefaultEncoder() error {
	uri, err := url.Parse("")
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}

	if err := RepositoryEncoders.Register(uri, new(DefaultRepositoryEncoder)); err != nil {
		return fmt.Errorf("error registering default encoder: %w", err)
	}

	return nil
}

// DefaultRepositoryEncoder is the encoder used when no other encoder can be found for the registry. It will assume
// that the data from the web request is already correctly formatted for upserting to data storage.
type DefaultRepositoryEncoder struct{}

// Encode will transform the data from arbitrary web API requests into a byte slice that can be passed to repository
// upsert methods.
func (dre *DefaultRepositoryEncoder) Encode(req http.Request, bytes []byte) ([]*proto.UpsertRequest, error) {
	table, err := tools.ParseDBTableFromURL(req)
	if err != nil {
		return nil, fmt.Errorf("error getting table from request: %w", err)
	}

	defaultUpsert := &proto.UpsertRequest{
		Table:    table,
		Data:     bytes,
		DataType: int32(tools.UpsertDataJSON),
	}

	return []*proto.UpsertRequest{defaultUpsert}, nil
}
