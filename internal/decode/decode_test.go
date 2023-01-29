package decode

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/alpstable/gidari/proto"
)

func TestBestFitDecodeType(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name   string
		accept string
		want   decodeType
	}{
		{
			name:   "empty",
			accept: "",
			want:   decodeTypeJSON,
		},
		{
			name:   "json",
			accept: "application/json",
			want:   decodeTypeJSON,
		},
		{
			name:   "json+protobuf",
			accept: "application/json, application/vnd.google.protobuf",
			want:   decodeTypeJSON,
		},
		{
			name:   "protobuf",
			accept: "application/vnd.google.protobuf",
			want:   decodeTypeUnknown,
		},
		{
			name:   "protobuf+json",
			accept: "application/vnd.google.protobuf, application/json",
			want:   decodeTypeJSON,
		},
		{
			name:   "protobuf+json+protobuf",
			accept: "application/vnd.google.protobuf, application/json, application/vnd.google.protobuf",
			want:   decodeTypeJSON,
		},
		{
			name:   "protobuf+json+qualityfactor",
			accept: "application/vnd.google.protobuf, application/json;q=0.5",
			want:   decodeTypeJSON,
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			if got := bestFitDecodeType(tcase.accept); got != tcase.want {
				t.Errorf("bestFitDecodeType(%q) = %v, want %v", tcase.accept, got, tcase.want)
			}
		})
	}
}

func TestMapToIR(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name   string
		mapArg map[string]interface{}
		cfg    irConfig
		want   []*proto.IteratorResult
	}{
		{
			name:   "empty",
			mapArg: map[string]interface{}{},
			cfg: irConfig{
				enc: func(interface{}) ([]byte, error) { return nil, nil },
				uri: func() *url.URL {
					u, _ := url.Parse("http://localhost:8080")

					return u
				}(),
			},
			want: []*proto.IteratorResult{},
		},
		{
			name: "json",
			mapArg: map[string]interface{}{
				"foo": "bar",
			},
			cfg: irConfig{
				enc: json.Marshal,
				uri: func() *url.URL {
					u, _ := url.Parse("http://localhost:8080")

					return u
				}(),
			},
			want: []*proto.IteratorResult{
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":"bar"}`),
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := mapToIR(tcase.mapArg, tcase.cfg)
			if err != nil {
				t.Fatalf("mapToIR() = %v", err)
			}

			if len(got) != len(tcase.want) {
				t.Fatalf("mapToIR() = %v, want %v", got, tcase.want)
			}

			for i := range got {
				if !reflect.DeepEqual(got[i].Data, tcase.want[i].Data) {
					t.Errorf("mapToIR() = %v, want %v", got, tcase.want)
				}

				if got[i].URL != tcase.want[i].URL {
					t.Errorf("mapToIR() = %v, want %v", got, tcase.want)
				}
			}
		})
	}
}

func TestSliceToIR(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name     string
		sliceArg interface{}
		cfg      irConfig
		want     []*proto.IteratorResult
	}{
		{
			name:     "empty",
			sliceArg: []interface{}{},
			cfg: irConfig{
				enc: func(interface{}) ([]byte, error) { return nil, nil },
				uri: func() *url.URL {
					u, _ := url.Parse("http://localhost:8080")

					return u
				}(),
			},
		},
		{
			name: "json",
			sliceArg: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"foo": "baz",
				},
			},
			cfg: irConfig{
				enc: json.Marshal,
				uri: func() *url.URL {
					u, _ := url.Parse("http://localhost:8080")

					return u
				}(),
			},
			want: []*proto.IteratorResult{
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":"bar"}`),
				},
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":"baz"}`),
				},
			},
		},
		{
			name: "json embedded slices",
			sliceArg: []interface{}{
				map[string]interface{}{
					"foo": []interface{}{
						map[string]interface{}{
							"bar": "baz",
						},
					},
				},
				[]interface{}{
					map[string]interface{}{
						"foo": "bar",
					},
					map[string]interface{}{
						"foo": "baz",
					},
				},
			},
			cfg: irConfig{
				enc: json.Marshal,
				uri: func() *url.URL {
					u, _ := url.Parse("http://localhost:8080")

					return u
				}(),
			},
			want: []*proto.IteratorResult{
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":[{"bar":"baz"}]}`),
				},
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":"bar"}`),
				},
				{
					URL:  "http://localhost:8080",
					Data: []byte(`{"foo":"baz"}`),
				},
			},
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := sliceToIR(reflect.ValueOf(tcase.sliceArg), tcase.cfg)
			if err != nil {
				t.Fatalf("sliceToIR() = %v", err)
			}

			if len(got) != len(tcase.want) {
				t.Fatalf("sliceToIR() = %v, want %v", got, tcase.want)
			}

			for i := range got {
				if !reflect.DeepEqual(got[i].Data, tcase.want[i].Data) {
					t.Errorf("sliceToIR() = %v, want %v", got, tcase.want)
				}

				if got[i].URL != tcase.want[i].URL {
					t.Errorf("sliceToIR() = %v, want %v", got, tcase.want)
				}
			}
		})
	}
}

