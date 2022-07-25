package coinbasepro

import (
	"encoding/json"
	"time"

	"github.com/alpine-hodler/web/internal/serial"
)

// * This is a generated file, do not edit

// Account holds data for trading account from the profile of the API key
type Account struct {
	Available      string `bson:"available" json:"available" sql:"available"`
	Balance        string `bson:"balance" json:"balance" sql:"balance"`
	Currency       string `bson:"currency" json:"currency" sql:"currency"`
	Hold           string `bson:"hold" json:"hold" sql:"hold"`
	ID             string `bson:"id" json:"id" sql:"id"`
	ProfileID      string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	TradingEnabled bool   `bson:"trading_enabled" json:"trading_enabled" sql:"trading_enabled"`
}

// AccountHold represents the hold on an account that belong to the same profile as the API key. Holds are placed on an
// account for any active orders or pending withdraw requests. As an order is filled, the hold amount is updated. If an
// order is canceled, any remaining hold is removed. For withdrawals, the hold is removed after it is completed.
type AccountHold struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at" sql:"created_at"`
	ID        string    `bson:"id" json:"id" sql:"id"`
	Ref       string    `bson:"ref" json:"ref" sql:"ref"`
	Type      string    `bson:"type" json:"type" sql:"type"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at" sql:"updated_at"`
}

// AccountLedger lists ledger activity for an account. This includes anything that would affect the accounts balance -
// transfers, trades, fees, etc.
type AccountLedger struct {
	Amount    string               `bson:"amount" json:"amount" sql:"amount"`
	Balance   string               `bson:"balance" json:"balance" sql:"balance"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at" sql:"created_at"`
	Details   AccountLedgerDetails `bson:"details" json:"details" sql:"details"`
	ID        string               `bson:"id" json:"id" sql:"id"`
	Type      EntryType            `bson:"type" json:"type" sql:"type"`
}

// AccountLedgerDetails are the details for account history.
type AccountLedgerDetails struct {
	OrderID   string `bson:"order_id" json:"order_id" sql:"order_id"`
	ProductID string `bson:"product_id" json:"product_id" sql:"product_id"`
	TradeID   string `bson:"trade_id" json:"trade_id" sql:"trade_id"`
}

// AccountTransferDetails are the details for an account transfer.
type AccountTransferDetails struct {
	CoinbaseAccountID       string `bson:"coinbase_account_id" json:"coinbase_account_id" sql:"coinbase_account_id"`
	CoinbasePaymentMethodID string `bson:"coinbase_payment_method_id" json:"coinbase_payment_method_id" sql:"coinbase_payment_method_id"`
	CoinbaseTransactionID   string `bson:"coinbase_transaction_id" json:"coinbase_transaction_id" sql:"coinbase_transaction_id"`
}

// Auction is an object of data concerning a book request.
type Auction struct {
	AuctionState string    `bson:"auction_state" json:"auction_state" sql:"auction_state"`
	BestAskPrice string    `bson:"best_ask_price" json:"best_ask_price" sql:"best_ask_price"`
	BestAskSize  string    `bson:"best_ask_size" json:"best_ask_size" sql:"best_ask_size"`
	BestBidPrice string    `bson:"best_bid_price" json:"best_bid_price" sql:"best_bid_price"`
	BestBidSize  string    `bson:"best_bid_size" json:"best_bid_size" sql:"best_bid_size"`
	CanOpen      string    `bson:"can_open" json:"can_open" sql:"can_open"`
	OpenPrice    string    `bson:"open_price" json:"open_price" sql:"open_price"`
	OpenSize     string    `bson:"open_size" json:"open_size" sql:"open_size"`
	Time         time.Time `bson:"time" json:"time" sql:"time"`
}

// AvailableBalance is the available balance on the coinbase account
type AvailableBalance struct {
	Amount   string `bson:"amount" json:"amount" sql:"amount"`
	Currency string `bson:"currency" json:"currency" sql:"currency"`
	Scale    string `bson:"scale" json:"scale" sql:"scale"`
}

// Balance is the balance for picker data
type Balance struct {
	Amount   string `bson:"amount" json:"amount" sql:"amount"`
	Currency string `bson:"currency" json:"currency" sql:"currency"`
}

// BankCountry are the name and code for the bank's country associated with a wallet
type BankCountry struct {
	Code string `bson:"code" json:"code" sql:"code"`
	Name string `bson:"name" json:"name" sql:"name"`
}

// BidAsk is a slice of bytes that represents the bids or asks for a given product. The term "bid" refers to the highest
// price a buyer will pay to buy a specified number of shares of a stock at any given time. The term "ask" refers to the
// lowest price at which a seller will sell the stock. The bid price will almost always be lower than the ask or
// “offer,” price.
type BidAsk []byte

// Book is a list of open orders for a product. The amount of detail shown can be customized with the level parameter.
type Book struct {
	Asks        BidAsk  `bson:"asks" json:"asks" sql:"asks"`
	Auction     Auction `bson:"auction" json:"auction" sql:"auction"`
	AuctionMode bool    `bson:"auction_mode" json:"auction_mode" sql:"auction_mode"`
	Bids        BidAsk  `bson:"bids" json:"bids" sql:"bids"`
	Sequence    float64 `bson:"sequence" json:"sequence" sql:"sequence"`
}

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

	// ProductID is the productID for the candle, e.g. BTC-ETH. This is not through the Coinbase Pro web API and is inteded
	// for use in data layers and business logic.
	ProductID string `bson:"product_id" json:"product_id" sql:"product_id"`

	// Unix is the bucket start time as an int64 Unix value.
	Unix int64 `bson:"unix" json:"unix" sql:"unix"`

	// Volumes is the volume of trading activity during the bucket interval.
	Volume float64 `bson:"volume" json:"volume" sql:"volume"`
}

