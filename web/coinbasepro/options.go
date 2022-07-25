package coinbasepro

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/alpine-hodler/web/internal"
)

// * This is a generated file, do not edit

// AccountHoldsOptions are options for API requests.
type AccountHoldsOptions struct {
	// After is used for pagination and sets end cursor to `after` date.
	After *string `bson:"after" json:"after" sql:"after"`

	// Before is used for pagination and sets start cursor to `before` date.
	Before *string `bson:"before" json:"before" sql:"before"`

	// Limit puts a limit on number of results to return.
	Limit *int `bson:"limit" json:"limit" sql:"limit"`
}

// AccountLedgerOptions are options for API requests.
type AccountLedgerOptions struct {
	// After is used for pagination. Sets end cursor to `after` date.
	After *int `bson:"after" json:"after" sql:"after"`

	// Before is used for pagination. Sets start cursor to `before` date.
	Before *int `bson:"before" json:"before" sql:"before"`

	// EndDate will filter results by maximum posted date.
	EndDate *string `bson:"end_date" json:"end_date" sql:"end_date"`

	// Limit puts a limit on number of results to return.
	Limit *int `bson:"limit" json:"limit" sql:"limit"`

	// StartDate will filter results by minimum posted date.
	StartDate *string `bson:"start_date" json:"start_date" sql:"start_date"`
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// BookOptions are options for API requests.
type BookOptions struct {
	// Levels 1 and 2 are aggregated. The size field is the sum of the size of the orders at that price, and num-orders is
	// the count of orders at that price; size should not be multiplied by num-orders. Level 3 is non-aggregated and returns
	// the entire order book. While the book is in an auction, the L1, L2 and L3 book will also contain the most recent
	// indicative quote disseminated during the auction, and auction_mode will be set to true. These indicative quote
	// messages are sent on an interval basis (approximately once a second) during the collection phase of an auction and
	// includes information about the tentative price and size affiliated with the completion. Level 1 and Level 2 are
	// recommended for polling. For the most up-to-date data, consider using the websocket stream. Level 3 is only
	// recommended for users wishing to maintain a full real-time order book using the websocket stream. Abuse of Level 3
	// via polling will cause your access to be limited or blocked.
	Level *int32 `bson:"level" json:"level" sql:"level"`
}

// CandlesOptions are options for API requests.
type CandlesOptions struct {
	// End is a timestamp for ending range of aggregations.
	End *string `bson:"end" json:"end" sql:"end"`

	// Granularity is one of the following values: {60, 300, 900, 3600, 21600, 86400}. Otherwise, your request will be
	// rejected. These values correspond to timeslices representing one minute, five minutes, fifteen minutes, one hour, six
	// hours, and one day, respectively.
	Granularity *Granularity `bson:"granularity" json:"granularity" sql:"granularity"`

	// Start is a timestamp for starting range of aggregations.
	Start *string `bson:"start" json:"start" sql:"start"`
}

// CreateOrderOptions are options for API requests.
type CreateOrderOptions struct {
	CancelAfter *CancelAfter `bson:"cancel_after" json:"cancel_after" sql:"cancel_after"`
	ClientOid   *string      `bson:"client_oid" json:"client_oid" sql:"client_oid"`
	Funds       *float64     `bson:"funds" json:"funds" sql:"funds"`
	PostOnly    *bool        `bson:"post_only" json:"post_only" sql:"post_only"`
	Price       *float64     `bson:"price" json:"price" sql:"price"`
	ProductID   string       `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID   *string      `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	STP         *STP         `bson:"stp" json:"stp" sql:"stp"`
	Side        Side         `bson:"side" json:"side" sql:"side"`
	Size        *float64     `bson:"size" json:"size" sql:"size"`
	Stop        *Stop        `bson:"stop" json:"stop" sql:"stop"`
	StopPrice   *float64     `bson:"stop_price" json:"stop_price" sql:"stop_price"`
	TimeInForce *TimeInForce `bson:"time_in_force" json:"time_in_force" sql:"time_in_force"`
	Type        *OrderType   `bson:"type" json:"type" sql:"type"`
}

// CreateReportOptions are options for API requests.
type CreateReportOptions struct {
	// Account - required for account-type reports
	AccountID *string `bson:"account_id" json:"account_id" sql:"account_id"`

	// Email to send generated report to
	Email *string `bson:"email" json:"email" sql:"email"`

	// End date for items to be included in report
	EndDate *string `bson:"end_date" json:"end_date" sql:"end_date"`

	// Portfolio - Which portfolio to generate the report for
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`

	// Product - required for fills-type reports
	ProductID *string `bson:"product_id" json:"product_id" sql:"product_id"`

	// Start date for items to be included in report.
	StartDate *string `bson:"start_date" json:"start_date" sql:"start_date"`

	// required for 1099k-transaction-history-type reports
	Year   *string     `bson:"year" json:"year" sql:"year"`
	Format *FileFormat `bson:"format" json:"format" sql:"format"`
	Type   ReportType  `bson:"type" json:"type" sql:"type"`
}

// ConvertCurrencyOptions are options for API requests.
type ConvertCurrencyOptions struct {
	Amount    float64 `bson:"amount" json:"amount" sql:"amount"`
	From      string  `bson:"from" json:"from" sql:"from"`
	Nonce     *string `bson:"nonce" json:"nonce" sql:"nonce"`
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	To        string  `bson:"to" json:"to" sql:"to"`
}

// CurrencyConversionOptions are options for API requests.
type CurrencyConversionOptions struct {
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// CoinbaseAccountDepositOptions are options for API requests.
type CoinbaseAccountDepositOptions struct {
	Amount            float64 `bson:"amount" json:"amount" sql:"amount"`
	CoinbaseAccountID string  `bson:"coinbase_account_id" json:"coinbase_account_id" sql:"coinbase_account_id"`
	Currency          string  `bson:"currency" json:"currency" sql:"currency"`
	ProfileID         *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// PaymentMethodDepositOptions are options for API requests.
type PaymentMethodDepositOptions struct {
	Amount          float64 `bson:"amount" json:"amount" sql:"amount"`
	Currency        string  `bson:"currency" json:"currency" sql:"currency"`
	PaymentMethodID string  `bson:"payment_method_id" json:"payment_method_id" sql:"payment_method_id"`
	ProfileID       *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// FillsOptions are options for API requests.
type FillsOptions struct {
	After     *int    `bson:"after" json:"after" sql:"after"`
	Before    *int    `bson:"before" json:"before" sql:"before"`
	Limit     *int    `bson:"limit" json:"limit" sql:"limit"`
	OrderID   *string `bson:"order_id" json:"order_id" sql:"order_id"`
	ProductID *string `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// OrdersOptions are options for API requests.
type OrdersOptions struct {
	After     *string  `bson:"after" json:"after" sql:"after"`
	Before    *string  `bson:"before" json:"before" sql:"before"`
	EndDate   *string  `bson:"end_date" json:"end_date" sql:"end_date"`
	Limit     int      `bson:"limit" json:"limit" sql:"limit"`
	ProductID *string  `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID *string  `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	SortedBy  *string  `bson:"sorted_by" json:"sortedBy" sql:"sorted_by"`
	Sorting   *string  `bson:"sorting" json:"sorting" sql:"sorting"`
	StartDate *string  `bson:"start_date" json:"start_date" sql:"start_date"`
	Status    []string `bson:"status" json:"status" sql:"status"`
}

// CancelOpenOrdersOptions are options for API requests.
type CancelOpenOrdersOptions struct {
	ProductID *string `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// ProductsOptions are options for API requests.
type ProductsOptions struct {
	Type *string `bson:"type" json:"type" sql:"type"`
}

// DeleteProfileOptions are options for API requests.
type DeleteProfileOptions struct {
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	To        *string `bson:"to" json:"to" sql:"to"`
}

// ProfileOptions are options for API requests.
type ProfileOptions struct {
	Active *bool `bson:"active" json:"active" sql:"active"`
}

// RenameProfileOptions are options for API requests.
type RenameProfileOptions struct {
	Name      *string `bson:"name" json:"name" sql:"name"`
	ProfileID *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// ProfilesOptions are options for API requests.
type ProfilesOptions struct {
	Active *bool `bson:"active" json:"active" sql:"active"`
}

// CreateProfileOptions are options for API requests.
type CreateProfileOptions struct {
	Name *string `bson:"name" json:"name" sql:"name"`
}

// CreateProfileTransferOptions are options for API requests.
type CreateProfileTransferOptions struct {
	Amount   *string `bson:"amount" json:"amount" sql:"amount"`
	Currency *string `bson:"currency" json:"currency" sql:"currency"`
	From     *string `bson:"from" json:"from" sql:"from"`
	To       *string `bson:"to" json:"to" sql:"to"`
}

// ReportsOptions are options for API requests.
type ReportsOptions struct {
	// Filter results after a specific date
	After *string `bson:"after" json:"after" sql:"after"`

	// Filter results by a specific profile_id
	PortfolioID *string `bson:"portfolio_id" json:"portfolio_id" sql:"portfolio_id"`

	// Filter results by type of report (fills or account) - otc_fills: real string is otc-fills -
	// type_1099k_transaction_history: real string is 1099-transaction-history - tax_invoice: real string is tax-invoice
	Type *ReportType `bson:"type" json:"type" sql:"type"`

	// Ignore expired results
	IgnoredExpired *bool `bson:"ignored_expired" json:"ignored_expired" sql:"ignored_expired"`

	// Limit results to a specific number
	Limit *int `bson:"limit" json:"limit" sql:"limit"`
}

// TradesOptions are options for API requests.
type TradesOptions struct {
	After  *int32 `bson:"after" json:"after" sql:"after"`
	Before *int32 `bson:"before" json:"before" sql:"before"`
	Limit  *int32 `bson:"limit" json:"limit" sql:"limit"`
}

// AccountTransfersOptions are options for API requests.
type AccountTransfersOptions struct {
	// After is used for pagination. Sets end cursor to `after` date.
	After *string `bson:"after" json:"after" sql:"after"`

	// Before is used for pagination. Sets start cursor to `before` date.
	Before *string `bson:"before" json:"before" sql:"before"`

	// Limit puts a limit on number of results to return.
	Limit *int            `bson:"limit" json:"limit" sql:"limit"`
	Type  *TransferMethod `bson:"type" json:"type" sql:"type"`
}

// CoinbaseAccountWithdrawalOptions are options for API requests.
type CoinbaseAccountWithdrawalOptions struct {
	Amount            float64 `bson:"amount" json:"amount" sql:"amount"`
	CoinbaseAccountID string  `bson:"coinbase_account_id" json:"coinbase_account_id" sql:"coinbase_account_id"`
	Currency          string  `bson:"currency" json:"currency" sql:"currency"`
	ProfileID         *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// CryptoWithdrawalOptions are options for API requests.
type CryptoWithdrawalOptions struct {
	Amount           float64  `bson:"amount" json:"amount" sql:"amount"`
	CryptoAddress    string   `bson:"crypto_address" json:"crypto_address" sql:"crypto_address"`
	Currency         string   `bson:"currency" json:"currency" sql:"currency"`
	DestinationTag   *string  `bson:"destination_tag" json:"destination_tag" sql:"destination_tag"`
	Fee              *float64 `bson:"fee" json:"fee" sql:"fee"`
	NoDestinationTag *bool    `bson:"no_destination_tag" json:"no_destination_tag" sql:"no_destination_tag"`
	Nonce            *int     `bson:"nonce" json:"nonce" sql:"nonce"`
	ProfileID        *string  `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	TwoFactorCode    *string  `bson:"two_factor_code" json:"two_factor_code" sql:"two_factor_code"`
}

// PaymentMethodWithdrawalOptions are options for API requests.
type PaymentMethodWithdrawalOptions struct {
	Amount          float64 `bson:"amount" json:"amount" sql:"amount"`
	Currency        string  `bson:"currency" json:"currency" sql:"currency"`
	PaymentMethodID string  `bson:"payment_method_id" json:"payment_method_id" sql:"payment_method_id"`
	ProfileID       *string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
}

// WithdrawalFeeEstimateOptions are options for API requests.
type WithdrawalFeeEstimateOptions struct {
	CryptoAddress *string `bson:"crypto_address" json:"crypto_address" sql:"crypto_address"`
	Currency      *string `bson:"currency" json:"currency" sql:"currency"`
}

func (opts *AccountHoldsOptions) EncodeBody() (buf io.Reader, err error)          { return }
func (opts *AccountLedgerOptions) EncodeBody() (buf io.Reader, err error)         { return }
func (opts *AccountTransfersOptions) EncodeBody() (buf io.Reader, err error)      { return }
func (opts *BookOptions) EncodeBody() (buf io.Reader, err error)                  { return }
func (opts *CancelOpenOrdersOptions) EncodeBody() (buf io.Reader, err error)      { return }
func (opts *CandlesOptions) EncodeBody() (buf io.Reader, err error)               { return }
func (opts *CurrencyConversionOptions) EncodeBody() (buf io.Reader, err error)    { return }
func (opts *DeleteProfileOptions) EncodeBody() (buf io.Reader, err error)         { return }
func (opts *FillsOptions) EncodeBody() (buf io.Reader, err error)                 { return }
func (opts *OrdersOptions) EncodeBody() (buf io.Reader, err error)                { return }
func (opts *ProductsOptions) EncodeBody() (buf io.Reader, err error)              { return }
func (opts *ProfileOptions) EncodeBody() (buf io.Reader, err error)               { return }
func (opts *ProfilesOptions) EncodeBody() (buf io.Reader, err error)              { return }
func (opts *RenameProfileOptions) EncodeBody() (buf io.Reader, err error)         { return }
func (opts *ReportsOptions) EncodeBody() (buf io.Reader, err error)               { return }
func (opts *TradesOptions) EncodeBody() (buf io.Reader, err error)                { return }
func (opts *WithdrawalFeeEstimateOptions) EncodeBody() (buf io.Reader, err error) { return }
func (opts *CoinbaseAccountDepositOptions) EncodeQuery(req *http.Request)         { return }
func (opts *CoinbaseAccountWithdrawalOptions) EncodeQuery(req *http.Request)      { return }
func (opts *ConvertCurrencyOptions) EncodeQuery(req *http.Request)                { return }
func (opts *CreateOrderOptions) EncodeQuery(req *http.Request)                    { return }
func (opts *CreateProfileOptions) EncodeQuery(req *http.Request)                  { return }
func (opts *CreateProfileTransferOptions) EncodeQuery(req *http.Request)          { return }
func (opts *CreateReportOptions) EncodeQuery(req *http.Request)                   { return }
func (opts *CryptoWithdrawalOptions) EncodeQuery(req *http.Request)               { return }
func (opts *PaymentMethodDepositOptions) EncodeQuery(req *http.Request)           { return }
func (opts *PaymentMethodWithdrawalOptions) EncodeQuery(req *http.Request)        { return }

// SetBefore sets the Before field on AccountHoldsOptions. Before is used for pagination and sets start cursor to
// `before` date.
func (opts *AccountHoldsOptions) SetBefore(Before string) *AccountHoldsOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on AccountHoldsOptions. After is used for pagination and sets end cursor to `after`
// date.
func (opts *AccountHoldsOptions) SetAfter(After string) *AccountHoldsOptions {
	opts.After = &After
	return opts
}

// SetLimit sets the Limit field on AccountHoldsOptions. Limit puts a limit on number of results to return.
func (opts *AccountHoldsOptions) SetLimit(Limit int) *AccountHoldsOptions {
	opts.Limit = &Limit
	return opts
}

// SetStartDate sets the StartDate field on AccountLedgerOptions. StartDate will filter results by minimum posted date.
func (opts *AccountLedgerOptions) SetStartDate(StartDate string) *AccountLedgerOptions {
	opts.StartDate = &StartDate
	return opts
}

// SetEndDate sets the EndDate field on AccountLedgerOptions. EndDate will filter results by maximum posted date.
func (opts *AccountLedgerOptions) SetEndDate(EndDate string) *AccountLedgerOptions {
	opts.EndDate = &EndDate
	return opts
}

// SetBefore sets the Before field on AccountLedgerOptions. Before is used for pagination. Sets start cursor to `before`
// date.
func (opts *AccountLedgerOptions) SetBefore(Before int) *AccountLedgerOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on AccountLedgerOptions. After is used for pagination. Sets end cursor to `after` date.
func (opts *AccountLedgerOptions) SetAfter(After int) *AccountLedgerOptions {
	opts.After = &After
	return opts
}

// SetProfileID sets the ProfileID field on AccountLedgerOptions.
func (opts *AccountLedgerOptions) SetProfileID(ProfileID string) *AccountLedgerOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetLimit sets the Limit field on AccountLedgerOptions. Limit puts a limit on number of results to return.
func (opts *AccountLedgerOptions) SetLimit(Limit int) *AccountLedgerOptions {
	opts.Limit = &Limit
	return opts
}

// SetBefore sets the Before field on AccountTransfersOptions. Before is used for pagination. Sets start cursor to
// `before` date.
func (opts *AccountTransfersOptions) SetBefore(Before string) *AccountTransfersOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on AccountTransfersOptions. After is used for pagination. Sets end cursor to `after`
// date.
func (opts *AccountTransfersOptions) SetAfter(After string) *AccountTransfersOptions {
	opts.After = &After
	return opts
}

// SetLimit sets the Limit field on AccountTransfersOptions. Limit puts a limit on number of results to return.
func (opts *AccountTransfersOptions) SetLimit(Limit int) *AccountTransfersOptions {
	opts.Limit = &Limit
	return opts
}

// SetType sets the Type field on AccountTransfersOptions.
func (opts *AccountTransfersOptions) SetType(Type TransferMethod) *AccountTransfersOptions {
	opts.Type = &Type
	return opts
}

// SetLevel sets the Level field on BookOptions. Levels 1 and 2 are aggregated. The size field is the sum of the size of
// the orders at that price, and num-orders is the count of orders at that price; size should not be multiplied by
// num-orders. Level 3 is non-aggregated and returns the entire order book. While the book is in an auction, the L1, L2
// and L3 book will also contain the most recent indicative quote disseminated during the auction, and auction_mode will
// be set to true. These indicative quote messages are sent on an interval basis (approximately once a second) during
// the collection phase of an auction and includes information about the tentative price and size affiliated with the
// completion. Level 1 and Level 2 are recommended for polling. For the most up-to-date data, consider using the
// websocket stream. Level 3 is only recommended for users wishing to maintain a full real-time order book using the
// websocket stream. Abuse of Level 3 via polling will cause your access to be limited or blocked.
func (opts *BookOptions) SetLevel(Level int32) *BookOptions {
	opts.Level = &Level
	return opts
}

// SetProfileID sets the ProfileID field on CancelOpenOrdersOptions.
func (opts *CancelOpenOrdersOptions) SetProfileID(ProfileID string) *CancelOpenOrdersOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetProductID sets the ProductID field on CancelOpenOrdersOptions.
func (opts *CancelOpenOrdersOptions) SetProductID(ProductID string) *CancelOpenOrdersOptions {
	opts.ProductID = &ProductID
	return opts
}

// SetGranularity sets the Granularity field on CandlesOptions. Granularity is one of the following values: {60, 300,
// 900, 3600, 21600, 86400}. Otherwise, your request will be rejected. These values correspond to timeslices
// representing one minute, five minutes, fifteen minutes, one hour, six hours, and one day, respectively.
func (opts *CandlesOptions) SetGranularity(Granularity Granularity) *CandlesOptions {
	opts.Granularity = &Granularity
	return opts
}

// SetStart sets the Start field on CandlesOptions. Start is a timestamp for starting range of aggregations.
func (opts *CandlesOptions) SetStart(Start string) *CandlesOptions {
	opts.Start = &Start
	return opts
}

// SetEnd sets the End field on CandlesOptions. End is a timestamp for ending range of aggregations.
func (opts *CandlesOptions) SetEnd(End string) *CandlesOptions {
	opts.End = &End
	return opts
}

// SetProfileID sets the ProfileID field on CoinbaseAccountDepositOptions.
func (opts *CoinbaseAccountDepositOptions) SetProfileID(ProfileID string) *CoinbaseAccountDepositOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetAmount sets the Amount field on CoinbaseAccountDepositOptions.
func (opts *CoinbaseAccountDepositOptions) SetAmount(Amount float64) *CoinbaseAccountDepositOptions {
	opts.Amount = Amount
	return opts
}

// SetCoinbaseAccountID sets the CoinbaseAccountID field on CoinbaseAccountDepositOptions.
func (opts *CoinbaseAccountDepositOptions) SetCoinbaseAccountID(CoinbaseAccountID string) *CoinbaseAccountDepositOptions {
	opts.CoinbaseAccountID = CoinbaseAccountID
	return opts
}

// SetCurrency sets the Currency field on CoinbaseAccountDepositOptions.
func (opts *CoinbaseAccountDepositOptions) SetCurrency(Currency string) *CoinbaseAccountDepositOptions {
	opts.Currency = Currency
	return opts
}

// SetProfileID sets the ProfileID field on CoinbaseAccountWithdrawalOptions.
func (opts *CoinbaseAccountWithdrawalOptions) SetProfileID(ProfileID string) *CoinbaseAccountWithdrawalOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetAmount sets the Amount field on CoinbaseAccountWithdrawalOptions.
func (opts *CoinbaseAccountWithdrawalOptions) SetAmount(Amount float64) *CoinbaseAccountWithdrawalOptions {
	opts.Amount = Amount
	return opts
}

// SetCoinbaseAccountID sets the CoinbaseAccountID field on CoinbaseAccountWithdrawalOptions.
func (opts *CoinbaseAccountWithdrawalOptions) SetCoinbaseAccountID(CoinbaseAccountID string) *CoinbaseAccountWithdrawalOptions {
	opts.CoinbaseAccountID = CoinbaseAccountID
	return opts
}

// SetCurrency sets the Currency field on CoinbaseAccountWithdrawalOptions.
func (opts *CoinbaseAccountWithdrawalOptions) SetCurrency(Currency string) *CoinbaseAccountWithdrawalOptions {
	opts.Currency = Currency
	return opts
}

// SetProfileID sets the ProfileID field on ConvertCurrencyOptions.
func (opts *ConvertCurrencyOptions) SetProfileID(ProfileID string) *ConvertCurrencyOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetFrom sets the From field on ConvertCurrencyOptions.
func (opts *ConvertCurrencyOptions) SetFrom(From string) *ConvertCurrencyOptions {
	opts.From = From
	return opts
}

// SetTo sets the To field on ConvertCurrencyOptions.
func (opts *ConvertCurrencyOptions) SetTo(To string) *ConvertCurrencyOptions {
	opts.To = To
	return opts
}

// SetAmount sets the Amount field on ConvertCurrencyOptions.
func (opts *ConvertCurrencyOptions) SetAmount(Amount float64) *ConvertCurrencyOptions {
	opts.Amount = Amount
	return opts
}

// SetNonce sets the Nonce field on ConvertCurrencyOptions.
func (opts *ConvertCurrencyOptions) SetNonce(Nonce string) *ConvertCurrencyOptions {
	opts.Nonce = &Nonce
	return opts
}

// SetProfileID sets the ProfileID field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetProfileID(ProfileID string) *CreateOrderOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetType sets the Type field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetType(Type OrderType) *CreateOrderOptions {
	opts.Type = &Type
	return opts
}

// SetSide sets the Side field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetSide(Side Side) *CreateOrderOptions {
	opts.Side = Side
	return opts
}

