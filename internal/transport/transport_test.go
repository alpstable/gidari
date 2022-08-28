package transport

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

		chunks, err := ts.chunks(u)
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
		require.Equal(t, expChunks, chunks)
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

		chunks, err := ts.chunks(u)
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
		require.Equal(t, expChunks, chunks)
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

		chunks, err := ts.chunks(u)
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
		require.Equal(t, expChunks, chunks)
	})
}
