// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/alpstable/gidari/proto"
)

func newIteratorConfig(t *testing.T, uri string) *Config {
	t.Helper()

	return &Config{
		URL: func() *url.URL {
			u, _ := url.Parse("http://localhost")
			return u
		}(),
	}
}

func newIteratorConfigWithHandler(t *testing.T, uri string, h WebResultAssigner) *Config {
	t.Helper()

	cfg := newIteratorConfig(t, uri)
	cfg.HTTPResponseHandler = h

	return cfg
}

func newIteratorHandler(t *testing.T) WebResultAssigner {
	t.Helper()

	return func(ctx context.Context, httpResponse WebResult) ([]*proto.IteratorResult, error) {
		return []*proto.IteratorResult{
			{
				URL:  "http://localhost",
				Data: []byte("test"),
			},
		}, nil
	}
}

//func TestNewIterator(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name string
//		cfg  *Config
//		want *Iterator
//		err  error
//	}{
//		{
//			name: "nil config",
//			cfg:  nil,
//			err:  ErrNilConfig,
//		},
//		{
//			name: "empty config",
//			cfg:  &Config{},
//			err:  ErrMissingURL,
//		},
//		{
//			name: "default HTTP Handler",
//			cfg:  newIteratorConfig(t, "http://localhost"),
//			want: &Iterator{
//				cfg:               newIteratorConfig(t, "http://localhost"),
//				webResultAssigner: assignWebResult,
//			},
//		},
//		{
//			name: "custom HTTP Handler",
//			cfg:  newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),
//			want: &Iterator{
//				cfg: newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),
//
//				webResultAssigner: newIteratorHandler(t),
//			},
//		},
//	} {
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, err := NewIterator(context.Background(), tcase.cfg)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if tcase.err != nil {
//				return
//			}
//
//			httpResponse := WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("test"))),
//				},
//			}
//
//			gotResults, err := got.webResultAssigner(context.Background(), httpResponse)
//			if err != nil {
//				t.Fatalf("got error %v, want nil", err)
//			}
//
//			wantResults, err := tcase.want.webResultAssigner(context.Background(), httpResponse)
//			if err != nil {
//				t.Fatalf("got error %v, want nil", err)
//			}
//
//			if len(gotResults) != len(wantResults) {
//				t.Fatalf("got %d results, want %d", len(gotResults), len(wantResults))
//			}
//
//			if !reflect.DeepEqual(gotResults, wantResults) {
//				t.Fatalf("got %v, want %v", gotResults[0], wantResults[0])
//			}
//		})
//	}
//}
//
//func TestSanitizeInvalidJSON(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name       string
//		data       []byte
//		clobColumn string
//		want       []byte
//		ok         bool
//		err        error
//	}{
//		{
//			name: "empty",
//			data: []byte{},
//			want: []byte{},
//		},
//		{
//			name: "invalid, no clob column",
//			data: []byte("invalid"),
//			ok:   false,
//		},
//		{
//			name:       "invalid, clob column",
//			data:       []byte("invalid"),
//			clobColumn: "test",
//			want:       []byte(`{"test":"invalid"}`),
//			ok:         true,
//		},
//		{
//			name: "valid, no clob column",
//			data: []byte(`{"test": "test"}`),
//			want: []byte(`{"test": "test"}`),
//			ok:   true,
//		},
//		{
//			name:       "valid, clob column",
//			data:       []byte(`{"test": "test"}`),
//			clobColumn: "doesnt matter",
//			want:       []byte(`{"test": "test"}`),
//			ok:         true,
//		},
//	} {
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, ok, err := sanitizeJSON(tcase.data, tcase.clobColumn)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if ok != tcase.ok {
//				t.Fatalf("got ok %v, want %v", ok, tcase.ok)
//			}
//
//			if !bytes.Equal(got, tcase.want) {
//				t.Fatalf("got %s, want %s", got, tcase.want)
//			}
//		})
//	}
//}
//
//func TestAssignWebResult(t *testing.T) {
//	t.Parallel()
//
//	for _, tcase := range []struct {
//		name string
//		rsp  WebResult
//		err  error
//		want []*proto.IteratorResult
//	}{
//		{
//			name: "nil body",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: nil,
//				},
//			},
//			want: []*proto.IteratorResult{},
//		},
//		{
//			name: "empty body",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte{})),
//				},
//			},
//			want: []*proto.IteratorResult{},
//		},
//		{
//			name: "invalid JSON",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("invalid"))),
//				},
//			},
//			want: nil,
//		},
//		{
//			name: "valid JSON",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"test": "test"}`))),
//				},
//
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "invalid JSON, clob column",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte("invalid"))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//				ClobColumn: "test",
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"invalid"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "valid JSON, clob column",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`{"test": "test"}`))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//				ClobColumn: "doesn't matter",
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//		{
//			name: "valid JSON array",
//			rsp: WebResult{
//				Response: &http.Response{
//					Body: ioutil.NopCloser(bytes.NewReader([]byte(`[{"test": "test"},
//{"test2": "test"}]`))),
//				},
//				URL: func() *url.URL {
//					u, _ := url.Parse("http://localhost")
//
//					return u
//				}(),
//			},
//			want: []*proto.IteratorResult{
//				{
//					Data: []byte(`{"test":"test"}`),
//					URL:  "http://localhost",
//				},
//				{
//					Data: []byte(`{"test2":"test"}`),
//					URL:  "http://localhost",
//				},
//			},
//		},
//	} {
//		tcase := tcase
//
//		t.Run(tcase.name, func(t *testing.T) {
//			t.Parallel()
//
//			got, err := assignWebResult(context.Background(), tcase.rsp)
//			if err != tcase.err {
//				t.Fatalf("got error %v, want %v", err, tcase.err)
//			}
//
//			if len(got) != len(tcase.want) {
//				t.Fatalf("got %d results, want %d", len(got), len(tcase.want))
//			}
//
//			for idx, result := range got {
//				if result.GetURL() != tcase.want[idx].GetURL() {
//					t.Fatalf("got %s, want %s", result.GetURL(), tcase.want[idx].GetURL())
//				}
//
//				if !bytes.Equal(result.GetData(), tcase.want[idx].GetData()) {
//					t.Fatalf("got %s, want %s", result.GetData(), tcase.want[idx].GetData())
//				}
//			}
//		})
//	}
//}

func mockResponseBodyArray(b *testing.B, size int) []byte {
	b.Helper()

	var buf bytes.Buffer

	buf.WriteString("[")

	for i := 0; i < size; i++ {
		buf.WriteString(`{"test": "test"}`)

		if i < size-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("]")

	return buf.Bytes()
}

func mockResponseBodyObject(b *testing.B) []byte {
	b.Helper()

	return []byte(`{"test": "test"}`)
}

func BenchmarkAssignWebResult(b *testing.B) {
	b.ReportAllocs()

	for _, tcase := range []struct {
		name string
		rsp  WebResult
	}{
		{
			name: "valid JSON array",
			rsp: WebResult{
				Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(mockResponseBodyArray(b, 100))),
				},
				URL: func() *url.URL {
					u, _ := url.Parse("http://localhost")

					return u
				}(),
			},
		},
		{
			name: "valid JSON object",
			rsp: WebResult{
				Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader(mockResponseBodyObject(b))),
				},
				URL: func() *url.URL {
					u, _ := url.Parse("http://localhost")

					return u
				}(),
			},
		},
	} {
		tcase := tcase

		b.Run(tcase.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				assignWebResult(context.Background(), tcase.rsp)
			}
		})
	}
}
