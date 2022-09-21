package coinbasepro

import (
	"encoding/json"
	"fmt"
)

// Candle represents the historic rate for a product at a point in time.
type Candle struct {
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

// Candles are the historic rates for a product. Rates are returned in grouped buckets. Candle schema is of the form
// `[timestamp, price_low, price_high, price_open, price_close]`
type Candles []*Candle

// UnmarshalJSON will deserialize bytes into a Candles model
func (candles *Candles) UnmarshalJSON(bytes []byte) error {
	var raw [][]float64
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return fmt.Errorf("error unmarshaling candles: %w", err)
	}
	for _, r := range raw {
		candle := new(Candle)
		candle.Unix = int64(r[0])
		candle.PriceLow = r[1]
		candle.PriceHigh = r[2]
		candle.PriceOpen = r[3]
		candle.PriceClose = r[4]
		candle.Volume = r[5]
		*candles = append(*candles, candle)
	}
	return nil
}