// Candles are the historic rates for a product. Rates are returned in grouped buckets. Candle schema is of the form
// `[timestamp, price_low, price_high, price_open, price_close]`
type Candles []*Candle

// CreateOrder is the server's response for placing a new order.
type CreateOrder struct {
	CreatedAt     time.Time   `bson:"created_at" json:"created_at" sql:"created_at"`
	DoneAt        time.Time   `bson:"done_at" json:"done_at" sql:"done_at"`
	DoneReason    string      `bson:"done_reason" json:"done_reason" sql:"done_reason"`
	ExpireTime    time.Time   `bson:"expire_time" json:"expire_time" sql:"expire_time"`
	FillFees      string      `bson:"fill_fees" json:"fill_fees" sql:"fill_fees"`
	FilledSize    string      `bson:"filled_size" json:"filled_size" sql:"filled_size"`
	FundingAmount string      `bson:"funding_amount" json:"funding_amount" sql:"funding_amount"`
	Funds         string      `bson:"funds" json:"funds" sql:"funds"`
	ID            string      `bson:"id" json:"id" sql:"id"`
	PostOnly      bool        `bson:"post_only" json:"post_only" sql:"post_only"`
	Price         string      `bson:"price" json:"price" sql:"price"`
	ProductID     string      `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID     string      `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	RejectReason  string      `bson:"reject_reason" json:"reject_reason" sql:"reject_reason"`
	Settled       bool        `bson:"settled" json:"settled" sql:"settled"`
	Side          Side        `bson:"side" json:"side" sql:"side"`
	Size          string      `bson:"size" json:"size" sql:"size"`
	SpecificFunds string      `bson:"specific_funds" json:"specific_funds" sql:"specific_funds"`
	Status        string      `bson:"status" json:"status" sql:"status"`
	Stop          Stop        `bson:"stop" json:"stop" sql:"stop"`
	StopPrice     string      `bson:"stop_price" json:"stop_price" sql:"stop_price"`
	TimeInForce   TimeInForce `bson:"time_in_force" json:"time_in_force" sql:"time_in_force"`
	Type          OrderType   `bson:"type" json:"type" sql:"type"`
}

// CreateReport represents information for a report created through the client.
type CreateReport struct {
	ID     string     `bson:"id" json:"id" sql:"id"`
	Status Status     `bson:"status" json:"status" sql:"status"`
	Type   ReportType `bson:"type" json:"type" sql:"type"`
}

// CryptoAccount references a crypto account that a CoinbasePaymentMethod belongs to
type CryptoAccount struct {
	ID           string `bson:"id" json:"id" sql:"id"`
	Resource     string `bson:"resource" json:"resource" sql:"resource"`
	ResourcePath string `bson:"resource_path" json:"resource_path" sql:"resource_path"`
}

// CryptoAddress is used for a one-time crypto address for depositing crypto.
type CryptoAddress struct {
	Address        string                  `bson:"address" json:"address" sql:"address"`
	AddressInfo    CryptoAddressInfo       `bson:"address_info" json:"address_info" sql:"address_info"`
	CallbackURL    string                  `bson:"callback_url" json:"callback_url" sql:"callback_url"`
	CreateAt       time.Time               `bson:"create_at" json:"create_at" sql:"create_at"`
	DepositUri     string                  `bson:"deposit_uri" json:"deposit_uri" sql:"deposit_uri"`
	DestinationTag string                  `bson:"destination_tag" json:"destination_tag" sql:"destination_tag"`
	ID             string                  `bson:"id" json:"id" sql:"id"`
	LegacyAddress  string                  `bson:"legacy_address" json:"legacy_address" sql:"legacy_address"`
	Name           string                  `bson:"name" json:"name" sql:"name"`
	Network        string                  `bson:"network" json:"network" sql:"network"`
	Resource       string                  `bson:"resource" json:"resource" sql:"resource"`
	ResourcePath   string                  `bson:"resource_path" json:"resource_path" sql:"resource_path"`
	UpdatedAt      time.Time               `bson:"updated_at" json:"updated_at" sql:"updated_at"`
	UriScheme      string                  `bson:"uri_scheme" json:"uri_scheme" sql:"uri_scheme"`
	Warnings       []*CryptoAddressWarning `bson:"warnings" json:"warnings" sql:"warnings"`
}

// CryptoAddressInfo holds info for a crypto address
type CryptoAddressInfo struct {
	Address        string `bson:"address" json:"address" sql:"address"`
	DestinationTag string `bson:"destination_tag" json:"destination_tag" sql:"destination_tag"`
}

// CryptoAddressWarning is a warning for generating a crypting address
type CryptoAddressWarning struct {
	Details  string `bson:"details" json:"details" sql:"details"`
	ImageURL string `bson:"image_url" json:"image_url" sql:"image_url"`
	Title    string `bson:"title" json:"title" sql:"title"`
}

// Currency is a currency that coinbase knows about. Not al currencies may be currently in use for trading.
type Currency struct {
	ConvertibleTo []string        `bson:"convertible_to" json:"convertible_to" sql:"convertible_to"`
	Details       CurrencyDetails `bson:"details" json:"details" sql:"details"`
	ID            string          `bson:"id" json:"id" sql:"id"`
	MaxPrecision  string          `bson:"max_precision" json:"max_precision" sql:"max_precision"`
	Message       string          `bson:"message" json:"message" sql:"message"`
	MinSize       string          `bson:"min_size" json:"min_size" sql:"min_size"`
	Name          string          `bson:"name" json:"name" sql:"name"`
	Status        string          `bson:"status" json:"status" sql:"status"`
}

