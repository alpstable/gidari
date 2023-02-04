// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	"encoding/json"
	"fmt"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

var ErrUnsupportedDecodeType = fmt.Errorf("unsupported decode type")

func decodeJSON(data []byte) (*structpb.ListValue, error) {
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

// DecodeUpsertRequest will a UpsertRequest into a structpb.ListValue for
// ease-of-use. This method will return an error if the provided "decodeType" is
// not supported.
func DecodeUpsertRequest(req *UpsertRequest) (*structpb.ListValue, error) {
	switch DecodeType(req.DataType) {
	case DecodeTypeJSON:
		return decodeJSON(req.Data)
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedDecodeType, req.DataType)
	}
}
