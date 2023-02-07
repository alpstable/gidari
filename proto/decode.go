// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"encoding/json"
	"fmt"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

var ErrUnsupportedDecodeType = fmt.Errorf("unsupported decode type")

type DecodeType int32

const (
	DecodeTypeUnknown DecodeType = iota
	DecodeTypeJSON
)

func decodeJSON(data []byte) (*structpb.ListValue, error) {
	// If there is no data, return an empty list.
	if len(data) == 0 {
		return &structpb.ListValue{}, nil
	}

	// Check if the first byte of the json is a '{' or '['
	if data[0] == '{' {
		// Unmarshal the json into a structpb.Struct
		record := &structpb.Struct{}
		if err := json.Unmarshal(data, record); err != nil {
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
	if err := json.Unmarshal(data, records); err != nil {
		panic(err)
	}

	return records, nil
}

// Decode will a UpsertRequest into a structpb.ListValue for
// ease-of-use. This method will return an error if the provided "decodeType" is
// not supported.
func Decode(dtype DecodeType, data []byte) (*structpb.ListValue, error) {
	switch dtype {
	case DecodeTypeJSON:
		return decodeJSON(data)
	case DecodeTypeUnknown:
		fallthrough
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedDecodeType, dtype)
	}
}