// SetSTP sets the STP field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetSTP(STP STP) *CreateOrderOptions {
	opts.STP = &STP
	return opts
}

// SetStop sets the Stop field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetStop(Stop Stop) *CreateOrderOptions {
	opts.Stop = &Stop
	return opts
}

// SetStopPrice sets the StopPrice field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetStopPrice(StopPrice float64) *CreateOrderOptions {
	opts.StopPrice = &StopPrice
	return opts
}

// SetPrice sets the Price field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetPrice(Price float64) *CreateOrderOptions {
	opts.Price = &Price
	return opts
}

// SetSize sets the Size field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetSize(Size float64) *CreateOrderOptions {
	opts.Size = &Size
	return opts
}

// SetFunds sets the Funds field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetFunds(Funds float64) *CreateOrderOptions {
	opts.Funds = &Funds
	return opts
}

// SetProductID sets the ProductID field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetProductID(ProductID string) *CreateOrderOptions {
	opts.ProductID = ProductID
	return opts
}

// SetTimeInForce sets the TimeInForce field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetTimeInForce(TimeInForce TimeInForce) *CreateOrderOptions {
	opts.TimeInForce = &TimeInForce
	return opts
}

// SetCancelAfter sets the CancelAfter field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetCancelAfter(CancelAfter CancelAfter) *CreateOrderOptions {
	opts.CancelAfter = &CancelAfter
	return opts
}

