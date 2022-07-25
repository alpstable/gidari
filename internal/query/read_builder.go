package query

import (
	"fmt"

	"github.com/alpine-hodler/driver/data/proto"
)

type ReadBuilderType uint8

const (
	CoinbaseProCandleMinutesReadBuilder ReadBuilderType = iota
	PolygonBarMinutesReadBuilder
)

// Bytes will conver the byte for readBuilderType to a slice of bytes, this is for protobuf since protobuf does not
// support a byte type, only bytes.
func (readBuilderType ReadBuilderType) Bytes() []byte {
	return []byte{uint8(readBuilderType)}
}

type StorageType uint8

const (
	PostgresStorage StorageType = iota
	MongoStorage
)

// Builder constructs the query and the arguments for the query to be executed by some storage operation. The data
// returned by this builder should be semi-agnostic. The query portion must be specific to a storage device, but the
// args portion should not.
type ReadBuilder interface {
	// ReaderArgs will return arguments to insert into the query, such as the input to a prepared statement.
	ReaderArgs(*proto.ReadRequest) ([]interface{}, error)

	// ReaderQuery is a factory for returing a query as a bytes array to be processed by the storage object.
	ReaderQuery(StorageType, ...interface{}) ([]byte, error)
}

// GetReadBuilder takes a ReadBuilderType and returns a ReadBuilder-interface object for building read queries for a
// semi-agnostic storage device.
func GetReadBuilder(readBuilderType ReadBuilderType) (ReadBuilder, error) {
	switch readBuilderType {
	case CoinbaseProCandleMinutesReadBuilder:
		return newCoinbaseProCandleMinutes()
	case PolygonBarMinutesReadBuilder:
		return newPolygonBarMinutes()
	default:
		return nil, fmt.Errorf("read builder %q not supported", readBuilderType)
	}
}
