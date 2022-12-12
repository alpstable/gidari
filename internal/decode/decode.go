package decode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/alpstable/gidari/proto"
)

var (
	ErrFailedToReadBody          = errors.New("failed to read body")
	ErrFailedToUnmarshalJSONData = errors.New("failed to unmarshal JSON data")
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

// Decode will use the "Accept" header from an http.Response object to determine the "best fit" decode function to use.
// If no decode function satisfies the "Accept" header, then an error is returned.
//
// See the "bestFitDecodeType" function and the "acceptSlice.Less" method for more information on how the "Accept"
// header is parsed as "best fit".
func Decode(rsp *http.Response) ([]*proto.IteratorResult, error) {
	acceptHeader := rsp.Header.Get("Accept")
	switch bestFitDecodeType(acceptHeader) {
	case decodeTypeJSON:
		return decodeJSON(rsp)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidResponseForDecoder, acceptHeader)
	}
}

// isDecodeTypeJSON will check if the provided "accept" struct is typed for decoding into JSON.
func isDecodeTypeJSON(accept accept) bool {
	return accept.typ == "application" && (accept.subtype == "json" || accept.subtype == "*") ||
		accept.typ == "*" && accept.subtype == "*"
}

// bestFitDecodeType will parse the provided Accept(-Charset|-Encoding|-Language) header and return the header that
// best fits the decoding algorithm. If the "Accept" header is not set, then this method will return a decodeTypeJSON.
// If the "Accept" header is set, but no match is found, then this method will return a decodeTypeUnkown.
//
// See the "acceptSlice.Less" method for more informaiton on how the "best fit" is determined.
func bestFitDecodeType(header string) decodeType {
	decodeType := decodeTypeUnknown
	for _, accept := range parseAcceptHeader(header) {
		if isDecodeTypeJSON(accept) {
			decodeType = decodeTypeJSON
		}
	}

	return decodeType
}

type irConfig struct {
	// enc is the encoding function for an iterator result, converting the data into a byte slice.
	enc func(interface{}) ([]byte, error)
	uri *url.URL
}

// mapToIR will convert a map[string]interface{} into an IteratorResult.
func mapToIR(data interface{}, cfg irConfig) ([]*proto.IteratorResult, error) {
	bytes, err := cfg.enc(data)
	if err != nil {
		return nil, err
	}

	if bytes == nil || len(bytes) == 0 {
		return []*proto.IteratorResult{}, nil
	}

	return []*proto.IteratorResult{
		{
			Data: bytes,
			URL:  cfg.uri.String(),
		},
	}, nil
}

// sliceToIR will attempt to create a flat slice of proto.IteratorResult objects from the dataValue provided.
func sliceToIR(dataValue reflect.Value, cfg irConfig) ([]*proto.IteratorResult, error) {
	out := []*proto.IteratorResult{}
	for i := 0; i < dataValue.Len(); i++ {
		ir, err := newIteratorResults(dataValue.Index(i).Interface(), cfg)
		if err != nil {
			return nil, err
		}

		out = append(out, ir...)
	}

	return out, nil
}

// makeIteratorResultSlice will create a new slice of interface{} constructs of the "data" type.
func newIteratorResults(data interface{}, cfg irConfig) ([]*proto.IteratorResult, error) {
	dataValue := reflect.ValueOf(data)
	switch dataValue.Kind() {
	case reflect.Slice:
		return sliceToIR(dataValue, cfg)
	case reflect.Map, reflect.Struct:
		return mapToIR(data, cfg)
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.String,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, dataValue.Kind())
	}

	return nil, nil
}

// decodeJSON will decode the body of an http.Response object into a slice of IteratorResult. If the response is a
// an array of JSON objects, then each object will be decoded into an IteratorResult.
func decodeJSON(rsp *http.Response) ([]*proto.IteratorResult, error) {
	// Check if the response body is nil, if it is then return an empty slice.
	if rsp.Body == nil {
		return []*proto.IteratorResult{}, nil
	}

	// read the body of the response using the io.Reader interface.
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToReadBody, err)
	}

	// Unmarshal the response body into an interface.
	var data interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSONData, err)
	}

	return nil, nil
}
