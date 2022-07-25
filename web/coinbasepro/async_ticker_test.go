package coinbasepro

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/alpine-hodler/web/pkg/websocket"
)

func TestAsyncTickerStream(t *testing.T) {
	t.Run("mock third-party usage", func(t *testing.T) {
		t.Run("should start and close without error", func(t *testing.T) {
			// some third-party goroutines
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			numThirdParties := 2
			wg := sync.WaitGroup{}
			wg.Add(numThirdParties)
			for i := 0; i < numThirdParties; i++ {
				go func() {
					ticker.Open()
					r := 1 + rand.Intn(1)
					time.Sleep(time.Duration(r) * time.Second)
					ticker.Close()
					wg.Done()
				}()
			}
			wg.Wait()
		})
	})

	t.Run("ticker#Close", func(t *testing.T) {
		t.Run("should close without error", func(t *testing.T) {
			treshold := 100
			for i := 0; i < treshold; i++ {
				mockC, _ := websocket.NewMock()
				ticker := newAsyncTicker(mockC)
				ticker.Open()
				go func() {
					tickers := []Ticker{}
					for ticker := range ticker.Channel() {
						tickers = append(tickers, ticker)
					}
				}()
				ticker.Close()
			}
		})

		t.Run("should closew without error on long runtime", func(t *testing.T) {
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			ticker.Open()
			go func() {
				tickers := []Ticker{}
				for ticker := range ticker.Channel() {
					tickers = append(tickers, ticker)
				}
			}()
			time.Sleep(2 * time.Second)
			ticker.Close()
			time.Sleep(2 * time.Millisecond)
		})

		t.Run("should do nothing when there is no stream", func(t *testing.T) {
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			ticker.Close()
		})
	})

	t.Run("ticker#Open", func(t *testing.T) {
		t.Run("should re-initialize channel data after each close", func(t *testing.T) {
			treshold := 100
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			for i := 0; i < treshold; i++ {
				ticker.Open()
				go func() {
					tickers := []Ticker{}
					for ticker := range ticker.Channel() {
						tickers = append(tickers, ticker)
					}
				}()
				ticker.Close()
			}
		})

		t.Run("should be able to start stream over again", func(t *testing.T) {
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			ticker.Open()
			go func() {
				tickers := []Ticker{}
				for ticker := range ticker.Channel() {
					tickers = append(tickers, ticker)
				}
			}()
			time.Sleep(2 * time.Microsecond)
			ticker.Open()
			go func() {
				tickers := []Ticker{}
				for ticker := range ticker.Channel() {
					tickers = append(tickers, ticker)
				}
			}()
			time.Sleep(1 * time.Microsecond)
			ticker.Close()
		})

		t.Run("should not fatal if you start streams concurrently", func(t *testing.T) {
			mockC, _ := websocket.NewMock()
			ticker := newAsyncTicker(mockC)
			treshold := 1000
			for j := 0; j < treshold; j++ {
				concurrentStreams := 100
				for i := 0; i < concurrentStreams; i++ {
					go ticker.Open()
				}
				go func() {
					tickers := []Ticker{}
					for ticker := range ticker.Channel() {
						tickers = append(tickers, ticker)
					}
				}()
				time.Sleep(1 * time.Microsecond)
				ticker.Close()
			}
		})
	})
}