// SetPostOnly sets the PostOnly field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetPostOnly(PostOnly bool) *CreateOrderOptions {
	opts.PostOnly = &PostOnly
	return opts
}

// SetClientOid sets the ClientOid field on CreateOrderOptions.
func (opts *CreateOrderOptions) SetClientOid(ClientOid string) *CreateOrderOptions {
	opts.ClientOid = &ClientOid
	return opts
}

// SetName sets the Name field on CreateProfileOptions.
func (opts *CreateProfileOptions) SetName(Name string) *CreateProfileOptions {
	opts.Name = &Name
	return opts
}

// SetFrom sets the From field on CreateProfileTransferOptions.
func (opts *CreateProfileTransferOptions) SetFrom(From string) *CreateProfileTransferOptions {
	opts.From = &From
	return opts
}

// SetTo sets the To field on CreateProfileTransferOptions.
func (opts *CreateProfileTransferOptions) SetTo(To string) *CreateProfileTransferOptions {
	opts.To = &To
	return opts
}

// SetCurrency sets the Currency field on CreateProfileTransferOptions.
func (opts *CreateProfileTransferOptions) SetCurrency(Currency string) *CreateProfileTransferOptions {
	opts.Currency = &Currency
	return opts
}

// SetAmount sets the Amount field on CreateProfileTransferOptions.
func (opts *CreateProfileTransferOptions) SetAmount(Amount string) *CreateProfileTransferOptions {
	opts.Amount = &Amount
	return opts
}

