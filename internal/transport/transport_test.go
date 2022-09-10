package transport

import (
	"context"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestTimeseries(t *testing.T) {
	t.Run("chunks where end date is before last iteration", func(t *testing.T) {
		ts := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		u, err := url.Parse("https//api.test.com/")
		require.NoError(t, err, "error parsing url: %v", err)

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T00:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		require.NoError(t, err, "error creating timeseries chunks: %v", err)

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 05, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 11, 0, 0, 0, 0, time.UTC),
			},
		}
		require.Equal(t, expChunks, ts.chunks)
	})

	t.Run("chunks where end date is equal to last iteration", func(t *testing.T) {
		ts := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		u, err := url.Parse("https//api.test.com/")
		require.NoError(t, err, "error parsing url: %v", err)

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T01:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		require.NoError(t, err, "error creating timeseries chunks: %v", err)

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 05, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 11, 1, 0, 0, 0, time.UTC),
			},
		}
		require.Equal(t, expChunks, ts.chunks)
	})

	t.Run("chunks where end date is after last iteration", func(t *testing.T) {
		ts := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		u, err := url.Parse("https//api.test.com/")
		require.NoError(t, err, "error parsing url: %v", err)

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T02:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		require.NoError(t, err, "error creating timeseries chunks: %v", err)

		expChunks := [][2]time.Time{
			{
				time.Date(2022, 05, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 5, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 15, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 11, 1, 0, 0, 0, time.UTC),
			},
			{
				time.Date(2022, 05, 11, 1, 0, 0, 0, time.UTC),
				time.Date(2022, 05, 11, 2, 0, 0, 0, time.UTC),
			},
		}
		require.Equal(t, expChunks, ts.chunks)
	})
}

func TestUpsert(t *testing.T) {
	// Create an auth map to fill in the authentication details for the fixture data.
	credentials, err := ioutil.ReadFile("/etc/alpine-hodler/cred.yml")
	if err != nil {
		t.Fatalf("error reading auth config: %v", err)
	}

	auth := make(map[string]Authentication)
	if err := yaml.Unmarshal(credentials, &auth); err != nil {
		t.Fatalf("error unmarshaling auth config: %v", err)
	}

	// Iterate over the fixtures/upsert directory and run each configuration file.
	fixtureRoot := "fixtures/upsert"
	fixtures, err := ioutil.ReadDir(fixtureRoot)
	if err != nil {
		t.Fatalf("error reading fixtures: %v", err)
	}
	for _, fixture := range fixtures {
		t.Run(fixture.Name(), func(t *testing.T) {
			path := filepath.Join(fixtureRoot, fixture.Name())

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("error reading fixture: %v", err)
			}

			var cfg Config
			if err := yaml.Unmarshal(bytes, &cfg); err != nil {
				t.Fatalf("error unmarshaling fixture: %v", err)
			}
			cfg.Logger = logrus.New()

			// Fill in the authentication details for the fixture.
			cfgAuth := cfg.Authentication
			if cfgAuth.APIKey != nil {
				cfg.Authentication = auth[cfgAuth.APIKey.Passphrase]
			}

			// Upsert the fixture.
			if err := Upsert(context.Background(), &cfg); err != nil {
				t.Fatalf("error upserting: %v", err)
			}
		})
	}
}
