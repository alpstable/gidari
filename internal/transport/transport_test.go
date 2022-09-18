package transport

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T00:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

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
		if !reflect.DeepEqual(expChunks, ts.chunks) {
			t.Fatalf("unexpected chunks: %v", ts.chunks)
		}
	})

	t.Run("chunks where end date is equal to last iteration", func(t *testing.T) {
		ts := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		u, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T01:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

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
		if !reflect.DeepEqual(expChunks, ts.chunks) {
			t.Fatalf("unexpected chunks: %v", ts.chunks)
		}
	})

	t.Run("chunks where end date is after last iteration", func(t *testing.T) {
		ts := &timeseries{
			StartName: "start",
			EndName:   "end",
			Period:    18000,
		}

		u, err := url.Parse("https//api.test.com/")
		if err != nil {
			t.Fatalf("error parsing url: %v", err)
		}

		query := u.Query()
		query.Set("start", "2022-05-10T00:00:00Z")
		query.Set("end", "2022-05-11T02:00:00Z")
		u.RawQuery = query.Encode()

		err = ts.setChunks(u)
		if err != nil {
			t.Fatalf("error setting chunks: %v", err)
		}

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
		if !reflect.DeepEqual(expChunks, ts.chunks) {
			t.Fatalf("unexpected chunks: %v", ts.chunks)
		}
	})
}

func TestUpsert(t *testing.T) {
	// if /etc/alpine-hodler/auth.env exists, load it
	if _, err := os.Stat("/etc/alpine-hodler/auth.env"); err == nil {
		err := godotenv.Load("/etc/alpine-hodler/auth.env")
		if err != nil {
			t.Fatalf("error loading auth.env: %v", err)
		}
	} else {
		t.Log("no auth.env found")
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
				// The "passhprase" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = Authentication{
					APIKey: &APIKey{
						Key:        os.Getenv(cfgAuth.APIKey.Key),
						Secret:     os.Getenv(cfgAuth.APIKey.Secret),
						Passphrase: os.Getenv(cfgAuth.APIKey.Passphrase),
					},
				}
			}

			if cfgAuth.Auth2 != nil {
				// The "bearer" field in the fixture should be the name of the auth map entry. That
				// is how we lookup which authentication details to use.
				cfg.Authentication = Authentication{
					Auth2: &Auth2{
						Bearer: os.Getenv(cfgAuth.Auth2.Bearer),
					},
				}
			}

			// Upsert the fixture.
			if err := Upsert(context.Background(), &cfg); err != nil {
				t.Fatalf("error upserting: %v", err)
			}
		})
	}
}