// SetStartDate sets the StartDate field on CreateReportOptions. Start date for items to be included in report.
func (opts *CreateReportOptions) SetStartDate(StartDate string) *CreateReportOptions {
	opts.StartDate = &StartDate
	return opts
}

// SetEndDate sets the EndDate field on CreateReportOptions. End date for items to be included in report
func (opts *CreateReportOptions) SetEndDate(EndDate string) *CreateReportOptions {
	opts.EndDate = &EndDate
	return opts
}

// SetType sets the Type field on CreateReportOptions.
func (opts *CreateReportOptions) SetType(Type ReportType) *CreateReportOptions {
	opts.Type = Type
	return opts
}

// SetYear sets the Year field on CreateReportOptions. required for 1099k-transaction-history-type reports
func (opts *CreateReportOptions) SetYear(Year string) *CreateReportOptions {
	opts.Year = &Year
	return opts
}

// SetFormat sets the Format field on CreateReportOptions.
func (opts *CreateReportOptions) SetFormat(Format FileFormat) *CreateReportOptions {
	opts.Format = &Format
	return opts
}

// SetProductID sets the ProductID field on CreateReportOptions. Product - required for fills-type reports
func (opts *CreateReportOptions) SetProductID(ProductID string) *CreateReportOptions {
	opts.ProductID = &ProductID
	return opts
}

