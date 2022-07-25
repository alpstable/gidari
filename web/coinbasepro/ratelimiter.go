package coinbasepro

import (
	"time"

	"golang.org/x/time/rate"
)

// * This is a generated file, do not edit

type ratelimiter uint8

const (
	_ ratelimiter = iota
	AccountHoldsRatelimiter
	AccountLedgerRatelimiter
	AccountRatelimiter
	AccountTransfersRatelimiter
	AccountsRatelimiter
	BookRatelimiter
	CancelOpenOrdersRatelimiter
	CancelOrderRatelimiter
	CandlesRatelimiter
	CoinbaseAccountDepositRatelimiter
	CoinbaseAccountWithdrawalRatelimiter
	ConvertCurrencyRatelimiter
	CreateOrderRatelimiter
	CreateProfileRatelimiter
	CreateProfileTransferRatelimiter
	CreateReportRatelimiter
	CryptoWithdrawalRatelimiter
	CurrenciesRatelimiter
	CurrencyConversionRatelimiter
	CurrencyRatelimiter
	DeleteProfileRatelimiter
	ExchangeLimitsRatelimiter
	FeesRatelimiter
	FillsRatelimiter
	GenerateCryptoAddressRatelimiter
	OrderRatelimiter
	OrdersRatelimiter
	PaymentMethodDepositRatelimiter
	PaymentMethodWithdrawalRatelimiter
	PaymentMethodsRatelimiter
	ProductRatelimiter
	ProductStatsRatelimiter
	ProductTickerRatelimiter
	ProductsRatelimiter
	ProfileRatelimiter
	ProfilesRatelimiter
	RenameProfileRatelimiter
	ReportRatelimiter
	ReportsRatelimiter
	SignedPricesRatelimiter
	TradesRatelimiter
	TransferRatelimiter
	TransfersRatelimiter
	WalletsRatelimiter
	WithdrawalFeeEstimateRatelimiter
)

var ratelimiters = [uint8(46)]*rate.Limiter{}

func init() {
	ratelimiters[AccountHoldsRatelimiter] = nil
	ratelimiters[AccountLedgerRatelimiter] = nil
	ratelimiters[AccountRatelimiter] = nil
	ratelimiters[AccountTransfersRatelimiter] = nil
	ratelimiters[AccountsRatelimiter] = nil
	ratelimiters[BookRatelimiter] = nil
	ratelimiters[CancelOpenOrdersRatelimiter] = nil
	ratelimiters[CancelOrderRatelimiter] = nil
	ratelimiters[CandlesRatelimiter] = nil
	ratelimiters[CoinbaseAccountDepositRatelimiter] = nil
	ratelimiters[CoinbaseAccountWithdrawalRatelimiter] = nil
	ratelimiters[ConvertCurrencyRatelimiter] = nil
	ratelimiters[CreateOrderRatelimiter] = nil
	ratelimiters[CreateProfileRatelimiter] = nil
	ratelimiters[CreateProfileTransferRatelimiter] = nil
	ratelimiters[CreateReportRatelimiter] = nil
	ratelimiters[CryptoWithdrawalRatelimiter] = nil
	ratelimiters[CurrenciesRatelimiter] = nil
	ratelimiters[CurrencyConversionRatelimiter] = nil
	ratelimiters[CurrencyRatelimiter] = nil
	ratelimiters[DeleteProfileRatelimiter] = nil
	ratelimiters[ExchangeLimitsRatelimiter] = nil
	ratelimiters[FeesRatelimiter] = nil
	ratelimiters[FillsRatelimiter] = nil
	ratelimiters[GenerateCryptoAddressRatelimiter] = nil
	ratelimiters[OrderRatelimiter] = nil
	ratelimiters[OrdersRatelimiter] = nil
	ratelimiters[PaymentMethodDepositRatelimiter] = nil
	ratelimiters[PaymentMethodWithdrawalRatelimiter] = nil
	ratelimiters[PaymentMethodsRatelimiter] = nil
	ratelimiters[ProductRatelimiter] = nil
	ratelimiters[ProductStatsRatelimiter] = nil
	ratelimiters[ProductTickerRatelimiter] = nil
	ratelimiters[ProductsRatelimiter] = nil
	ratelimiters[ProfileRatelimiter] = nil
	ratelimiters[ProfilesRatelimiter] = nil
	ratelimiters[RenameProfileRatelimiter] = nil
	ratelimiters[ReportRatelimiter] = nil
	ratelimiters[ReportsRatelimiter] = nil
	ratelimiters[SignedPricesRatelimiter] = nil
	ratelimiters[TradesRatelimiter] = nil
	ratelimiters[TransferRatelimiter] = nil
	ratelimiters[TransfersRatelimiter] = nil
	ratelimiters[WalletsRatelimiter] = nil
	ratelimiters[WithdrawalFeeEstimateRatelimiter] = nil
}

// getRateLimiter will load the rate limiter for a specific request, lazy loaded.
func getRateLimiter(rl ratelimiter, b int) *rate.Limiter {
	if ratelimiters[rl] == nil {
		ratelimiters[rl] = rate.NewLimiter(rate.Every(1*time.Second), b)
	}
	return ratelimiters[rl]
}
