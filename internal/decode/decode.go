package decode

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/alpstable/gidari/proto"
)

var (
	ErrFailedToReadBody          = errors.New("failed to read body")
	ErrInvalidResponseForDecoder = errors.New("invalid response for decoder")
	ErrUnsupportedDataType       = fmt.Errorf("unsupported data type")
)

// decodeType will define the type of decoder to use based on the provided Accept header. See the "bestFitDecodeType"
// function for usage.
type decodeType int

const (
	// decodeTypeUnknown is returned when the Accept header does not match any of the provided types.
	decodeTypeUnknown decodeType = iota

	// decodeTypeJSON will decode the request body as JSON.
	decodeTypeJSON
)

// decodeFunc is a function that decodes a response into a slice of IteratorResult.
type decodeFunc func(rsp *http.Response) ([]*proto.IteratorResult, error)

// Decode will use the "Accept" header from an http.Response object to determine which decode function to use. If no
// decode function satisfies the "Accept" header, then an error is returned.
func Decode(rsp *http.Response) ([]*proto.IteratorResult, error) {
	acceptHeader := rsp.Header.Get("Accept")

	switch bestFitDecodeType(acceptHeader) {
	case decodeTypeJSON:
		return decodeJSON(rsp)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponseForDecoder, acceptHeader)
	}
}

// bestFitDecodeType will parse the provided Accept(-Charset|-Encoding|-Language) header and return the header that
// best fits the decoding algorithm. If the "Accept" header is not set, then this method will return a decodeTypeJSON.
// If the "Accept" header is set, but no match is found, then this method will return a decodeTypeUnkown.
//
// See the "acceptSlice.Less" method for more informaiton on how the "best fit" is determined.
func bestFitDecodeType(header string) decodeType {
	accepted := parseAcceptHeader(header)

	// If the header is empty, we default to JSON.
	if len(accepted) == 0 {
		return decodeTypeJSON
	}

	for _, accept := range accepted {
		// If the type is "*" and the subtype is "*", then we default to JSON.
		if accept.typ == "*" && accept.subtype == "*" {
			return decodeTypeJSON
		}

		// If the type is "application" and the subtype is "json" or "*", then we use JSON.
		if accept.typ == "application" && (accept.subtype == "json" || accept.subtype == "*") {
			return decodeTypeJSON
		}
	}

	return decodeTypeUnknown
}

// newInterfaceSlice will create a new slice of interface{} constructs of the "data" type.
func newInterfaceSlice(data interface{}) ([]interface{}, error) {
	var out []interface{}

	dataValue := reflect.ValueOf(data)
	switch dataValue.Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			out = append(out, dataValue.Index(i).Interface())
		}
	case reflect.Map:
		out = append(out, dataValue.Interface())
	case reflect.Struct:
		out = append(out, dataValue.Interface())
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.String,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, dataValue.Kind())
	}

	return out, nil
}

// decodeJSON will decode the body of an http.Response object into a slice of IteratorResult. If the response is a
// an array of JSON objects, then each object will be decoded into an IteratorResult.
func decodeJSON(rsp *http.Response) ([]*proto.IteratorResult, error) {
	// Check if the response body is nil, if it is then return an empty slice.
	if rsp.Body == nil {
		return []*proto.IteratorResult{}, nil
	}

	// read the body of the response using the io.Reader interface.
	_, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToReadBody, err)
	}

	return nil, nil
}
