package coinbasepro

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/alpine-hodler/web/pkg/transport"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestStress(t *testing.T) {
	return // ! don't run these atm

	godotenv.Load(".simple-test.env")
	os.Setenv("CB_PRO_URL", "https://api-public.sandbox.exchange.coinbase.com") // safety check

	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	client, err := NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	require.NoError(t, err)

	t.Run("Client.Candles intensive looping", func(t *testing.T) {
		t0 := time.Date(2019, 01, 01, 0, 0, 0, 0, time.UTC)
		for t0.Before(time.Now()) {
			next := t0.AddDate(0, 0, 75)
			opts := new(CandlesOptions).
				SetGranularity(Granularity21600).
				SetStart(t0.Format(time.RFC3339)).
				SetEnd(next.Format(time.RFC3339))

			t0 = next
			_, err := client.Candles("BTC-USD", opts)
			require.NoError(t, err)
		}
	})
	t.Run("Client.Candles concurrent intensive looping", func(t *testing.T) {
		var worker = func(id int,
			client *Client,
			jobs <-chan *CandlesOptions,
			results chan<- *Candles) {
			for options := range jobs {
				fmt.Println("range", *options.Start, *options.End)
				candles, err := client.Candles("BTC-USD", options)
				if err != nil {
					err = fmt.Errorf("error fetching candles from web api (%v, %v, %v): %v",
						*options.Start, *options.End, options.Granularity, err)
					panic(err)
				}
				results <- candles
			}
		}

		chunks := timechunks(18000,
			time.Date(2022, 05, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2022, 05, 30, 0, 0, 0, 0, time.UTC))

		jobs := make(chan *CandlesOptions, len(chunks))
		results := make(chan *Candles, len(chunks))

		for w := 1; w <= runtime.NumCPU(); w++ {
			go worker(w, client, jobs, results)
		}

		for _, chunk := range chunks {
			jobs <- new(CandlesOptions).
				SetGranularity(Granularity60).
				SetStart(chunk[0].Format(time.RFC3339)).
				SetEnd(chunk[1].Format(time.RFC3339))
		}
		close(jobs)
		for a := 1; a <= len(chunks); a++ {
			<-results
		}
	})
}

func timechunks(period int64, start, end time.Time) [][2]time.Time {
	chunks := [][2]time.Time{}
	t := start
	for t.Before(end) {
		next := t.Add(time.Second * time.Duration(period))
		if next.Before(end) {
			chunks = append(chunks, [2]time.Time{t, next})
		} else {
			chunks = append(chunks, [2]time.Time{t, end})
		}
		t = next
	}
	return chunks
}
