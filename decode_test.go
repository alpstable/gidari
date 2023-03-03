// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0

package gidari

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name            string
		data            []byte
		dataType        DecodeType
		expectedResults []interface{}
		err             error
	}{
		{
			name:     "empty data",
			dataType: DecodeTypeJSON,
		},
		{
			name:     "json object",
			data:     []byte(`{"foo": "bar"}`),
			dataType: DecodeTypeJSON,
			expectedResults: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		{
			name:     "json array",
			data:     []byte(`[{"foo": "bar"}]`),
			dataType: DecodeTypeJSON,
			expectedResults: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		{
			name:     "json array with multiple objects",
			data:     []byte(`[{"foo": "bar"}, {"foo": "baz"}]`),
			dataType: DecodeTypeJSON,
			expectedResults: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"foo": "baz",
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			var decFunc DecodeFunc

			switch tcase.dataType {
			case DecodeTypeJSON:
				decFunc = decodeFuncJSON(&http.Response{
					Body:          io.NopCloser(bytes.NewReader(tcase.data)),
					ContentLength: int64(len(tcase.data)),
				})
			case DecodeTypeUnknown:
				fallthrough
			default:
				t.Fatalf("unsupported decode type: %v", tcase.dataType)
			}

			// Then we call the DecodeUpsertRequest function.
			list := &structpb.ListValue{}
			err := decFunc(list)
			if !errors.Is(err, tcase.err) {
				t.Fatalf("unexpected error: %v", err)
			}

			// Convert the expected result into a list.
			expectedList, err := structpb.NewList(tcase.expectedResults)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Compare the expected list with the actual list.
			if !proto.Equal(expectedList, list) {
				t.Fatalf("unexpected list: %v", list)
			}
		})
	}
}

func BenchmarkDecodeUpsertRequest(b *testing.B) {
	// Create a very large JSON object.
	data := []byte(`{`)
	for i := 0; i < 1000; i++ {
		data = append(data, []byte(fmt.Sprintf(`"foo%d": "bar%d",`, i, i))...)
	}

	data = append(data, []byte(`"foo1000": "bar1000"}`)...)

	httpResponse := func() *http.Response {
		return &http.Response{
			Body:          io.NopCloser(bytes.NewReader(data)),
			ContentLength: int64(len(data)),
		}
	}

	// Run the benchmark.
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := decodeFuncJSON(httpResponse())(&structpb.ListValue{})
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}