// func TestNewInterfaceSlice(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name string
//		in   interface{}
//		want []interface{}
//		err  error
//	}{
//		{
//			name: "empty",
//			in:   []string{},
//			want: []interface{}{},
//			err:  nil,
//		},
//		{
//			name: "struct",
//			in: struct {
//				Test string
//			}{
//				Test: "test",
//			},
//			want: []interface{}{
//				struct {
//					Test string
//				}{
//					Test: "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "struct slice",
//			in: []struct {
//				Test string
//			}{
//				{
//					Test: "test",
//				},
//			},
//			want: []interface{}{
//				struct {
//					Test string
//				}{
//					Test: "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "struct slice ptr",
//			in: []*struct {
//				Test string
//			}{
//				{
//					Test: "test",
//				},
//			},
//			want: []interface{}{
//				&struct {
//					Test string
//				}{
//					Test: "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "map",
//			in: map[string]string{
//				"test": "test",
//			},
//			want: []interface{}{
//				map[string]string{
//					"test": "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "map slice",
//			in: []map[string]string{
//				{
//					"test": "test",
//				},
//			},
//			want: []interface{}{
//				map[string]string{
//					"test": "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "map slice ptr",
//			in: []*map[string]string{
//				{
//					"test": "test",
//				},
//			},
//			want: []interface{}{
//				&map[string]string{
//					"test": "test",
//				},
//			},
//			err: nil,
//		},
//		{
//			name: "invalid",
//			in:   "test",
//			want: nil,
//			err:  ErrUnsupportedDataType,
//		},
//	} {
//		tcase := tcase
//
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, err := makeIteratorResultSlice(tcase.in)
//			if !errors.Is(err, tcase.err) {
//				t.Errorf("newInterfaceSlice(%v) = %v, want %v", tcase.in, err, tcase.err)
//			}
//
//			if len(got) != len(tcase.want) {
//				t.Errorf("newInterfaceSlice(%v) = %v, want %v", tcase.in, got, tcase.want)
//			}
//
//			sanitizedWant := make([]interface{}, len(tcase.want))
//			for i, v := range tcase.want {
//				sanitizedWant[i] = reflect.ValueOf(v).Interface()
//			}
//
//			sanitizedGot := make([]interface{}, len(got))
//			for i, v := range got {
//				sanitizedGot[i] = reflect.ValueOf(v).Interface()
//			}
//
//			if !reflect.DeepEqual(sanitizedGot, sanitizedWant) {
//				t.Errorf("newInterfaceSlice(%v) = %v, want %v", tcase.in, sanitizedGot, sanitizedWant)
//			}
//		})
//	}
//}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		rsp  *http.Response
		want []*proto.IteratorResult
		err  error
	}{
		{
			name: "empty",
			rsp: &http.Response{
				Body:       nil,
				StatusCode: http.StatusOK,
			},
			want: []*proto.IteratorResult{},
			err:  nil,
		},
		//{
		//	name: "json object",
		//	rsp: &http.Response{
		//		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"test":"test"}`))),
		//		StatusCode: http.StatusOK,
		//		Request: &http.Request{
		//			URL: func() *url.URL {
		//				u, _ := url.Parse("http://localhost:8080")

		//				return u
		//			}(),
		//		},
		//	},
		//	want: []*proto.IteratorResult{
		//		{
		//			Data: []byte(`{"test":"test"}`),
		//			URL:  "http://localhost:8080",
		//		},
		//	},
		// },
	} {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := decodeJSON(tcase.rsp)
			if err != tcase.err {
				t.Fatalf("got error %v, want %v", err, tcase.err)
			}

			if tcase.want == nil && got != nil {
				t.Fatalf("got %v, want %v", got, tcase.want)
			}

			if len(got) != len(tcase.want) {
				t.Fatalf("got %d results, want %d", len(got), len(tcase.want))
			}

			for i := range got {
				for j := range got[i].Data {
					if got[i].Data[j] != tcase.want[i].Data[j] {
						t.Fatalf("got %v, want %v", got[i].Data[j], tcase.want[i].Data[j])
					}
				}

				if got[i].URL != tcase.want[i].URL {
					t.Fatalf("got %v, want %v", got[i].URL, tcase.want[i].URL)
				}
			}
		})
	}
}

// func BenchmarkNewInterfaceSlice(b *testing.B) {
//	for _, tcase := range []struct {
//		name string
//		in   interface{}
//	}{
//		{
//			name: "slice",
//			in: []struct {
//				Test string `json:"test"`
//			}{
//				{
//					Test: "test",
//				},
//				{
//					Test: "test",
//				},
//			},
//		},
//	} {
//		b.Run(tcase.name, func(b *testing.B) {
//			b.ReportAllocs()
//
//			for i := 0; i < b.N; i++ {
//				_, err := makeIteratorResultSlice(tcase.in)
//				if err != nil {
//					b.Fatal(err)
//				}
//			}
//		})
//	}
//}
