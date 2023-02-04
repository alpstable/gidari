package proto

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestDecodeUpsertRequest(t *testing.T) {
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

			// First we create the UpsertRequest object.
			req := &UpsertRequest{
				Data:     tcase.data,
				DataType: int32(tcase.dataType),
			}

			// Then we call the DecodeUpsertRequest function.
			list, err := DecodeUpsertRequest(req)
			if err != tcase.err {
				t.Fatalf("unexpected error: %v", err)
			}

			// Then we check the result.
			if list == nil {
				t.Fatalf("unexpected nil list")
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

	// Create the UpsertRequest object.
	req := &UpsertRequest{
		Data:     data,
		DataType: int32(DecodeTypeJSON),
	}

	// Run the benchmark.
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := DecodeUpsertRequest(req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
