package repository

import (
	"context"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/internal/query"
	"github.com/alpine-hodler/driver/tools"
	"github.com/alpine-hodler/driver/web/coinbasepro"
	"google.golang.org/protobuf/types/known/structpb"
)

// CoinbasePro is the repository wrapper for creating CRUD operations on a storage object.
type CoinbasePro interface {
	tools.GenericStorage

	ReadCandleMinutes(context.Context, string, *coinbasepro.CandlesOptions, *proto.ReadResponse) error
	UpsertCandleMinutes(context.Context, coinbasepro.Candles, *proto.CreateResponse) (coinbasepro.Candles, error)
}

type cbp struct{ *storage }

// NewCoinbaseProCandles will return a new service for performing CRUD operations on the Coinbase Pro candles web API.
func NewCoinbasePro(_ context.Context, r tools.GenericStorage) CoinbasePro {
	stg := new(storage)
	stg.r = newStorage(r)
	return &cbp{storage: stg}
}

// ReadCandleMinutes will query storage for coinbase pro candles to the 60-granularity for a given productID.
func (svc *cbp) ReadCandleMinutes(ctx context.Context, productID string, opts *coinbasepro.CandlesOptions,
	rsp *proto.ReadResponse) error {

	req := new(proto.ReadRequest)
	req.ReaderBuilder = query.CoinbaseProCandleMinutesReadBuilder.Bytes()
	req.Table = "candle_minutes"
	if err := tools.AssignReadOptions(req, opts); err != nil {
		return err
	}
	if err := tools.AssignReadRequired(req, "product_id", productID); err != nil {
		return err
	}
	return svc.r.Read(ctx, req, rsp)
}

// UpsertCandleMinutes will upsert candles to the 60-granularity.
func (svc *cbp) UpsertCandleMinutes(ctx context.Context, candles coinbasepro.Candles,
	rsp *proto.CreateResponse) (coinbasepro.Candles, error) {

	var records []*structpb.Struct
	if err := tools.MakeRecordsRequest(candles, &records); err != nil {
		return nil, err
	}
	req := new(proto.UpsertRequest)
	req.Table = "candle_minutes"
	req.Records = records
	return candles, svc.r.Upsert(ctx, req, rsp)
}