// CurrencyConversion is the response that converts funds from from currency to to currency. Funds are converted on the
// from account in the profile_id profile.
type CurrencyConversion struct {
	Amount        string `bson:"amount" json:"amount" sql:"amount"`
	From          string `bson:"from" json:"from" sql:"from"`
	FromAccountID string `bson:"from_account_id" json:"from_account_id" sql:"from_account_id"`
	ID            string `bson:"id" json:"id" sql:"id"`
	To            string `bson:"to" json:"to" sql:"to"`
	ToAccountID   string `bson:"to_account_id" json:"to_account_id" sql:"to_account_id"`
}

// CurrencyDetails are the details for a currency that coinbase knows about
type CurrencyDetails struct {
	CryptoAddressLink     string   `bson:"crypto_address_link" json:"crypto_address_link" sql:"crypto_address_link"`
	CryptoTransactionLink string   `bson:"crypto_transaction_link" json:"crypto_transaction_link" sql:"crypto_transaction_link"`
	DisplayName           string   `bson:"display_name" json:"display_name" sql:"display_name"`
	GroupTypes            []string `bson:"group_types" json:"group_types" sql:"group_types"`
	MaxWithdrawalAmount   float64  `bson:"max_withdrawal_amount" json:"max_withdrawal_amount" sql:"max_withdrawal_amount"`
	MinWithdrawalAmount   float64  `bson:"min_withdrawal_amount" json:"min_withdrawal_amount" sql:"min_withdrawal_amount"`
	NetworkConfirmations  int      `bson:"network_confirmations" json:"network_confirmations" sql:"network_confirmations"`
	ProcessingTimeSeconds float64  `bson:"processing_time_seconds" json:"processing_time_seconds" sql:"processing_time_seconds"`
	PushPaymentMethods    []string `bson:"push_payment_methods" json:"push_payment_methods" sql:"push_payment_methods"`
	SortOrder             int      `bson:"sort_order" json:"sort_order" sql:"sort_order"`
	Symbol                string   `bson:"symbol" json:"symbol" sql:"symbol"`
	Type                  string   `bson:"type" json:"type" sql:"type"`
}

// CurrencyTransferLimit encapsulates ACH data for a currency via Max/Remaining amounts.
type CurrencyTransferLimit struct {
	Max       float64 `bson:"max" json:"max" sql:"max"`
	Remaining float64 `bson:"remaining" json:"remaining" sql:"remaining"`
}

// CurrencyTransferLimits encapsulates ACH data for many currencies.
type CurrencyTransferLimits map[string]CurrencyTransferLimit

// Deposit is the response for deposited funds from a www.coinbase.com wallet to the specified profile_id.
type Deposit struct {
	Amount   string `bson:"amount" json:"amount" sql:"amount"`
	Currency string `bson:"currency" json:"currency" sql:"currency"`
	Fee      string `bson:"fee" json:"fee" sql:"fee"`
	ID       string `bson:"id" json:"id" sql:"id"`
	PayoutAt string `bson:"payout_at" json:"payout_at" sql:"payout_at"`
	Subtotal string `bson:"subtotal" json:"subtotal" sql:"subtotal"`
}

// ExchangeLimits represents exchange limit information for a single user.
type ExchangeLimits struct {
	LimitCurrency  string         `bson:"limit_currency" json:"limit_currency" sql:"limit_currency"`
	TransferLimits TransferLimits `bson:"transfer_limits" json:"transfer_limits" sql:"transfer_limits"`
}

// Fees are fees rates and 30 days trailing volume.
type Fees struct {
	MakerFeeRate string `bson:"maker_fee_rate" json:"maker_fee_rate" sql:"maker_fee_rate"`
	TakerFeeRate string `bson:"taker_fee_rate" json:"taker_fee_rate" sql:"taker_fee_rate"`
	UsdVolume    string `bson:"usd_volume" json:"usd_volume" sql:"usd_volume"`
}

// FIATAccount references a FIAT account thata CoinbasePaymentMethod belongs to
type FIATAccount struct {
	ID           string `bson:"id" json:"id" sql:"id"`
	Resource     string `bson:"resource" json:"resource" sql:"resource"`
	ResourcePath string `bson:"resource_path" json:"resource_path" sql:"resource_path"`
}

// TODO: Get fill description
type Fill struct {
	Fee       string `bson:"fee" json:"fee" sql:"fee"`
	Liquidity string `bson:"liquidity" json:"liquidity" sql:"liquidity"`
	OrderID   string `bson:"order_id" json:"order_id" sql:"order_id"`
	Price     string `bson:"price" json:"price" sql:"price"`
	ProductID string `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID string `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	Settled   bool   `bson:"settled" json:"settled" sql:"settled"`
	Side      string `bson:"side" json:"side" sql:"side"`
	Size      string `bson:"size" json:"size" sql:"size"`
	TradeID   int    `bson:"trade_id" json:"trade_id" sql:"trade_id"`
	UsdVolume string `bson:"usd_volume" json:"usd_volume" sql:"usd_volume"`
	UserID    string `bson:"user_id" json:"user_id" sql:"user_id"`
}

// TODO
type Flags struct{}

// Limits defines limits for a payment method
type Limits struct {
	Name string `bson:"name" json:"name" sql:"name"`
	Type string `bson:"type" json:"type" sql:"type"`
}