// SetAccountID sets the AccountID field on CreateReportOptions. Account - required for account-type reports
func (opts *CreateReportOptions) SetAccountID(AccountID string) *CreateReportOptions {
	opts.AccountID = &AccountID
	return opts
}

// SetEmail sets the Email field on CreateReportOptions. Email to send generated report to
func (opts *CreateReportOptions) SetEmail(Email string) *CreateReportOptions {
	opts.Email = &Email
	return opts
}

// SetProfileID sets the ProfileID field on CreateReportOptions. Portfolio - Which portfolio to generate the report for
func (opts *CreateReportOptions) SetProfileID(ProfileID string) *CreateReportOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetProfileID sets the ProfileID field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetProfileID(ProfileID string) *CryptoWithdrawalOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetAmount sets the Amount field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetAmount(Amount float64) *CryptoWithdrawalOptions {
	opts.Amount = Amount
	return opts
}

// SetCryptoAddress sets the CryptoAddress field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetCryptoAddress(CryptoAddress string) *CryptoWithdrawalOptions {
	opts.CryptoAddress = CryptoAddress
	return opts
}

// SetCurrency sets the Currency field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetCurrency(Currency string) *CryptoWithdrawalOptions {
	opts.Currency = Currency
	return opts
}

// SetDestinationTag sets the DestinationTag field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetDestinationTag(DestinationTag string) *CryptoWithdrawalOptions {
	opts.DestinationTag = &DestinationTag
	return opts
}

