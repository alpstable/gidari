package populate

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/data/repository"
	"github.com/alpine-hodler/driver/transport/option"
	"github.com/alpine-hodler/driver/web/coinbasepro"
)

// CoinbaseProCandlesDirector is a builder to contstruct requests on populating storage repositories with data from the
// Coinbase Pro web API concerning candles.
type CoinbaseProCandlesDirector struct {
	opts *option.CoinbaseProCandles
}

// NewCoinbaseProCandlesDirector will return a new director object for populating storage.
func NewCoinbaseProCandlesDirector(options ...func(*option.CoinbaseProCandles)) CoinbaseProCandlesDirector {
	c := option.NewCoinbaseProCandles()
	for _, o := range options {
		o(c)
	}
	d := new(CoinbaseProCandlesDirector)
	d.opts = c
	return *d
}

// cbpcClient will return a valid web API to populate storage with.
func cbpcClient(ctx context.Context, opts option.CoinbaseProCandles) (*coinbasepro.Client, error) {
	return coinbasepro.NewClient(ctx, opts.RoundTripper)
}

// cbpcEndtime will return the end time for populating historical data. If this configuration hasn't be set, it will
// default to time.Now().
func cbpcEndtime(opts option.CoinbaseProCandles) time.Time {
	if opts.End.IsZero() {
		return time.Now()
	}
	return opts.End
}

// cbpcStarttime will return the start time for populating historical data. If this configuration hasn't be set, it will
// default to five years from the current time (i.e. time.Now()).
func cbpcStarttime(opts option.CoinbaseProCandles) time.Time {
	if opts.Start.IsZero() {
		return time.Now().AddDate(-5, 0, 0)
	}
	return opts.Start
}