// Oracle is cryptographically signed price-info ready to be posted on-chain using Compound's Open Oracle smart
// contract.
type Oracle struct {
	// Messages are an array contains abi-encoded values [kind, timestamp, key, value], where kind always equals to
	// 'prices', timestamp is the time when the price was obtained, key is asset ticker (e.g. 'eth') and value is asset
	// price
	Messages []string `bson:"messages" json:"messages" sql:"messages"`

	// Prices contains human-readable asset prices
	Prices OraclePrices `bson:"prices" json:"prices" sql:"prices"`

	// Signatures are an array of Ethereum-compatible ECDSA signatures for each message
	Signatures []string `bson:"signatures" json:"signatures" sql:"signatures"`

	// Timestamp indicates when the latest datapoint was obtained
	Timestamp time.Time `bson:"timestamp" json:"timestamp" sql:"timestamp"`
}

// OraclePrices contain human-readable asset prices.
type OraclePrices struct {
	AdditionalProp string `bson:"additional_prop" json:"additionalProp" sql:"additional_prop"`
}

// Order is an open order.
type Order struct {
	CreatedAt      time.Time   `bson:"created_at" json:"created_at" sql:"created_at"`
	DoneAt         time.Time   `bson:"done_at" json:"done_at" sql:"done_at"`
	DoneReason     string      `bson:"done_reason" json:"done_reason" sql:"done_reason"`
	ExecutedValue  string      `bson:"executed_value" json:"executed_value" sql:"executed_value"`
	ExpireTime     time.Time   `bson:"expire_time" json:"expire_time" sql:"expire_time"`
	FillFees       string      `bson:"fill_fees" json:"fill_fees" sql:"fill_fees"`
	FilledSize     string      `bson:"filled_size" json:"filled_size" sql:"filled_size"`
	FundingAmount  string      `bson:"funding_amount" json:"funding_amount" sql:"funding_amount"`
	Funds          string      `bson:"funds" json:"funds" sql:"funds"`
	ID             string      `bson:"id" json:"id" sql:"id"`
	PostOnly       bool        `bson:"post_only" json:"post_only" sql:"post_only"`
	Price          string      `bson:"price" json:"price" sql:"price"`
	ProductID      string      `bson:"product_id" json:"product_id" sql:"product_id"`
	RejectReason   string      `bson:"reject_reason" json:"reject_reason" sql:"reject_reason"`
	Settled        bool        `bson:"settled" json:"settled" sql:"settled"`
	Side           Side        `bson:"side" json:"side" sql:"side"`
	Size           string      `bson:"size" json:"size" sql:"size"`
	SpecifiedFunds string      `bson:"specified_funds" json:"specified_funds" sql:"specified_funds"`
	Status         string      `bson:"status" json:"status" sql:"status"`
	Stop           string      `bson:"stop" json:"stop" sql:"stop"`
	StopPrice      string      `bson:"stop_price" json:"stop_price" sql:"stop_price"`
	TimeInForce    TimeInForce `bson:"time_in_force" json:"time_in_force" sql:"time_in_force"`
	Type           OrderType   `bson:"type" json:"type" sql:"type"`
}

// PaymentMethod is a payment method used on coinbase
type PaymentMethod struct {
	AllowBuy           bool                `bson:"allow_buy" json:"allow_buy" sql:"allow_buy"`
	AllowDeposit       bool                `bson:"allow_deposit" json:"allow_deposit" sql:"allow_deposit"`
	AllowSell          bool                `bson:"allow_sell" json:"allow_sell" sql:"allow_sell"`
	AllowWithdraw      bool                `bson:"allow_withdraw" json:"allow_withdraw" sql:"allow_withdraw"`
	AvailableBalance   AvailableBalance    `bson:"available_balance" json:"available_balance" sql:"available_balance"`
	CdvStatus          string              `bson:"cdv_status" json:"cdv_status" sql:"cdv_status"`
	CreateAt           time.Time           `bson:"create_at" json:"create_at" sql:"create_at"`
	CryptoAccount      CryptoAccount       `bson:"crypto_account" json:"crypto_account" sql:"crypto_account"`
	Currency           string              `bson:"currency" json:"currency" sql:"currency"`
	FIATAccount        FIATAccount         `bson:"fiat_account" json:"fiat_account" sql:"fiat_account"`
	HoldBusinessDays   int                 `bson:"hold_business_days" json:"hold_business_days" sql:"hold_business_days"`
	HoldDays           int                 `bson:"hold_days" json:"hold_days" sql:"hold_days"`
	ID                 string              `bson:"id" json:"id" sql:"id"`
	InstantBuy         bool                `bson:"instant_buy" json:"instant_buy" sql:"instant_buy"`
	InstantSale        bool                `bson:"instant_sale" json:"instant_sale" sql:"instant_sale"`
	Limits             Limits              `bson:"limits" json:"limits" sql:"limits"`
	Name               string              `bson:"name" json:"name" sql:"name"`
	PickerData         PickerData          `bson:"picker_data" json:"picker_data" sql:"picker_data"`
	PrimaryBuy         bool                `bson:"primary_buy" json:"primary_buy" sql:"primary_buy"`
	PrimarySell        bool                `bson:"primary_sell" json:"primary_sell" sql:"primary_sell"`
	RecurringOptions   []*RecurringOptions `bson:"recurring_options" json:"recurring_options" sql:"recurring_options"`
	Resource           string              `bson:"resource" json:"resource" sql:"resource"`
	ResourcePath       string              `bson:"resource_path" json:"resource_path" sql:"resource_path"`
	Type               string              `bson:"type" json:"type" sql:"type"`
	UpdatedAt          time.Time           `bson:"updated_at" json:"updated_at" sql:"updated_at"`
	VerificationMethod string              `bson:"verification_method" json:"verification_method" sql:"verification_method"`
	Verified           bool                `bson:"verified" json:"verified" sql:"verified"`
}

