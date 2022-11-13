package gidari

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
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

func newIteratorConfigWithHandler(t *testing.T, uri string, h HTTPResponseHandler) *Config {
	t.Helper()

	cfg := newIteratorConfig(t, uri)
	cfg.HTTPResponseHandler = h

	return cfg
}

func newIteratorHandler(t *testing.T) HTTPResponseHandler {
	t.Helper()

	return func(ctx context.Context, httpResponse HTTPResponse) ([]*proto.IteratorResult, error) {
		return []*proto.IteratorResult{
			{
				URL:  "http://localhost",
				Data: []byte("test"),
			},
		}, nil
	}
}

func TestNewIterator(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		cfg  *Config
		want *Iterator
		err  error
	}{
		{
			name: "nil config",
			cfg:  nil,
			err:  ErrNilConfig,
		},
		{
			name: "empty config",
			cfg:  &Config{},
			err:  ErrMissingURL,
		},
		{
			name: "default HTTP Handler",
			cfg:  newIteratorConfig(t, "http://localhost"),
			want: &Iterator{
				cfg:                 newIteratorConfig(t, "http://localhost"),
				httpResponseHandler: defaultHTTPResponseHandler,
			},
		},
		{
			name: "custom HTTP Handler",
			cfg:  newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),
			want: &Iterator{
				cfg: newIteratorConfigWithHandler(t, "http://localhost", newIteratorHandler(t)),

				httpResponseHandler: newIteratorHandler(t),
			},
		},
	} {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewIterator(context.Background(), tcase.cfg)
			if err != tcase.err {
				t.Fatalf("got error %v, want %v", err, tcase.err)
			}

			if tcase.err != nil {
				return
			}

			httpResponse := HTTPResponse{
				Response: &http.Response{
					Body: ioutil.NopCloser(bytes.NewReader([]byte("test"))),
				},
			}

			gotResults, err := got.httpResponseHandler(context.Background(), httpResponse)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}

			wantResults, err := tcase.want.httpResponseHandler(context.Background(), httpResponse)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}

			if len(gotResults) != len(wantResults) {
				t.Fatalf("got %d results, want %d", len(gotResults), len(wantResults))
			}

			if !reflect.DeepEqual(gotResults, wantResults) {
				t.Fatalf("got %v, want %v", gotResults[0], wantResults[0])
			}
		})
	}
}