// cbpcTimechunks will partition the starttime until time.Now() into period-sized chunks. The period is in terms of
// seconds.
func cbpcTimechunks(period int64, opts option.CoinbaseProCandles) [][2]time.Time {
	chunks := [][2]time.Time{}
	end := cbpcEndtime(opts)
	t := cbpcStarttime(opts)
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

// Populate Will populate all of the relevant repositories with filters over the option data.
func (d CoinbaseProCandlesDirector) Populate(ctx context.Context) error {
	for _, builder := range d.opts.Builders {
		if err := builder.Populate(ctx, *d.opts); err != nil {
			return err
		}
	}
	return nil
}

// candlesWorker fetches candle data from the Coinbase Pro web API.
func candlesWorker(id int, productID string, client *coinbasepro.Client, jobs <-chan *coinbasepro.CandlesOptions,
	results chan<- *coinbasepro.Candles) {
	// range over the options channel and make requests to the candles web API. In theory, we would expect 300 results to
	// be returned per request, but there is no actual guaruntee of that; we can only limit the rate.
	for options := range jobs {
		candles, err := client.Candles(productID, options)
		if err != nil {
			errs := fmt.Errorf("error fetching candles from web api (%v, %v, %v): %v", *options.Start, *options.End,
				options.Granularity, err)
			panic(errs)
		}
		results <- candles
	}
}

// cbpcRepositoryWorkers conrrecntly schedules on job per repository, flushing the data from the buffer.
func cbpcRepositoryWorkers(logHeader string, repos []repository.CoinbasePro, buf <-chan coinbasepro.Candles,
	done chan<- bool) {
	for b := range buf {
		wg := sync.WaitGroup{}
		wg.Add(len(repos))
		for _, repo := range repos {
			go func(logHeader string, repo repository.CoinbasePro) {
				lh := fmt.Sprintf("%s (%s):", logHeader, repo.Name())

				rsp := new(proto.CreateResponse)
				_, err := repo.UpsertCandleMinutes(context.TODO(), b, rsp)
				if err != nil {
					panic(err)
				}
				wg.Done()
				logger.Sugar().Infof("%s buffer flushed for %q", lh, repo.Name())
			}(logHeader, repo)
		}
		wg.Wait()
		done <- true
	}
}

// cbpcProductWorker concurrently flushes candle data from the Coinbase Pro web API to the repositories listed in
// option.CoinbaseProCandles.
func cbpcProductWorker(ctx context.Context, logHeader string, granularity coinbasepro.Granularity, productID string,
	popch chan<- bool, opts option.CoinbaseProCandles) error {

	logHeader = fmt.Sprintf("%s (%s):", logHeader, productID)
	granularityInt, err := strconv.ParseInt(granularity.String(), 0, 32)
	if err != nil {
		panic(err)
	}
	logger.Sugar().Infof("%s populating cb pro candles for %q (%v second intervals)",
		logHeader, productID, granularity)

	client, err := cbpcClient(ctx, opts)
	if err != nil {
		panic(err)
	}

	// period is the maximum amount of time we can query the web API before we exceed the limi.  Coinbase Pro  only
	// returns 300 results at a time, exceeding this will cause a 500 status.
	period := granularityInt * 300

	//  to partition our date range into period-sized chunks to run in parallel.
	partitions := cbpcTimechunks(period, opts)
	logger.Sugar().Debugf("%s batching %v partitions over period %vs", logHeader, len(partitions), period)

	jobs := make(chan *coinbasepro.CandlesOptions, len(partitions))
	results := make(chan *coinbasepro.Candles, len(partitions))

	// Start the same number of workers as cores on the machine.
	for w := 1; w <= runtime.NumCPU(); w++ {
		go candlesWorker(w, productID, client, jobs, results)
	}

	// Send partitions to the web API worker.
	for _, partition := range partitions {
		jobs <- new(coinbasepro.CandlesOptions).
			SetGranularity(granularity).
			SetStart(partition[0].Format(time.RFC3339)).
			SetEnd(partition[1].Format(time.RFC3339))
	}
	close(jobs)

	tolerance := 1500 // TODO: figure out how to abstract tolerance.
	flushcount := int(math.Ceil(float64(len(partitions)) * 300.0 / float64(tolerance)))
	logger.Sugar().Debugf("%s expecting storage to populate with %v flushes (tolerance=%v)",
		logHeader, flushcount, tolerance)

	bufch := make(chan coinbasepro.Candles, flushcount)
	done := make(chan bool, flushcount)

	// Start workers to handle flushing the data from memory into the storage device defined by the repo passed into
	// the bulder.
	for w := 1; w <= flushcount; w++ {
		go cbpcRepositoryWorkers(logHeader, opts.Repositories, bufch, done)
	}

	buf := coinbasepro.Candles{} // Wait for the results to come in, writing them to storage as a buffer fills up.
	flushed := 0                 // keep track of the number of flushed data
	for a := 1; a <= len(partitions); a++ {
		for _, candle := range *<-results {
			candle.ProductID = productID
			buf = append(buf, candle)
		}
		logger.Sugar().Debugf("%s partition %v/%v with buffer size: %v ", logHeader, a, len(partitions), len(buf))

		var reset bool
		if a == len(partitions) || len(buf) >= tolerance {
			// Flush to DB if we are on the last iteration OR the buffer has exceeded the tolerance.
			flushed++
			bufch <- buf
			buf = coinbasepro.Candles{}
			reset = true
			logger.Sugar().Debugf("%s flushed %v/%v", logHeader, flushed, flushcount)
		}
		if len(buf) == 0 && (!reset || a == len(partitions)) {
			// It's possible that we get to this point without the "done" buffer having been completely saturated. This
			// will happen when there are consistently < 300 results from a single request. Since these can build up
			// non-deterministically, we need to completely empty the channel here. That is if the buf is 0 and we have
			// hit the end of the loop OR the buffer has not been reset via the storage routine, then attempt to saturate
			// the "done" buffer.
			logger.Sugar().Debugf("%s forcing saturation", logHeader)
			done <- true
		}
	}
	close(bufch)

	// Wait for all of the candles to flush
	for a := 1; a <= flushcount; a++ {
		<-done
	}
	popch <- true
	logger.Sugar().Infof("%s populating cb pro candles product %q completed", logHeader, productID)
	return nil
}

// cbpcPopulate concurrently schedules one job per productID to flush from the Coinbase Pro Web API to the repositories
// listed in option.CoinbaseProCandles.
func cbpcPopulate(ctx context.Context, granularity coinbasepro.Granularity, opts option.CoinbaseProCandles) error {
	defer logger.Sync()
	id := time.Now().Unix()
	logHeader := fmt.Sprintf("%s %v:", populateCoinbaseProCandlesFn, id)

	done := make(chan bool, len(opts.Products))
	for _, productID := range opts.Products {
		go cbpcProductWorker(ctx, logHeader, granularity, productID, done, opts)
	}

	// Wait for all of the request data to go flush with the repo slice.
	for a := 1; a <= len(opts.Products); a++ {
		<-done
	}

	logger.Sugar().Infof("%s populating coinbase pro candles completed", logHeader)
	return nil
}

type CoinbaseProCandles60 struct{}

// populate will populate historical candle data.
func (c60 *CoinbaseProCandles60) Populate(ctx context.Context, opts option.CoinbaseProCandles) error {
	return cbpcPopulate(ctx, coinbasepro.Granularity60, opts)
}