// SetNoDestinationTag sets the NoDestinationTag field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetNoDestinationTag(NoDestinationTag bool) *CryptoWithdrawalOptions {
	opts.NoDestinationTag = &NoDestinationTag
	return opts
}

// SetTwoFactorCode sets the TwoFactorCode field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetTwoFactorCode(TwoFactorCode string) *CryptoWithdrawalOptions {
	opts.TwoFactorCode = &TwoFactorCode
	return opts
}

// SetNonce sets the Nonce field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetNonce(Nonce int) *CryptoWithdrawalOptions {
	opts.Nonce = &Nonce
	return opts
}

// SetFee sets the Fee field on CryptoWithdrawalOptions.
func (opts *CryptoWithdrawalOptions) SetFee(Fee float64) *CryptoWithdrawalOptions {
	opts.Fee = &Fee
	return opts
}

// SetProfileID sets the ProfileID field on CurrencyConversionOptions.
func (opts *CurrencyConversionOptions) SetProfileID(ProfileID string) *CurrencyConversionOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetProfileID sets the ProfileID field on DeleteProfileOptions.
func (opts *DeleteProfileOptions) SetProfileID(ProfileID string) *DeleteProfileOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetTo sets the To field on DeleteProfileOptions.
func (opts *DeleteProfileOptions) SetTo(To string) *DeleteProfileOptions {
	opts.To = &To
	return opts
}

// SetOrderID sets the OrderID field on FillsOptions.
func (opts *FillsOptions) SetOrderID(OrderID string) *FillsOptions {
	opts.OrderID = &OrderID
	return opts
}

// SetProductID sets the ProductID field on FillsOptions.
func (opts *FillsOptions) SetProductID(ProductID string) *FillsOptions {
	opts.ProductID = &ProductID
	return opts
}

// SetProfileID sets the ProfileID field on FillsOptions.
func (opts *FillsOptions) SetProfileID(ProfileID string) *FillsOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetLimit sets the Limit field on FillsOptions.
func (opts *FillsOptions) SetLimit(Limit int) *FillsOptions {
	opts.Limit = &Limit
	return opts
}

// SetBefore sets the Before field on FillsOptions.
func (opts *FillsOptions) SetBefore(Before int) *FillsOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on FillsOptions.
func (opts *FillsOptions) SetAfter(After int) *FillsOptions {
	opts.After = &After
	return opts
}

// SetProfileID sets the ProfileID field on OrdersOptions.
func (opts *OrdersOptions) SetProfileID(ProfileID string) *OrdersOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetProductID sets the ProductID field on OrdersOptions.
func (opts *OrdersOptions) SetProductID(ProductID string) *OrdersOptions {
	opts.ProductID = &ProductID
	return opts
}

// SetSortedBy sets the SortedBy field on OrdersOptions.
func (opts *OrdersOptions) SetSortedBy(SortedBy string) *OrdersOptions {
	opts.SortedBy = &SortedBy
	return opts
}

// SetSorting sets the Sorting field on OrdersOptions.
func (opts *OrdersOptions) SetSorting(Sorting string) *OrdersOptions {
	opts.Sorting = &Sorting
	return opts
}

// SetStartDate sets the StartDate field on OrdersOptions.
func (opts *OrdersOptions) SetStartDate(StartDate string) *OrdersOptions {
	opts.StartDate = &StartDate
	return opts
}

// SetEndDate sets the EndDate field on OrdersOptions.
func (opts *OrdersOptions) SetEndDate(EndDate string) *OrdersOptions {
	opts.EndDate = &EndDate
	return opts
}

// SetBefore sets the Before field on OrdersOptions.
func (opts *OrdersOptions) SetBefore(Before string) *OrdersOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on OrdersOptions.
func (opts *OrdersOptions) SetAfter(After string) *OrdersOptions {
	opts.After = &After
	return opts
}

// SetLimit sets the Limit field on OrdersOptions.
func (opts *OrdersOptions) SetLimit(Limit int) *OrdersOptions {
	opts.Limit = Limit
	return opts
}

// SetStatus sets the Status field on OrdersOptions.
func (opts *OrdersOptions) SetStatus(Status []string) *OrdersOptions {
	opts.Status = Status
	return opts
}

// SetProfileID sets the ProfileID field on PaymentMethodDepositOptions.
func (opts *PaymentMethodDepositOptions) SetProfileID(ProfileID string) *PaymentMethodDepositOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetAmount sets the Amount field on PaymentMethodDepositOptions.
func (opts *PaymentMethodDepositOptions) SetAmount(Amount float64) *PaymentMethodDepositOptions {
	opts.Amount = Amount
	return opts
}

// SetPaymentMethodID sets the PaymentMethodID field on PaymentMethodDepositOptions.
func (opts *PaymentMethodDepositOptions) SetPaymentMethodID(PaymentMethodID string) *PaymentMethodDepositOptions {
	opts.PaymentMethodID = PaymentMethodID
	return opts
}

// SetCurrency sets the Currency field on PaymentMethodDepositOptions.
func (opts *PaymentMethodDepositOptions) SetCurrency(Currency string) *PaymentMethodDepositOptions {
	opts.Currency = Currency
	return opts
}

