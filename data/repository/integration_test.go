package repository_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/data/repository"
	"github.com/alpine-hodler/driver/data/storage"
	"github.com/alpine-hodler/driver/tools"
	"github.com/alpine-hodler/driver/web/coinbasepro"
	"github.com/alpine-hodler/driver/web/polygon"
	"github.com/stretchr/testify/require"
)

type teststg uint8

const (
	mongoteststg teststg = iota
	postgresstg
)

type testrepo uint8

const (
	coinbaseprotestrepo testrepo = iota
	polygontestrepo
)

type testop func(context.Context, *testing.T, tools.GenericStorage, testrepo)
type testsession func(*testing.T, string, testop) bool

type testconfig struct {
	tr       testrepo
	ts       teststg
	name     string
	port     string
	database string
	host     string
	username string
}

func testexecutor(t *testing.T, config testconfig) []interface{} {
	var ok bool
	var err error

	ctx := context.Background()

	// set the storage object from factory
	var stg tools.GenericStorage
	switch config.ts {
	case mongoteststg:
		uri, _ := tools.MongoURI(config.host, "", "", config.port, config.database)
		stg, err = storage.NewMongo(ctx, uri)
	case postgresstg:
		uri, _ := tools.PostgresURI(config.host, config.username, "", config.port, config.database)
		stg, err = storage.NewPostgres(ctx, uri)
	}
	require.NoError(t, err)

	// Set the repo executor from factory
	var exec func(context.Context, func(context.Context, tools.GenericStorage) (bool, error)) error
	switch config.tr {
	case coinbaseprotestrepo:
		repo := repository.NewCoinbasePro(ctx, stg)
		exec = repo.ExecTx
	case polygontestrepo:
		repo := repository.NewPolygon(ctx, stg)
		exec = repo.ExecTx
	}
	require.NoError(t, err)

	return []interface{}{config.name, testsession(func(t *testing.T, name string, op testop) bool {
		err := exec(ctx, func(ctx context.Context, stg tools.GenericStorage) (bool, error) {
			ok = t.Run(name, func(t *testing.T) { op(ctx, t, stg, config.tr) })
			return false, nil
		})
		require.NoError(t, err)
		return ok
	})}
}

func testcases(t *testing.T) [][]interface{} {
	return [][]interface{}{
		testexecutor(t, testconfig{
			tr:       coinbaseprotestrepo,
			ts:       mongoteststg,
			name:     "mongo-coinbasepro",
			host:     "mongo-coinbasepro",
			port:     "27017",
			database: "coinbasepro",
		}),

		testexecutor(t, testconfig{
			tr:       coinbaseprotestrepo,
			ts:       postgresstg,
			name:     "postgres-coinbasepro",
			host:     "postgres-coinbasepro",
			username: "postgres",
			port:     "5432",
			database: "coinbasepro",
		}),

		testexecutor(t, testconfig{
			tr:       polygontestrepo,
			ts:       postgresstg,
			name:     "postgres-polygon",
			host:     "postgres-polygon",
			username: "postgres",
			port:     "5433",
			database: "polygon",
		}),
	}
}

