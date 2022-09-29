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
func RegisterEncoders(rer RepositoryEncoderRegistry, encoders ...func(RepositoryEncoderRegistry) error) error {
	for _, encoder := range encoders {
		if err := encoder(rer); err != nil {
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

type RegisteredRepositoryEncoder struct {
	RepositoryEncoder
}

// Lookup will lookup the "RepositoryEncoder" using a URL and table.
func (rer RepositoryEncoderRegistry) Lookup(u *url.URL) *RegisteredRepositoryEncoder {
	key := NewRepositoryEncoderKey(u)
	if encoder := rer[key]; encoder != nil {
		return &RegisteredRepositoryEncoder{encoder}
	}

	return &RegisteredRepositoryEncoder{new(DefaultRepositoryEncoder)}
}

// RegisterDefaultEncoder will register custom encoder specific to the default project.
func RegisterDefaultEncoder(rer RepositoryEncoderRegistry) error {
	uri, err := url.Parse("")
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}

	if err := rer.Register(uri, new(DefaultRepositoryEncoder)); err != nil {
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
