package option

import (
	"context"
	"time"

	"github.com/alpine-hodler/driver/data/repository"
	"github.com/alpine-hodler/driver/web/transport"
)

type CoinbaseProCandlesPopulater interface {
	Populate(context.Context, CoinbaseProCandles) error
}

type CoinbaseProCandles struct {
	Builders     []CoinbaseProCandlesPopulater
	Products     []string
	Repositories []repository.CoinbasePro
	RoundTripper transport.T
	Start        time.Time
	End          time.Time
}

func NewCoinbaseProCandles() *CoinbaseProCandles {
	return new(CoinbaseProCandles)
}

// CBPCWithBuilder will include a builder to populate storage.
func CBPCWithPopulater(p CoinbaseProCandlesPopulater) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.Builders = append(opts.Builders, p)
	}
}

// CBPCWithducts will set the coinbase pro repository to populate CBPCWith web API data for a slice of products.
func CBPCWithProduct(product string) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.Products = append(opts.Products, product)
	}
}

// CBPCWithRepository will set the coinbase pro repository to populate CBPCWith web API data.
func CBPCWithRepository(repo repository.CoinbasePro) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.Repositories = append(opts.Repositories, repo)
	}
}

// CBPCWithRoundTripper will set the web API http request rount tripper for authentication.
func CBPCWithRoundTripper(rt transport.T) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.RoundTripper = rt
	}
}

// CBPCWithEndTime will set the end time for populating historical data. If this value is not set, all historical
// builders will presume the end time is the current time (i.e. time.Now()).
func CBPCWithEndTime(end time.Time) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.End = end
	}
}

// CBPCWithStartTime will set the start time for populating historical data. If this value is not set, all historical
// builders will presume the start time is five years from the current time (i.e. time.Now()).
func CBPCWithStartTime(start time.Time) func(*CoinbaseProCandles) {
	return func(opts *CoinbaseProCandles) {
		opts.Start = start
	}
}
