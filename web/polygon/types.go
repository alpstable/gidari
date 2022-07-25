package polygon

// * This is a generated file, do not edit

type Status string

const (
	StatusOk         Status = "ok"
	StatusEarlyClose Status = "early-close"
	StatusClosed     Status = "closed"
)

// String will convert a Status into a string.
func (Status *Status) String() string {
	if Status != nil {
		return string(*Status)
	}
	return ""
}

// AssetClass is an identifier for a group of similar financial instruments.
type AssetClass string

const (
	AssetClassStocks  AssetClass = "stocks"
	AssetClassOptions AssetClass = "options"
	AssetClassCrypto  AssetClass = "crypto"
	AssetClassFx      AssetClass = "fx"
)

// String will convert a AssetClass into a string.
func (AssetClass *AssetClass) String() string {
	if AssetClass != nil {
		return string(*AssetClass)
	}
	return ""
}

// ConditionType is an identifier for a collection of related conditions.
type ConditionType string

const (
	ConditionTypeSaleCondition                 ConditionType = "sale_condition"
	ConditionTypeQuoteCondition                ConditionType = "quote_condition"
	ConditionTypeSipGeneratedFlag              ConditionType = "sip_generated_flag"
	ConditionTypeFinancialStatusIndicator      ConditionType = "financial_status_indicator"
	ConditionTypeShortSaleRestrictionIndicator ConditionType = "short_sale_restriction_indicator"
	ConditionTypeSettlementCondition           ConditionType = "settlement_condition"
	ConditionTypeMarketCondition               ConditionType = "market_condition"
	ConditionTypeTradeThruExempt               ConditionType = "trade_thru_exempt"
)

// String will convert a ConditionType into a string.
func (ConditionType *ConditionType) String() string {
	if ConditionType != nil {
		return string(*ConditionType)
	}
	return ""
}

// Timespan is the size of the time window.
type Timespan string

const (
	TimespanMinute  Timespan = "minute"
	TimespanHour    Timespan = "hour"
	TimespanDay     Timespan = "day"
	TimespanWeek    Timespan = "week"
	TimespanMonth   Timespan = "month"
	TimespanQuarter Timespan = "quarter"
	TimespanYear    Timespan = "year"
)

// String will convert a Timespan into a string.
func (Timespan *Timespan) String() string {
	if Timespan != nil {
		return string(*Timespan)
	}
	return ""
}
