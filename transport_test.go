package gidari

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestTransport(t *testing.T) {
	t.Parallel()

	for _, tcase := range []struct {
		name                              string
		expectedNumberOfUpsertsPerStorage int
		cfg                               *Config
		err                               error
	}{
		{
			name: "no configuration",
			err:  ErrNilConfig,
		},
		{
			name: "no requests",
			cfg:  &Config{},
			err:  ErrNoRequests,
		},
		{
			name: "single request with single storage",
			cfg: newMockConfig(mockConfigOptions{
				reqCount: 1,
				stgCount: 1,
			}),
			expectedNumberOfUpsertsPerStorage: 1,
		},
		{
			name: "single request with multiple storages",
			cfg: newMockConfig(mockConfigOptions{
				reqCount: 1,
				stgCount: 3,
			}),
			expectedNumberOfUpsertsPerStorage: 1,
		},
		{
			name: "multiple requests with single storage",
			cfg: newMockConfig(mockConfigOptions{
				reqCount: 3,
				stgCount: 1,
			}),
			expectedNumberOfUpsertsPerStorage: 3,
		},
		{
			name: "multiple requests with multiple storages",
			cfg: newMockConfig(mockConfigOptions{
				reqCount: 3,
				stgCount: 3,
			}),
			expectedNumberOfUpsertsPerStorage: 3,
		},
		{
			name: "voluminous requests with multiple storages",
			cfg: newMockConfig(mockConfigOptions{
				reqCount:    10_000,
				stgCount:    3,
				rateLimiter: rate.NewLimiter(rate.Limit(1*time.Second), 10_000),
			}),
			expectedNumberOfUpsertsPerStorage: 10_000,
		},
	} {
		tcase := tcase

		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			err := Transport(context.Background(), tcase.cfg)
			if tcase.err != nil && !errors.Is(err, tcase.err) {
				t.Errorf("expected error %v, got %v", tcase.err, err)
			}

			if tcase.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			// If there is no configuration then we can termiante
			// the test here.
			if tcase.cfg == nil {
				return
			}

			cfg := tcase.cfg

			// If there is no mock storage then we can terminate
			// the test here.
			if len(cfg.Storage) == 0 {
				return
			}

			// We need to validate various operation for each
			// storage object.
			for _, stg := range cfg.Storage {
				mockStorage, ok := stg.Storage.(*mockStorage)
				if !ok {
					t.Errorf("expected mock storage, got %T", stg)
				}

				// The number of upserts should be equal to the
				// expected number of upserts. Note that there
				// can be less requests than upserts, for
				// example a timeseries request could be broken
				// into multipe flattened requests for upsert.
				if mockStorage.upsertCount != tcase.expectedNumberOfUpsertsPerStorage {
					t.Errorf("expected %d upserts, got %d", tcase.expectedNumberOfUpsertsPerStorage,
						mockStorage.upsertCount)
				}
			}

		})
	}
}
