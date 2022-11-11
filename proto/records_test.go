// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package proto

import (
	reflect "reflect"
	"testing"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestDecodeUpsertBinaryRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		req *UpsertBinaryRequest
	}

	testCases := []struct {
		name string
		args args
		want []*structpb.Struct
	}{
		{
			name: "valid",
			args: args{
				req: &UpsertBinaryRequest{
					Table:        &Table{Name: "tests1"},
					BinaryColumn: "test_string",
					PrimaryKeyMap: map[string]string{
						"id": "1",
					},
					Data: []byte(`{"x":1}`),
				},
			},
			want: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"test_string": {
							Kind: &structpb.Value_StringValue{
								StringValue: `{"x":1}`,
							},
						},
					},
				},
			},
		},
	}

	for _, tcase := range testCases {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := DecodeUpsertBinaryRequest(tcase.args.req)
			if err != nil {
				t.Errorf("DecodeUpsertBinaryRequest() error = %v, wantErr %v", err, false)
			}

			if !reflect.DeepEqual(got, tcase.want) {
				t.Errorf("DecodeUpsertBinaryRequest() = %v, want %v", got, tcase.want)
			}
		})
	}
}

func TestDecodeIteratorResults(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name      string
		endpoint  string
		jsonBytes []byte
		want      []*IteratorResult
	}{
		{
			name:      "object",
			endpoint:  "tests1",
			jsonBytes: []byte(`{"x":1}`),
			want: []*IteratorResult{
				{Data: []byte(`{"x":1}`), Endpoint: "tests1"},
			},
		},
		{
			name:      "array",
			endpoint:  "tests1",
			jsonBytes: []byte(`[{"x":1},{"x":2}]`),
			want: []*IteratorResult{
				{Data: []byte(`{"x":1}`), Endpoint: "tests1"},
				{Data: []byte(`{"x":2}`), Endpoint: "tests1"},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := DecodeIteratorResults(tcase.endpoint, tcase.jsonBytes)
			if err != nil {
				t.Errorf("DecodeIteratorResults() error = %v, wantErr %v", err, false)
			}

			if !reflect.DeepEqual(got, tcase.want) {
				t.Errorf("DecodeIteratorResults() = %v, want %v", got, tcase.want)
			}
		})
	}
}
