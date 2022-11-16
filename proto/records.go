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
	"net/http"
	"reflect"

	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrUnsupportedDataType     = fmt.Errorf("unsupported data type")
	ErrFailedToAssertInterface = fmt.Errorf("failed to assert interface")
	ErrFailedToMarshalJSON     = fmt.Errorf("failed to marshal json")
	ErrFailedToUnmarshalJSON   = fmt.Errorf("failed to unmarshal json")
	ErrFailedToCreateStruct    = fmt.Errorf("failed to create struct")
	ErrFailedToScanRow         = fmt.Errorf("failed to scan row")
	ErrFailedToParseFloat      = fmt.Errorf("failed to parse float")
	ErrFailedToDecodeRecords   = fmt.Errorf("failed to decode records")
	ErrFailedToGetColumns      = fmt.Errorf("failed to get columns")
)

func newSlice(data interface{}) ([]interface{}, error) {
	var out []interface{}

	dataValue := reflect.ValueOf(data)
	switch dataValue.Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			out = append(out, dataValue.Index(i).Interface())
		}
	case reflect.Map:
		out = append(out, dataValue.Interface())
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.String, reflect.Struct,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, dataValue.Kind())
	}

	return out, nil
}

// decodeRecords will parse a slice of data into a records slice.
func decodeRecords(data interface{}) ([]*structpb.Struct, error) {
	out, err := newSlice(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToCreateStruct, err)
	}

	records := make([]*structpb.Struct, 0)

	for _, r := range out {
		record, err := json.Marshal(r)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
		}

		rec := new(structpb.Struct)

		err = rec.UnmarshalJSON(record)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
		}

		records = append(records, rec)
	}

	return records, nil
}

// UpsertDataType are the supported types for decoding upsert records.
type UpsertDataType uint8

const (
	// UpsertDataJSON is the default upsert data type.
	UpsertDataJSON UpsertDataType = iota
)

// DecodeUpsertRequest will decode the records from the upsert request into a slice of structs.
func DecodeUpsertRequest(req *UpsertRequest) ([]*structpb.Struct, error) {
	if UpsertDataType(req.DataType) == UpsertDataJSON {
		var data interface{}
		if err := json.Unmarshal(req.Data, &data); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
		}

		records, err := decodeRecords(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToDecodeRecords, err)
		}

		return records, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, req.DataType)
}

// DecodeUpsertBinaryRequest will decode the records from an upsert binary request into a slice of structs.
func DecodeUpsertBinaryRequest(req *UpsertBinaryRequest) ([]*structpb.Struct, error) {
	// Create an "UpsertRequest" and decode the records normally.
	upsertReq := &UpsertRequest{
		DataType: int32(UpsertDataJSON),
		Data:     req.Data,
	}

	records, err := DecodeUpsertRequest(upsertReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedToDecodeRecords, err)
	}

	// Get the map for each record and encode it as JSON binary.
	binRecords := make([]*structpb.Struct, len(records))

	for idx, record := range records {
		binRecordData, err := json.Marshal(record.AsMap())
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
		}

		// Create the binary record and add it to the slice.
		binRecords[idx] = &structpb.Struct{
			Fields: map[string]*structpb.Value{
				req.BinaryColumn: {
					Kind: &structpb.Value_StringValue{
						StringValue: string(binRecordData),
					},
				},
			},
		}

		// Now add the primary keys to the binary record.
		if req.PrimaryKeyMap != nil {
			for jsonCol, primaryKey := range req.PrimaryKeyMap {
				if val, ok := record.Fields[jsonCol]; ok {
					binRecords[idx].Fields[primaryKey] = val
				}
			}
		} else {
			binRecords[idx].Fields["id"] = record.Fields["id"]
		}
	}

	return binRecords, nil
}

// DecodeWebResult
func DecodeWebResponse(resp *http.Response) ([]*IteratorResult, error) {
	// determine which format the response is in

	//var data interface{}
	//if err := json.Unmarshal(jsonBytes, &data); err != nil {
	//	return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
	//}

	//out, err := newSlice(data)
	//if err != nil {
	//	return nil, fmt.Errorf("%w: %v", ErrFailedToCreateStruct, err)
	//}

	//results := make([]*IteratorResult, len(out))

	//for idx, r := range out {
	//	record, err := json.Marshal(r)
	//	if err != nil {
	//		return nil, fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
	//	}

	//	results[idx] = &IteratorResult{
	//		Data: record,
	//		URL:  req.URL.String(),
	//	}
	//}

	//return results, nil
}

// PartitionStructs ensures that the request structures are partitioned into size n or less-sized chunks of data, to
// comply with insert requirements.
func PartitionStructs(size int, slice []*structpb.Struct) [][]*structpb.Struct {
	var chunks [][]*structpb.Struct

	for len(slice) > 0 {
		if len(slice) < size {
			size = len(slice)
		}

		chunks = append(chunks, slice[0:size])

		slice = slice[size:]
	}

	return chunks
}