// PickerData ??
type PickerData struct {
	AccountName           string  `bson:"account_name" json:"account_name" sql:"account_name"`
	AccountNumber         string  `bson:"account_number" json:"account_number" sql:"account_number"`
	AccountType           string  `bson:"account_type" json:"account_type" sql:"account_type"`
	Balance               Balance `bson:"balance" json:"balance" sql:"balance"`
	BankName              string  `bson:"bank_name" json:"bank_name" sql:"bank_name"`
	BranchName            string  `bson:"branch_name" json:"branch_name" sql:"branch_name"`
	CustomerName          string  `bson:"customer_name" json:"customer_name" sql:"customer_name"`
	Iban                  string  `bson:"iban" json:"iban" sql:"iban"`
	IconURL               string  `bson:"icon_url" json:"icon_url" sql:"icon_url"`
	InstitutionCode       string  `bson:"institution_code" json:"institution_code" sql:"institution_code"`
	InstitutionIdentifier string  `bson:"institution_identifier" json:"institution_identifier" sql:"institution_identifier"`
	InstitutionName       string  `bson:"institution_name" json:"institution_name" sql:"institution_name"`
	PaypalEmail           string  `bson:"paypal_email" json:"paypal_email" sql:"paypal_email"`
	PaypalOwner           string  `bson:"paypal_owner" json:"paypal_owner" sql:"paypal_owner"`
	RoutingNumber         string  `bson:"routing_number" json:"routing_number" sql:"routing_number"`
	SWIFT                 string  `bson:"swift" json:"swift" sql:"swift"`
	Symbol                string  `bson:"symbol" json:"symbol" sql:"symbol"`
}

// Product represents a currency pair available for trading.
type Product struct {
	AuctionMode           bool   `bson:"auction_mode" json:"auction_mode" sql:"auction_mode"`
	BaseCurrency          string `bson:"base_currency" json:"base_currency" sql:"base_currency"`
	BaseIncrement         string `bson:"base_increment" json:"base_increment" sql:"base_increment"`
	BaseMaxSize           string `bson:"base_max_size" json:"base_max_size" sql:"base_max_size"`
	BaseMinSize           string `bson:"base_min_size" json:"base_min_size" sql:"base_min_size"`
	CancelOnly            bool   `bson:"cancel_only" json:"cancel_only" sql:"cancel_only"`
	DisplayName           string `bson:"display_name" json:"display_name" sql:"display_name"`
	FxStablecoin          bool   `bson:"fx_stablecoin" json:"fx_stablecoin" sql:"fx_stablecoin"`
	ID                    string `bson:"id" json:"id" sql:"id"`
	LimitOnly             bool   `bson:"limit_only" json:"limit_only" sql:"limit_only"`
	MarginEnabled         bool   `bson:"margin_enabled" json:"margin_enabled" sql:"margin_enabled"`
	MaxMarketFunds        string `bson:"max_market_funds" json:"max_market_funds" sql:"max_market_funds"`
	MaxSlippagePercentage string `bson:"max_slippage_percentage" json:"max_slippage_percentage" sql:"max_slippage_percentage"`
	MinMarketFunds        string `bson:"min_market_funds" json:"min_market_funds" sql:"min_market_funds"`
	PostOnly              bool   `bson:"post_only" json:"post_only" sql:"post_only"`
	QuoteCurrency         string `bson:"quote_currency" json:"quote_currency" sql:"quote_currency"`
	QuoteIncrement        string `bson:"quote_increment" json:"quote_increment" sql:"quote_increment"`
	Status                Status `bson:"status" json:"status" sql:"status"`
	StatusMessage         string `bson:"status_message" json:"status_message" sql:"status_message"`
	TradingDisabled       bool   `bson:"trading_disabled" json:"trading_disabled" sql:"trading_disabled"`
}

// ProductStats are 30day and 24hour stats for a product.
type ProductStats struct {
	High        string `bson:"high" json:"high" sql:"high"`
	Last        string `bson:"last" json:"last" sql:"last"`
	Low         string `bson:"low" json:"low" sql:"low"`
	Open        string `bson:"open" json:"open" sql:"open"`
	Volume      string `bson:"volume" json:"volume" sql:"volume"`
	Volume30day string `bson:"volume_30day" json:"volume_30day" sql:"volume_30day"`
}

// ProductTicker is a snapshot information about the last trade (tick), best bid/ask and 24h volume.
type ProductTicker struct {
	Ask     string    `bson:"ask" json:"ask" sql:"ask"`
	Bid     string    `bson:"bid" json:"bid" sql:"bid"`
	Price   string    `bson:"price" json:"price" sql:"price"`
	Size    string    `bson:"size" json:"size" sql:"size"`
	Time    time.Time `bson:"time" json:"time" sql:"time"`
	TradeID int       `bson:"trade_id" json:"trade_id" sql:"trade_id"`
	Volume  string    `bson:"volume" json:"volume" sql:"volume"`
}

// Profile represents a profile to interact with the API.
type Profile struct {
	Active    bool      `bson:"active" json:"active" sql:"active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at" sql:"created_at"`
	HasMargin bool      `bson:"has_margin" json:"has_margin" sql:"has_margin"`
	ID        string    `bson:"id" json:"id" sql:"id"`
	IsDefault bool      `bson:"is_default" json:"is_default" sql:"is_default"`
	Name      string    `bson:"name" json:"name" sql:"name"`
	UserID    string    `bson:"user_id" json:"user_id" sql:"user_id"`
}

// RecurringOptions ??
type RecurringOptions struct {
	Label  string `bson:"label" json:"label" sql:"label"`
	Period string `bson:"period" json:"period" sql:"period"`
}

