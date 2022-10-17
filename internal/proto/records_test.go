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
					Table:        "tests1",
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
		t.Run(tcase.name, func(t *testing.T) {
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
