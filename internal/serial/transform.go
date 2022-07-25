package serial

import (
	"time"
)

// Transform interfaces data that can be deserialized/unmarshalled by type.
type Transform interface {
	// unmarshal is a wrapper for setting data in model unmarhallers
	Unmarshal(key string, fn func(interface{}) error) error
	UnmarshalBool(name string, v *bool) error
	UnmarshalFloatString(name string, v *float64) error
	UnmarshalFloat(name string, v *float64) error
	UnmarshalInt64(name string, v *int64) error
	UnmarshalInt32(name string, v *int32) error
	UnmarshalInt(name string, v *int) error
	UnmarshalStringSlice(name string, v *[]string) error
	UnmarshalString(name string, v *string) error
	UnmarshalStructSlice(name string, v _jsonStructSlice, template interface{}) error
	UnmarshalStruct(name string, v interface{}) error
	UnmarshalTime(layout string, name string, v *time.Time) error
	UnmarshalUnixString(name string, v *time.Time) error
	Value(key string) interface{}
}

// NewJSONTransform will return a new json transform that will create a map of objects from a byte
// stream of JSON to be deserialized.
func NewJSONTransform(d []byte) (Transform, error) {
	return make_jsonTransform(d)
}
