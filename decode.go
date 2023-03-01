// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0

package gidari

import (
	"encoding/json"
	"fmt"
	"net/http"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

// ErrUnsupportedDecodeType is returned when the provided decode type is not
// supported.
var ErrUnsupportedDecodeType = fmt.Errorf("unsupported decode type")

// ErrUnsupportedProtoType is returned when the provided proto type is not
// supported.
var ErrUnsupportedProtoType = fmt.Errorf("unsupported proto type")

// DecodeType is an enum that represents the type of data that is being decoded.
type DecodeType int32

const (
	// DecodeTypeUnknown is the default value for the DecodeType enum.
	DecodeTypeUnknown DecodeType = iota

	// DecodeTypeJSON is used to decode JSON data.
	DecodeTypeJSON
)

func addValue(list *structpb.ListValue, val *structpb.Value) error {
	switch val.Kind.(type) {
	case *structpb.Value_StructValue:
		list.Values = append(list.Values, val)
	case *structpb.Value_ListValue:
		list.Values = append(list.Values, val.GetListValue().Values...)
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedProtoType, val.Kind)
	}

	return nil
}

// DecodeFunc is a function that will decode the results of a request into a
// the target.
type DecodeFunc func(list *structpb.ListValue) error

func decodeFuncJSON(rsp *http.Response) DecodeFunc {
	return func(list *structpb.ListValue) error {
		defer func() {
			if err := rsp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		// Check to see if the response is empty. If it is, then we can
		// just return nil.
		if rsp.ContentLength == 0 {
			return nil
		}

		// Decode the response into a list of values.
		dec := json.NewDecoder(rsp.Body)

		for dec.More() {
			val := &structpb.Value{}
			if err := dec.Decode(val); err != nil {
				return fmt.Errorf("failed to decode json: %w", err)
			}

			if err := addValue(list, val); err != nil {
				return fmt.Errorf("failed to add value to list: %w", err)
			}
		}

		return nil
	}
}
