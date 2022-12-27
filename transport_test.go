package gidari

import (
	"context"
	"errors"
	"testing"
)

func TestTransport(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name string
		cfg  *Config
		err  error
	}{
		{
			name: "nil",
			err:  ErrNilConfig,
		},
	} {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			err := Transport(context.Background(), tcase.cfg)
			if tcase.err != nil && !errors.Is(err, tcase.err) {
				t.Errorf("expected error %v, got %v", tcase.err, err)
			}

			if tcase.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