func TestIntegration(t *testing.T) {
	// tolerance is the number of times to retry all tests
	tolerance := 3

	// upsertbenchmark are the number of records to attempt an upsert
	upsertbenchmark := func() int {
		rand.Seed(time.Now().UnixNano())
		return int(1e4) + rand.Intn(10)
	}

	for i := 1; i <= tolerance; i++ {
		for _, tc := range testcases(t) {
			name := tc[0].(string)
			t.Run(fmt.Sprintf("test %s %v", name, i), func(t *testing.T) {
				sessionfn := tc[1].(testsession)
				sessionfn(t, "UpsertCandleMinutes single",
					func(ctx context.Context, t *testing.T, stg tools.GenericStorage, tr testrepo) {
						if tr != coinbaseprotestrepo {
							return
						}
						repo := repository.NewCoinbasePro(ctx, stg)

						req := &proto.TruncateTablesRequest{Tables: []string{"candle_minutes"}}
						repo.TruncateTables(ctx, req)

						candles := coinbasepro.Candles{
							{PriceClose: 1.0, Unix: 1654315091, ProductID: "BTC-USD"},
						}
						rsp := new(proto.CreateResponse)
						c, err := repo.UpsertCandleMinutes(ctx, candles, rsp)
						require.NoErrorf(t, err, "error reading Coinbase Pro candles: %v", err)
						require.Len(t, c, 1)

						readrsp := new(proto.ReadResponse)
						err = repo.ReadCandleMinutes(ctx, "BTC-USD", nil, readrsp)
						require.NoError(t, err)

						rcandles := new(coinbasepro.Candles)
						err = tools.AssignReadResponseRecords(readrsp, rcandles)
						require.NoError(t, err)
						require.Equal(t, candles, *rcandles)
					})

				sessionfn(t, "UpsertCandleMinutes benchmark successive",
					func(ctx context.Context, t *testing.T, stg tools.GenericStorage, tr testrepo) {
						if tr != coinbaseprotestrepo {
							return
						}
						repo := repository.NewCoinbasePro(ctx, stg)

						// Clean up any extra data before testing.
						req := &proto.TruncateTablesRequest{Tables: []string{"candle_minutes"}}
						repo.TruncateTables(ctx, req)

						candles := coinbasepro.Candles{}
						ub := upsertbenchmark()
						for i := 1; i <= ub; i++ {
							candles = append(candles, &coinbasepro.Candle{
								Unix:       int64(i),
								ProductID:  "BTC-USD",
								PriceClose: 1.0,
								PriceHigh:  2.0,
								PriceLow:   3.0,
								PriceOpen:  4.0,
								Volume:     5.0,
							})
						}
						rsp := new(proto.CreateResponse)
						c, err := repo.UpsertCandleMinutes(ctx, candles, rsp)
						require.NoErrorf(t, err, "error reading Coinbase Pro candles: %v", err)
						require.Len(t, c, ub)

						readrsp := new(proto.ReadResponse)
						err = repo.ReadCandleMinutes(ctx, "BTC-USD", nil, readrsp)
						require.NoError(t, err)

						rcandles := new(coinbasepro.Candles)
						err = tools.AssignReadResponseRecords(readrsp, rcandles)
						require.NoError(t, err)
						require.Equal(t, candles, *rcandles)
					})

				sessionfn(t, "UpsertBarMinutes single",
					func(ctx context.Context, t *testing.T, stg tools.GenericStorage, tr testrepo) {
						if tr != polygontestrepo {
							return
						}
						repo := repository.NewPolygon(ctx, stg)

						// Clean up any extra data before testing.
						req := &proto.TruncateTablesRequest{Tables: []string{"bar_minutes"}}
						repo.TruncateTables(ctx, req)

						bars := []*polygon.BarResult{
							{Adjusted: true, C: 0.1, H: 0.2, L: 0.3, N: 0.4, O: 0.5,
								T: 1658375674, Ticker: "TEST", V: 0.7, Vw: 0.8},
						}

						rsp := new(proto.CreateResponse)
						c, err := repo.UpsertBarMinutes(ctx, bars, rsp)
						require.NoErrorf(t, err, "error reading Polygon Bars: %v", err)
						require.Len(t, c, 1)

						readrsp := new(proto.ReadResponse)

						start := time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)
						end := time.Now()

						err = repo.ReadBarMinutes(ctx, "TEST", true, start, end, readrsp)
						require.NoError(t, err)

						rbars := []*polygon.BarResult{}
						err = tools.AssignReadResponseRecords(readrsp, &rbars)
						require.NoError(t, err)
						require.Equal(t, bars, rbars)
					})

				sessionfn(t, "UpsertBarMinutes benchmark successive",
					func(ctx context.Context, t *testing.T, stg tools.GenericStorage, tr testrepo) {
						if tr != polygontestrepo {
							return
						}
						repo := repository.NewPolygon(ctx, stg)

						// Clean up any extra data before testing.
						req := &proto.TruncateTablesRequest{Tables: []string{"bar_minutes"}}
						repo.TruncateTables(ctx, req)

						bars := []*polygon.BarResult{}
						ub := upsertbenchmark()
						for i := 1; i <= ub; i++ {
							bars = append(bars, &polygon.BarResult{
								Adjusted: true,
								C:        0.1,
								H:        0.2,
								L:        0.3,
								N:        0.4,
								O:        0.5,
								T:        1658375674 + float64(i),
								Ticker:   "TEST",
								V:        0.7,
								Vw:       0.8,
							})
						}
						rsp := new(proto.CreateResponse)
						c, err := repo.UpsertBarMinutes(ctx, bars, rsp)
						require.NoErrorf(t, err, "error reading Polygon Bars: %v", err)
						require.Len(t, c, ub,
							"expected len to be %v, got %v", upsertbenchmark, len(c))

						readrsp := new(proto.ReadResponse)

						start := time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)
						end := time.Date(9999, 12, 31, 01, 01, 01, 01, time.UTC)

						err = repo.ReadBarMinutes(ctx, "TEST", true, start, end, readrsp)
						require.NoError(t, err)

						rbars := []*polygon.BarResult{}
						err = tools.AssignReadResponseRecords(readrsp, &rbars)
						require.NoError(t, err)
						require.Equal(t, bars, rbars)
					})
			})
		}
	}
}
