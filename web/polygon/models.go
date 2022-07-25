package polygon

// * This is a generated file, do not edit

// Bar represents an aggregate bar for a stock over a given date range in custom time window sizes.
type Bar struct {
	// A request id assigned by the server.
	RequestID string `bson:"request_id" json:"request_id" sql:"request_id"`

	// An array of conditions that match your query.
	Results []*BarResult `bson:"results" json:"results" sql:"results"`

	// Status is the status of the market on the holiday.
	Status Status `bson:"status" json:"status" sql:"status"`

	// The exchange symbol that this item is traded under.
	Ticker string `bson:"ticker" json:"ticker" sql:"ticker"`

	// The number of aggregates (minute or day) used to generate the response.
	QueryCount int `bson:"query_count" json:"queryCount" sql:"query_count"`

	// The total number of results for this request.
	ResultsCount int `bson:"results_count" json:"resultsCount" sql:"results_count"`

	// Whether or not this response was adjusted for splits.
	Adjusted bool `bson:"adjusted" json:"adjusted" sql:"adjusted"`
}

// BarResult are conditions that match a bar query.
type BarResult struct {
	// Adjusted is whether or not the results are adjusted for splits.
	Adjusted bool `bson:"adjusted" json:"adjusted" sql:"adjusted"`

	// C is the close price for the symbol in the given time period.
	C float64 `bson:"c" json:"c" sql:"c"`

	// H is the highest price for the symbol in the given time period.
	H float64 `bson:"h" json:"h" sql:"h"`

	// L is the lowest price for the symbol in the given time period.
	L float64 `bson:"l" json:"l" sql:"l"`

	// N is the number of transactions in the aggregate window.
	N float64 `bson:"n" json:"n" sql:"n"`

	// O is the open price for the symbol in the given time period.
	O float64 `bson:"o" json:"o" sql:"o"`

	// T is the Unix Msec timestamp for the start of the aggregate window.
	T float64 `bson:"t" json:"t" sql:"t"`

	// Ticker is the ticker symbol of the stock/equity
	Ticker string `bson:"ticker" json:"ticker" sql:"ticker"`

	// V is the trading volume of the symbol in the given time period.
	V float64 `bson:"v" json:"v" sql:"v"`

	// VW is the volume weighted average price.
	Vw float64 `bson:"vw" json:"vw" sql:"vw"`
}

// Consolidated describes the aggregation rules on a consolidated (all exchanges) basis.
type Consolidated struct {
	// Whether or not trades with this condition update the high/low.
	UpdatesHighLow bool `bson:"updates_high_low" json:"updates_high_low" sql:"updates_high_low"`

	// Whether or not trades with this condition update the open/close.
	UpdatesOpenClose bool `bson:"updates_open_close" json:"updates_open_close" sql:"updates_open_close"`

	// Whether or not trades with this condition update the volume.
	UpdatesVolume bool `bson:"updates_volume" json:"updates_volume" sql:"updates_volume"`
}

// MarketCenter describes aggregation rules on a per-market-center basis.
type MarketCenter struct {
	// Whether or not trades with this condition update the high/low.
	UpdatesHighLow bool `bson:"updates_high_low" json:"updates_high_low" sql:"updates_high_low"`

	// Whether or not trades with this condition update the open/close.
	UpdatesOpenClose bool `bson:"updates_open_close" json:"updates_open_close" sql:"updates_open_close"`

	// Whether or not trades with this condition update the volume.
	UpdatesVolume bool `bson:"updates_volume" json:"updates_volume" sql:"updates_volume"`
}

// SkipMapping is a mapping to a symbol for each SIP that has this condition.
type SkipMapping struct {
	CTA  string `bson:"cta" json:"CTA" sql:"cta"`
	OPRA string `bson:"opra" json:"OPRA" sql:"opra"`
	UTP  string `bson:"utp" json:"UTP" sql:"utp"`
}

// StockTicker are conditions that match a stock ticker query.
type StockTicker struct {
	// Abbreviation is a commonly-used abbreviation for this condition.
	Abbreviation string `bson:"abbreviation" json:"abbreviation" sql:"abbreviation"`

	// AssetClass is an identifier for a group of similar financial instruments.
	AssetClass AssetClass `bson:"asset_class" json:"asset_class" sql:"asset_class"`

	// DataTypes are data types that this condition applies to.
	DataTypes []string `bson:"data_types" json:"data_types" sql:"data_types"`

	// Description is a short description of the semantics of this condition.
	Description string `bson:"description" json:"description" sql:"description"`

	// Exchange, if present, is a mapping this condition from a Polygon.io code to a SIP symbol depends on this attribute.
	// In other words, data with this condition attached comes exclusively from the given exchange.
	Exchange int `bson:"exchange" json:"exchange" sql:"exchange"`

	// ID is an identifier used by Polygon.io for this condition. Unique per data type.
	ID int `bson:"id" json:"id" sql:"id"`

	// Legacy, if true, is this condition is from an old version of the SIPs' specs and no longer is used. Other conditions
	// may or may not reuse the same symbol as this one.
	Legacy bool `bson:"legacy" json:"legacy" sql:"legacy"`

	// MarketCenter describes aggregation rules on a per-market-center basis.
	MarketCenter MarketCenter `bson:"market_center" json:"market_center" sql:"market_center"`

	// Name is the name of this condition.
	Name string `bson:"name" json:"name" sql:"name"`

	// SipMapping is a mapping to a symbol for each SIP that has this condition.
	SipMapping SkipMapping `bson:"sip_mapping" json:"sip_mapping" sql:"sip_mapping"`

	// Type is an identifier for a collection of related conditions.
	Type ConditionType `bson:"type" json:"type" sql:"type"`

	// UpdateRules are a list of aggregation rules.
	UpdateRules UpdateRules `bson:"update_rules" json:"update_rules" sql:"update_rules"`
}

// Upcoming returns information concerning market holidays and their open/close times.
type Upcoming struct {
	// Close is the market close time on the holiday (if it's not closed).
	Close string `bson:"close" json:"close" sql:"close"`

	// Date is the date of the holiday.
	Date string `bson:"date" json:"date" sql:"date"`

	// Exchange is market the record is for.
	Exchange string `bson:"exchange" json:"exchange" sql:"exchange"`

	// Name is the name of the holiday.
	Name string `bson:"name" json:"name" sql:"name"`

	// Open is the market open time on the holiday (if it's not closed).
	Open string `bson:"open" json:"open" sql:"open"`

	// Status is the status of the market on the holiday.
	Status string `bson:"status" json:"status" sql:"status"`
}

// UpdateRules are a list of aggregation rules.
type UpdateRules struct {
	// Describes aggregation rules on a consolidated (all exchanges) basis.
	Consolidated Consolidated `bson:"consolidated" json:"consolidated" sql:"consolidated"`
}