// Report represents a list of past fills/account reports.
type Report struct {
	CreatedAt time.Time     `bson:"created_at" json:"created_at" sql:"created_at"`
	ExpiresAt time.Time     `bson:"expires_at" json:"expires_at" sql:"expires_at"`
	FileCount string        `bson:"file_count" json:"file_count" sql:"file_count"`
	FileURL   string        `bson:"file_url" json:"file_url" sql:"file_url"`
	ID        string        `bson:"id" json:"id" sql:"id"`
	Params    ReportsParams `bson:"params" json:"params" sql:"params"`
	Status    Status        `bson:"status" json:"status" sql:"status"`
	Type      ReportType    `bson:"type" json:"type" sql:"type"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at" sql:"updated_at"`
	UserID    string        `bson:"user_id" json:"user_id" sql:"user_id"`
}

// TODO
type ReportsParams struct {
	AccountID    string     `bson:"account_id" json:"account_id" sql:"account_id"`
	Email        string     `bson:"email" json:"email" sql:"email"`
	EndDate      time.Time  `bson:"end_date" json:"end_date" sql:"end_date"`
	Format       FileFormat `bson:"format" json:"format" sql:"format"`
	NewYorkState bool       `bson:"new_york_state" json:"new_york_state" sql:"new_york_state"`
	ProductID    string     `bson:"product_id" json:"product_id" sql:"product_id"`
	ProfileID    string     `bson:"profile_id" json:"profile_id" sql:"profile_id"`
	StartDate    time.Time  `bson:"start_date" json:"start_date" sql:"start_date"`
	User         User       `bson:"user" json:"user" sql:"user"`
}

// TODO
type Role struct{}

// SEPADepositInformation information regarding a wallet's deposits. A SEPA credit transfer is a single transfer of
// Euros from one person or organisation to another. For example, this could be to pay the deposit for a holiday rental
// or to settle an invoice. A SEPA direct debit is a recurring payment, for example to pay monthly rent or for a service
// like a mobile phone contract.
type SEPADepositInformation struct {
	AccountAddress string      `bson:"account_address" json:"account_address" sql:"account_address"`
	AccountName    string      `bson:"account_name" json:"account_name" sql:"account_name"`
	BankAddress    string      `bson:"bank_address" json:"bank_address" sql:"bank_address"`
	BankCountry    BankCountry `bson:"bank_country" json:"bank_country" sql:"bank_country"`
	BankName       string      `bson:"bank_name" json:"bank_name" sql:"bank_name"`
	Iban           string      `bson:"iban" json:"iban" sql:"iban"`
	Reference      string      `bson:"reference" json:"reference" sql:"reference"`
	SWIFT          string      `bson:"swift" json:"swift" sql:"swift"`
}

// SWIFTDepositInformation information regarding a wallet's deposits. SWIFT stands for Society for Worldwide Interbank
// Financial Telecommunications. Basically, it's a computer network that connects over 900 banks around the world – and
// enables them to transfer money. ING is part of this network. There is no fee for accepting deposits into your account
// with ING.
type SWIFTDepositInformation struct {
	AccountAddress string      `bson:"account_address" json:"account_address" sql:"account_address"`
	AccountName    string      `bson:"account_name" json:"account_name" sql:"account_name"`
	AccountNumber  string      `bson:"account_number" json:"account_number" sql:"account_number"`
	BankAddress    string      `bson:"bank_address" json:"bank_address" sql:"bank_address"`
	BankCountry    BankCountry `bson:"bank_country" json:"bank_country" sql:"bank_country"`
	BankName       string      `bson:"bank_name" json:"bank_name" sql:"bank_name"`
	Reference      string      `bson:"reference" json:"reference" sql:"reference"`
}

// Ticker is real-time price updates every time a match happens. It batches updates in case of cascading matches,
// greatly reducing bandwidth requirements.
type Ticker struct {
	BestAsk   string    `bson:"best_ask" json:"best_ask" sql:"best_ask"`
	BestBid   string    `bson:"best_bid" json:"best_bid" sql:"best_bid"`
	LastSize  string    `bson:"last_size" json:"last_size" sql:"last_size"`
	Price     string    `bson:"price" json:"price" sql:"price"`
	ProductID string    `bson:"product_id" json:"product_id" sql:"product_id"`
	Sequence  int       `bson:"sequence" json:"sequence" sql:"sequence"`
	Side      string    `bson:"side" json:"side" sql:"side"`
	Time      time.Time `bson:"time" json:"time" sql:"time"`
	TradeID   int       `bson:"trade_id" json:"trade_id" sql:"trade_id"`
	Type      string    `bson:"type" json:"type" sql:"type"`
}

// Trade is the list the latest trades for a product.
type Trade struct {
	Price   string    `bson:"price" json:"price" sql:"price"`
	Side    Side      `bson:"side" json:"side" sql:"side"`
	Size    string    `bson:"size" json:"size" sql:"size"`
	Time    time.Time `bson:"time" json:"time" sql:"time"`
	TradeID int32     `bson:"trade_id" json:"trade_id" sql:"trade_id"`
}

// Transfer will lists past withdrawals and deposits for an account.
type Transfer struct {
	Amount      string                 `bson:"amount" json:"amount" sql:"amount"`
	CanceledAt  time.Time              `bson:"canceled_at" json:"canceled_at" sql:"canceled_at"`
	CompletedAt time.Time              `bson:"completed_at" json:"completed_at" sql:"completed_at"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at" sql:"created_at"`
	Details     AccountTransferDetails `bson:"details" json:"details" sql:"details"`
	ID          string                 `bson:"id" json:"id" sql:"id"`
	ProcessedAt time.Time              `bson:"processed_at" json:"processed_at" sql:"processed_at"`
	Type        string                 `bson:"type" json:"type" sql:"type"`
	UserNonce   string                 `bson:"user_nonce" json:"user_nonce" sql:"user_nonce"`
}