// SetProfileID sets the ProfileID field on PaymentMethodWithdrawalOptions.
func (opts *PaymentMethodWithdrawalOptions) SetProfileID(ProfileID string) *PaymentMethodWithdrawalOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetAmount sets the Amount field on PaymentMethodWithdrawalOptions.
func (opts *PaymentMethodWithdrawalOptions) SetAmount(Amount float64) *PaymentMethodWithdrawalOptions {
	opts.Amount = Amount
	return opts
}

// SetPaymentMethodID sets the PaymentMethodID field on PaymentMethodWithdrawalOptions.
func (opts *PaymentMethodWithdrawalOptions) SetPaymentMethodID(PaymentMethodID string) *PaymentMethodWithdrawalOptions {
	opts.PaymentMethodID = PaymentMethodID
	return opts
}

// SetCurrency sets the Currency field on PaymentMethodWithdrawalOptions.
func (opts *PaymentMethodWithdrawalOptions) SetCurrency(Currency string) *PaymentMethodWithdrawalOptions {
	opts.Currency = Currency
	return opts
}

// SetType sets the Type field on ProductsOptions.
func (opts *ProductsOptions) SetType(Type string) *ProductsOptions {
	opts.Type = &Type
	return opts
}

// SetActive sets the Active field on ProfileOptions.
func (opts *ProfileOptions) SetActive(Active bool) *ProfileOptions {
	opts.Active = &Active
	return opts
}

// SetActive sets the Active field on ProfilesOptions.
func (opts *ProfilesOptions) SetActive(Active bool) *ProfilesOptions {
	opts.Active = &Active
	return opts
}

// SetProfileID sets the ProfileID field on RenameProfileOptions.
func (opts *RenameProfileOptions) SetProfileID(ProfileID string) *RenameProfileOptions {
	opts.ProfileID = &ProfileID
	return opts
}

// SetName sets the Name field on RenameProfileOptions.
func (opts *RenameProfileOptions) SetName(Name string) *RenameProfileOptions {
	opts.Name = &Name
	return opts
}

// SetPortfolioID sets the PortfolioID field on ReportsOptions. Filter results by a specific profile_id
func (opts *ReportsOptions) SetPortfolioID(PortfolioID string) *ReportsOptions {
	opts.PortfolioID = &PortfolioID
	return opts
}

// SetAfter sets the After field on ReportsOptions. Filter results after a specific date
func (opts *ReportsOptions) SetAfter(After string) *ReportsOptions {
	opts.After = &After
	return opts
}

// SetLimit sets the Limit field on ReportsOptions. Limit results to a specific number
func (opts *ReportsOptions) SetLimit(Limit int) *ReportsOptions {
	opts.Limit = &Limit
	return opts
}

// SetType sets the Type field on ReportsOptions. Filter results by type of report (fills or account) - otc_fills: real
// string is otc-fills - type_1099k_transaction_history: real string is 1099-transaction-history - tax_invoice: real
// string is tax-invoice
func (opts *ReportsOptions) SetType(Type ReportType) *ReportsOptions {
	opts.Type = &Type
	return opts
}

// SetIgnoredExpired sets the IgnoredExpired field on ReportsOptions. Ignore expired results
func (opts *ReportsOptions) SetIgnoredExpired(IgnoredExpired bool) *ReportsOptions {
	opts.IgnoredExpired = &IgnoredExpired
	return opts
}

// SetLimit sets the Limit field on TradesOptions.
func (opts *TradesOptions) SetLimit(Limit int32) *TradesOptions {
	opts.Limit = &Limit
	return opts
}

// SetBefore sets the Before field on TradesOptions.
func (opts *TradesOptions) SetBefore(Before int32) *TradesOptions {
	opts.Before = &Before
	return opts
}

// SetAfter sets the After field on TradesOptions.
func (opts *TradesOptions) SetAfter(After int32) *TradesOptions {
	opts.After = &After
	return opts
}

// SetCurrency sets the Currency field on WithdrawalFeeEstimateOptions.
func (opts *WithdrawalFeeEstimateOptions) SetCurrency(Currency string) *WithdrawalFeeEstimateOptions {
	opts.Currency = &Currency
	return opts
}

// SetCryptoAddress sets the CryptoAddress field on WithdrawalFeeEstimateOptions.
func (opts *WithdrawalFeeEstimateOptions) SetCryptoAddress(CryptoAddress string) *WithdrawalFeeEstimateOptions {
	opts.CryptoAddress = &CryptoAddress
	return opts
}

func (opts *CoinbaseAccountDepositOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "coinbase_account_id", opts.CoinbaseAccountID)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CoinbaseAccountWithdrawalOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "coinbase_account_id", opts.CoinbaseAccountID)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *ConvertCurrencyOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "from", opts.From)
		internal.HTTPBodyFragment(body, "to", opts.To)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "nonce", opts.Nonce)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CreateOrderOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "type", opts.Type)
		internal.HTTPBodyFragment(body, "side", opts.Side)
		internal.HTTPBodyFragment(body, "stp", opts.STP)
		internal.HTTPBodyFragment(body, "stop", opts.Stop)
		internal.HTTPBodyFragment(body, "stop_price", opts.StopPrice)
		internal.HTTPBodyFragment(body, "price", opts.Price)
		internal.HTTPBodyFragment(body, "size", opts.Size)
		internal.HTTPBodyFragment(body, "funds", opts.Funds)
		internal.HTTPBodyFragment(body, "product_id", opts.ProductID)
		internal.HTTPBodyFragment(body, "time_in_force", opts.TimeInForce)
		internal.HTTPBodyFragment(body, "cancel_after", opts.CancelAfter)
		internal.HTTPBodyFragment(body, "post_only", opts.PostOnly)
		internal.HTTPBodyFragment(body, "client_oid", opts.ClientOid)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CreateProfileOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "name", opts.Name)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CreateProfileTransferOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "from", opts.From)
		internal.HTTPBodyFragment(body, "to", opts.To)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CreateReportOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "start_date", opts.StartDate)
		internal.HTTPBodyFragment(body, "end_date", opts.EndDate)
		internal.HTTPBodyFragment(body, "type", opts.Type)
		internal.HTTPBodyFragment(body, "year", opts.Year)
		internal.HTTPBodyFragment(body, "format", opts.Format)
		internal.HTTPBodyFragment(body, "product_id", opts.ProductID)
		internal.HTTPBodyFragment(body, "account_id", opts.AccountID)
		internal.HTTPBodyFragment(body, "email", opts.Email)
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *CryptoWithdrawalOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "crypto_address", opts.CryptoAddress)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		internal.HTTPBodyFragment(body, "destination_tag", opts.DestinationTag)
		internal.HTTPBodyFragment(body, "no_destination_tag", opts.NoDestinationTag)
		internal.HTTPBodyFragment(body, "two_factor_code", opts.TwoFactorCode)
		internal.HTTPBodyFragment(body, "nonce", opts.Nonce)
		internal.HTTPBodyFragment(body, "fee", opts.Fee)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *PaymentMethodDepositOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "payment_method_id", opts.PaymentMethodID)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *PaymentMethodWithdrawalOptions) EncodeBody() (buf io.Reader, err error) {
	if opts != nil {
		body := make(map[string]interface{})
		internal.HTTPBodyFragment(body, "profile_id", opts.ProfileID)
		internal.HTTPBodyFragment(body, "amount", opts.Amount)
		internal.HTTPBodyFragment(body, "payment_method_id", opts.PaymentMethodID)
		internal.HTTPBodyFragment(body, "currency", opts.Currency)
		raw, err := json.Marshal(body)
		if err == nil {
			buf = bytes.NewBuffer(raw)
		}
	}
	return
}

func (opts *AccountHoldsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "before", opts.Before)
		internal.HTTPQueryEncodeString(req, "after", opts.After)
		internal.HTTPQueryEncodeInt(req, "limit", opts.Limit)
	}
	return
}

func (opts *AccountLedgerOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "start_date", opts.StartDate)
		internal.HTTPQueryEncodeString(req, "end_date", opts.EndDate)
		internal.HTTPQueryEncodeInt(req, "before", opts.Before)
		internal.HTTPQueryEncodeInt(req, "after", opts.After)
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeInt(req, "limit", opts.Limit)
	}
	return
}

func (opts *AccountTransfersOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "before", opts.Before)
		internal.HTTPQueryEncodeString(req, "after", opts.After)
		internal.HTTPQueryEncodeInt(req, "limit", opts.Limit)
		internal.HTTPQueryEncodeStringer(req, "type", opts.Type)
	}
	return
}

func (opts *BookOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeInt32(req, "level", opts.Level)
	}
	return
}

func (opts *CancelOpenOrdersOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeString(req, "product_id", opts.ProductID)
	}
	return
}

func (opts *CandlesOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeStringer(req, "granularity", opts.Granularity)
		internal.HTTPQueryEncodeString(req, "start", opts.Start)
		internal.HTTPQueryEncodeString(req, "end", opts.End)
	}
	return
}

func (opts *CurrencyConversionOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
	}
	return
}

func (opts *DeleteProfileOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeString(req, "to", opts.To)
	}
	return
}

func (opts *FillsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "order_id", opts.OrderID)
		internal.HTTPQueryEncodeString(req, "product_id", opts.ProductID)
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeInt(req, "limit", opts.Limit)
		internal.HTTPQueryEncodeInt(req, "before", opts.Before)
		internal.HTTPQueryEncodeInt(req, "after", opts.After)
	}
	return
}

func (opts *OrdersOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeString(req, "product_id", opts.ProductID)
		internal.HTTPQueryEncodeString(req, "sortedBy", opts.SortedBy)
		internal.HTTPQueryEncodeString(req, "sorting", opts.Sorting)
		internal.HTTPQueryEncodeString(req, "start_date", opts.StartDate)
		internal.HTTPQueryEncodeString(req, "end_date", opts.EndDate)
		internal.HTTPQueryEncodeString(req, "before", opts.Before)
		internal.HTTPQueryEncodeString(req, "after", opts.After)
		internal.HTTPQueryEncodeInt(req, "limit", &opts.Limit)

		internal.HTTPQueryEncodeStrings(req, "status", opts.Status)
	}
	return
}

func (opts *ProductsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "type", opts.Type)
	}
	return
}

func (opts *ProfileOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeBool(req, "active", opts.Active)
	}
	return
}

func (opts *ProfilesOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeBool(req, "active", opts.Active)
	}
	return
}

func (opts *RenameProfileOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "profile_id", opts.ProfileID)
		internal.HTTPQueryEncodeString(req, "name", opts.Name)
	}
	return
}

func (opts *ReportsOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "portfolio_id", opts.PortfolioID)
		internal.HTTPQueryEncodeString(req, "after", opts.After)
		internal.HTTPQueryEncodeInt(req, "limit", opts.Limit)
		internal.HTTPQueryEncodeStringer(req, "type", opts.Type)
		internal.HTTPQueryEncodeBool(req, "ignored_expired", opts.IgnoredExpired)
	}
	return
}

func (opts *TradesOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeInt32(req, "limit", opts.Limit)
		internal.HTTPQueryEncodeInt32(req, "before", opts.Before)
		internal.HTTPQueryEncodeInt32(req, "after", opts.After)
	}
	return
}

func (opts *WithdrawalFeeEstimateOptions) EncodeQuery(req *http.Request) {
	if opts != nil {
		internal.HTTPQueryEncodeString(req, "currency", opts.Currency)
		internal.HTTPQueryEncodeString(req, "crypto_address", opts.CryptoAddress)
	}
	return
}
