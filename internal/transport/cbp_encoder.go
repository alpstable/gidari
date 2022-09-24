package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/alpine-hodler/gidari/proto"
	"github.com/alpine-hodler/gidari/tools"
)

const cbpCandlesTable = "candles"
const cbpCandlesMinutesGranularity = "60"

// cbpCandle represents the historic rate for a product at a point in time.
type cbpCandle struct {
	// PriceClose is the closing price (last trade) in the bucket interval.
	PriceClose float64 `bson:"price_close" json:"price_close" sql:"price_close"`

	// PriceHigh is the highest price during the bucket interval.
	PriceHigh float64 `bson:"price_high" json:"price_high" sql:"price_high"`

	// PriceLow is the lowest price during the bucket interval.
	PriceLow float64 `bson:"price_low" json:"price_low" sql:"price_low"`

	// PriceOpen is the opening price (first trade) in the bucket interval.
	PriceOpen float64 `bson:"price_open" json:"price_open" sql:"price_open"`

	// ProductID is the productID for the candle, e.g. BTC-ETH. This is not through the Coinbase Pro web API and is
	// intended for use in data layers and business logic.
	ProductID string `bson:"product_id" json:"product_id" sql:"product_id"`

	// Unix is the bucket start time as an int64 Unix value.
	Unix int64 `bson:"unix" json:"unix" sql:"unix"`

	// Volumes is the volume of trading activity during the bucket interval.
	Volume float64 `bson:"volume" json:"volume" sql:"volume"`
}

// andles are the historic rates for a product. Rates are returned in grouped buckets. Candle schema is of the form
// `[timestamp, price_low, price_high, price_open, price_close]`.
type cbpCandles []*cbpCandle

// UnmarshalJSON will deserialize bytes into a Candles model.
func (c *cbpCandles) UnmarshalJSON(bytes []byte) error {
	var rawNumbers [][]float64
	if err := json.Unmarshal(bytes, &rawNumbers); err != nil {
		return fmt.Errorf("error unmarshaling candles: %w", err)
	}

	for _, rnum := range rawNumbers {
		candle := new(cbpCandle)
		candle.Unix = int64(rnum[0])
		candle.PriceLow = rnum[1]
		candle.PriceHigh = rnum[2]
		candle.PriceOpen = rnum[3]
		candle.PriceClose = rnum[4]
		candle.Volume = rnum[5]
		*c = append(*c, candle)
	}

	return nil
}

// RegisterCBPEncoder will register the Coinbase Pro encoder.
func RegisterCBPEncoder() error {
	uri, err := url.Parse("https://api-public.sandbox.exchange.coinbase.com/candles")
	if err != nil {
		return fmt.Errorf("error parsing url: %w", err)
	}

	if err := RepositoryEncoders.Register(uri, new(CBPEncoder)); err != nil {
		return fmt.Errorf("error registering encoder: %w", err)
	}

	return nil
}

func cbpEncodeCandles(req http.Request, bytes []byte) ([]*proto.UpsertRequest, error) {
	// The default name for this encoder is "candles".
	table := cbpCandlesTable

	granularity := req.URL.Query()["granularity"][0]
	if granularity == cbpCandlesMinutesGranularity {
		table = "candle_minutes"
	}

	productID := tools.SplitURLPath(req)[1]

	// initialize the slice of candleSlice.
	var candleSlice cbpCandles

	if err := json.Unmarshal(bytes, &candleSlice); err != nil {
		return nil, fmt.Errorf("error unmarshaling candles: %w", err)
	}

	for _, candle := range candleSlice {
		candle.ProductID = productID
	}

	updatedBytes, err := json.Marshal(candleSlice)
	if err != nil {
		return nil, fmt.Errorf("error marshaling candles: %w", err)
	}

	upsertRequest := &proto.UpsertRequest{
		Table:    table,
		Data:     updatedBytes,
		DataType: int32(tools.UpsertDataJSON),
	}

	return []*proto.UpsertRequest{upsertRequest}, nil
}

// CBPEncoder is the encoder used to transform data from Coinbase Pro web requests into bytes that can be processed by
// repository upsert methods.
type CBPEncoder struct{}

// Encode will transform the data from Coinbase Pro Sandbox web requests into a byte slice that can be passed to
// repository.
func (ccre *CBPEncoder) Encode(req http.Request, bytes []byte) ([]*proto.UpsertRequest, error) {
	table, err := tools.ParseDBTableFromURL(req)
	if err != nil {
		return nil, fmt.Errorf("error getting table from request: %w", err)
	}

	switch table {
	case cbpCandlesTable:
		return cbpEncodeCandles(req, bytes)
	default:
		u, err := url.Parse("")
		if err != nil {
			return nil, fmt.Errorf("error parsing url: %w", err)
		}

		upsertRequest, err := RepositoryEncoders.Lookup(u).Encode(req, bytes)
		if err != nil {
			return nil, fmt.Errorf("error encoding data: %w", err)
		}

		return upsertRequest, nil
	}
}
