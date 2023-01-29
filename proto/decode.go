package proto

import (
	"fmt"

	"github.com/alpstable/gidari/third_party/accept"
	"google.golang.org/protobuf/encoding/protojson"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

type DecodeType int32

const (
	DecodeTypeUnknown DecodeType = iota
	DecodeTypeJSON
)

// isDecodeTypeJSON will check if the provided "accept" struct is typed for
// decoding into JSON.
func isDecodeTypeJSON(accept accept.Accept) bool {
	return accept.Typ == "application" && (accept.Subtype == "json" || accept.Subtype == "*") ||
		accept.Typ == "*" && accept.Subtype == "*"
}

// BestFitDecodeType will parse the provided Accept(-Charset|-Encoding|-Language)
// header and return the header that best fits the decoding algorithm. If the
// "Accept" header is not set, then this method will return a decodeTypeJSON.
// If the "Accept" header is set, but no match is found, then this method will
// return a decodeTypeUnkown.
//
// See the "acceptSlice.Less" method in the "third_party/accept" package for
// more informaiton on how the "best fit" is determined.
func BestFitDecodeType(header string) DecodeType {
	decodeType := DecodeTypeUnknown
	for _, accept := range accept.ParseAcceptHeader(header) {
		if isDecodeTypeJSON(accept) {
			decodeType = DecodeTypeJSON

			break
		}
	}

	return decodeType
}

func decodeJSON(data []byte) (*structpb.ListValue, error) {
	// Check if the first byte of the json is a '{' or '['
	if data[0] == '{' {
		// Unmarshal the json into a structpb.Struct
		record := &structpb.Struct{}
		if err := protojson.Unmarshal(data, record); err != nil {
			panic(err)
		}

		return &structpb.ListValue{
			Values: []*structpb.Value{
				{
					Kind: &structpb.Value_StructValue{
						StructValue: record,
					},
				},
			},
		}, nil
	}

	records := &structpb.ListValue{}
	if err := protojson.Unmarshal(data, records); err != nil {
		panic(err)
	}

	return records, nil
}

func DecodeUpsertRequest(req *UpsertRequest) (*structpb.ListValue, error) {
	switch DecodeType(req.DataType) {
	case DecodeTypeJSON:
		return decodeJSON(req.Data)
	default:
		return nil, fmt.Errorf("unsupported data type: %d", req.DataType)
	}
}