// TODO
type TransferLimits struct {
	ACH                   CurrencyTransferLimits `bson:"ach" json:"ach" sql:"ach"`
	ACHNoBalance          CurrencyTransferLimits `bson:"ach_no_balance" json:"ach_no_balance" sql:"ach_no_balance"`
	Buy                   CurrencyTransferLimits `bson:"buy" json:"buy" sql:"buy"`
	CreditDebitCard       CurrencyTransferLimits `bson:"credit_debit_card" json:"credit_debit_card" sql:"credit_debit_card"`
	ExchangeWithdraw      CurrencyTransferLimits `bson:"exchange_withdraw" json:"exchange_withdraw" sql:"exchange_withdraw"`
	IdealDeposit          CurrencyTransferLimits `bson:"ideal_deposit" json:"ideal_deposit" sql:"ideal_deposit"`
	InstanceACHWithdrawal CurrencyTransferLimits `bson:"instance_ach_withdrawal" json:"instance_ach_withdrawal" sql:"instance_ach_withdrawal"`
	PaypalBuy             CurrencyTransferLimits `bson:"paypal_buy" json:"paypal_buy" sql:"paypal_buy"`
	PaypalWithdrawal      CurrencyTransferLimits `bson:"paypal_withdrawal" json:"paypal_withdrawal" sql:"paypal_withdrawal"`
	Secure3dBuy           CurrencyTransferLimits `bson:"secure3d_buy" json:"secure3d_buy" sql:"secure3d_buy"`
	Sell                  CurrencyTransferLimits `bson:"sell" json:"sell" sql:"sell"`
	SofortDeposit         CurrencyTransferLimits `bson:"sofort_deposit" json:"sofort_deposit" sql:"sofort_deposit"`
}

// UKDepositInformation information regarding a wallet's deposits.
type UKDepositInformation struct {
	AccountAddress string      `bson:"account_address" json:"account_address" sql:"account_address"`
	AccountName    string      `bson:"account_name" json:"account_name" sql:"account_name"`
	AccountNumber  string      `bson:"account_number" json:"account_number" sql:"account_number"`
	BankAddress    string      `bson:"bank_address" json:"bank_address" sql:"bank_address"`
	BankCountry    BankCountry `bson:"bank_country" json:"bank_country" sql:"bank_country"`
	BankName       string      `bson:"bank_name" json:"bank_name" sql:"bank_name"`
	Reference      string      `bson:"reference" json:"reference" sql:"reference"`
}

// TODO
type User struct {
	ActiveAt                  time.Time       `bson:"active_at" json:"active_at" sql:"active_at"`
	CbDataFromCache           bool            `bson:"cb_data_from_cache" json:"cb_data_from_cache" sql:"cb_data_from_cache"`
	CreatedAt                 time.Time       `bson:"created_at" json:"created_at" sql:"created_at"`
	Details                   UserDetails     `bson:"details" json:"details" sql:"details"`
	Flags                     Flags           `bson:"flags" json:"flags" sql:"flags"`
	FulfillsNewRequirements   bool            `bson:"fulfills_new_requirements" json:"fulfills_new_requirements" sql:"fulfills_new_requirements"`
	HasClawbackPaymentPending bool            `bson:"has_clawback_payment_pending" json:"has_clawback_payment_pending" sql:"has_clawback_payment_pending"`
	HasDefault                bool            `bson:"has_default" json:"has_default" sql:"has_default"`
	HasRestrictedAssets       bool            `bson:"has_restricted_assets" json:"has_restricted_assets" sql:"has_restricted_assets"`
	ID                        string          `bson:"id" json:"id" sql:"id"`
	IsBanned                  bool            `bson:"is_banned" json:"is_banned" sql:"is_banned"`
	LegalName                 string          `bson:"legal_name" json:"legal_name" sql:"legal_name"`
	Name                      string          `bson:"name" json:"name" sql:"name"`
	Preferences               UserPreferences `bson:"preferences" json:"preferences" sql:"preferences"`
	Roles                     []*Role         `bson:"roles" json:"roles" sql:"roles"`
	StateCode                 string          `bson:"state_code" json:"state_code" sql:"state_code"`
	TermsAccepted             time.Time       `bson:"terms_accepted" json:"terms_accepted" sql:"terms_accepted"`
	TwoFactorMethod           string          `bson:"two_factor_method" json:"two_factor_method" sql:"two_factor_method"`
	UserType                  string          `bson:"user_type" json:"user_type" sql:"user_type"`
}

// TODO
type UserDetails struct{}

// TODO
type UserPreferences struct{}

// Wallet represents a user's available Coinbase wallet (These are the wallets/accounts that are used for buying and
// selling on www.coinbase.com)
type Wallet struct {
	Active                  bool                    `bson:"active" json:"active" sql:"active"`
	AvailableOnConsumer     bool                    `bson:"available_on_consumer" json:"available_on_consumer" sql:"available_on_consumer"`
	Balance                 string                  `bson:"balance" json:"balance" sql:"balance"`
	Currency                string                  `bson:"currency" json:"currency" sql:"currency"`
	DestinationTagName      string                  `bson:"destination_tag_name" json:"destination_tag_name" sql:"destination_tag_name"`
	DestinationTagRegex     string                  `bson:"destination_tag_regex" json:"destination_tag_regex" sql:"destination_tag_regex"`
	HoldBalance             string                  `bson:"hold_balance" json:"hold_balance" sql:"hold_balance"`
	HoldCurrency            string                  `bson:"hold_currency" json:"hold_currency" sql:"hold_currency"`
	ID                      string                  `bson:"id" json:"id" sql:"id"`
	Name                    string                  `bson:"name" json:"name" sql:"name"`
	Primary                 bool                    `bson:"primary" json:"primary" sql:"primary"`
	Ready                   bool                    `bson:"ready" json:"ready" sql:"ready"`
	SEPADepositInformation  SEPADepositInformation  `bson:"sepa_deposit_information" json:"sepa_deposit_information" sql:"sepa_deposit_information"`
	SWIFTDepositInformation SWIFTDepositInformation `bson:"swift_deposit_information" json:"swift_deposit_information" sql:"swift_deposit_information"`
	Type                    string                  `bson:"type" json:"type" sql:"type"`
	UKDepositInformation    UKDepositInformation    `bson:"uk_deposit_information" json:"uk_deposit_information" sql:"uk_deposit_information"`
	WireDepositInformation  WireDepositInformation  `bson:"wire_deposit_information" json:"wire_deposit_information" sql:"wire_deposit_information"`
}

// WireDepositInformation information regarding a wallet's deposits
type WireDepositInformation struct {
	AccountAddress string      `bson:"account_address" json:"account_address" sql:"account_address"`
	AccountName    string      `bson:"account_name" json:"account_name" sql:"account_name"`
	AccountNumber  string      `bson:"account_number" json:"account_number" sql:"account_number"`
	BankAddress    string      `bson:"bank_address" json:"bank_address" sql:"bank_address"`
	BankCountry    BankCountry `bson:"bank_country" json:"bank_country" sql:"bank_country"`
	BankName       string      `bson:"bank_name" json:"bank_name" sql:"bank_name"`
	Reference      string      `bson:"reference" json:"reference" sql:"reference"`
	RoutingNumber  string      `bson:"routing_number" json:"routing_number" sql:"routing_number"`
}

// Withdrawal is data concerning withdrawing funds from the specified profile_id to a www.coinbase.com wallet.
type Withdrawal struct {
	Amount   string `bson:"amount" json:"amount" sql:"amount"`
	Currency string `bson:"currency" json:"currency" sql:"currency"`
	Fee      string `bson:"fee" json:"fee" sql:"fee"`
	ID       string `bson:"id" json:"id" sql:"id"`
	PayoutAt string `bson:"payout_at" json:"payout_at" sql:"payout_at"`
	Subtotal string `bson:"subtotal" json:"subtotal" sql:"subtotal"`
}

// WithdrawalFeeEstimate is a fee estimate for the crypto withdrawal to crypto address
type WithdrawalFeeEstimate struct {
	Fee float64 `bson:"fee" json:"fee" sql:"fee"`
}

// UnmarshalJSON will deserialize bytes into a BidAsk model
func (ba *BidAsk) UnmarshalJSON(b []byte) error { *ba = BidAsk(b); return nil }

// UnmarshalJSON will deserialize bytes into a Candles model
func (candles *Candles) UnmarshalJSON(bytes []byte) error {
	var raw [][]float64
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
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

// UnmarshalJSON will deserialize bytes into a Oracle model
func (Oracle *Oracle) UnmarshalJSON(d []byte) error {
	const (
		timestampJSONTag  = "timestamp"
		messagesJSONTag   = "messages"
		signaturesJSONTag = "signatures"
		pricesJSONTag     = "prices"
	)
	data, err := serial.NewJSONTransform(d)
	if err != nil {
		return err
	}
	Oracle.Prices = OraclePrices{}
	if err := data.UnmarshalStruct(pricesJSONTag, &Oracle.Prices); err != nil {
		return err
	}
	data.UnmarshalStringSlice(messagesJSONTag, &Oracle.Messages)
	data.UnmarshalStringSlice(signaturesJSONTag, &Oracle.Signatures)
	err = data.UnmarshalUnixString(timestampJSONTag, &Oracle.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON will deserialize bytes into a Transfer model
func (Transfer *Transfer) UnmarshalJSON(d []byte) error {
	const (
		IDJSONTag          = "id"
		typeJSONTag        = "type"
		createdAtJSONTag   = "created_at"
		completedAtJSONTag = "completed_at"
		canceledAtJSONTag  = "canceled_at"
		processedAtJSONTag = "processed_at"
		amountJSONTag      = "amount"
		userNonceJSONTag   = "user_nonce"
		detailsJSONTag     = "details"
	)
	data, err := serial.NewJSONTransform(d)
	if err != nil {
		return err
	}
	Transfer.Details = AccountTransferDetails{}
	if err := data.UnmarshalStruct(detailsJSONTag, &Transfer.Details); err != nil {
		return err
	}
	data.UnmarshalString(IDJSONTag, &Transfer.ID)
	data.UnmarshalString(amountJSONTag, &Transfer.Amount)
	data.UnmarshalString(typeJSONTag, &Transfer.Type)
	data.UnmarshalString(userNonceJSONTag, &Transfer.UserNonce)
	err = data.UnmarshalTime(coinbaseTimeLayout1, canceledAtJSONTag, &Transfer.CanceledAt)
	if err != nil {
		return err
	}
	err = data.UnmarshalTime(coinbaseTimeLayout1, completedAtJSONTag, &Transfer.CompletedAt)
	if err != nil {
		return err
	}
	err = data.UnmarshalTime(coinbaseTimeLayout1, createdAtJSONTag, &Transfer.CreatedAt)
	if err != nil {
		return err
	}
	err = data.UnmarshalTime(coinbaseTimeLayout1, processedAtJSONTag, &Transfer.ProcessedAt)
	if err != nil {
		return err
	}
	return nil
}
