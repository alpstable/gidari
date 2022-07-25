package coinbasepro

import "path"

// * This is a generated file, do not edit

type rawPath uint8

const (
	_ rawPath = iota
	AccountHoldsPath
	AccountLedgerPath
	AccountPath
	AccountTransfersPath
	AccountsPath
	BookPath
	CancelOpenOrdersPath
	CancelOrderPath
	CandlesPath
	CoinbaseAccountDepositPath
	CoinbaseAccountWithdrawalPath
	ConvertCurrencyPath
	CreateOrderPath
	CreateProfilePath
	CreateProfileTransferPath
	CreateReportPath
	CryptoWithdrawalPath
	CurrenciesPath
	CurrencyConversionPath
	CurrencyPath
	DeleteProfilePath
	ExchangeLimitsPath
	FeesPath
	FillsPath
	GenerateCryptoAddressPath
	OrderPath
	OrdersPath
	PaymentMethodDepositPath
	PaymentMethodWithdrawalPath
	PaymentMethodsPath
	ProductPath
	ProductStatsPath
	ProductTickerPath
	ProductsPath
	ProfilePath
	ProfilesPath
	RenameProfilePath
	ReportPath
	ReportsPath
	SignedPricesPath
	TradesPath
	TransferPath
	TransfersPath
	WalletsPath
	WithdrawalFeeEstimatePath
)

// Account will return data for a single account. Use this endpoint when you know the account_id. API key must belong to
// the same profile as the account.
func getAccountPath(params map[string]string) string {
	return path.Join("/accounts", params["account_id"])
}

// AccountHolds will return the holds of an account that belong to the same profile as the API key. Holds are placed on
// an account for any active orders or pending withdraw requests. As an order is filled, the hold amount is updated. If
// an order is canceled, any remaining hold is removed. For withdrawals, the hold is removed after it is completed.
func getAccountHoldsPath(params map[string]string) string {
	return path.Join("/accounts", params["account_id"], "holds")
}

// AccountLedger returns ledger activity for an account. This includes anything that would affect the accounts balance -
// transfers, trades, fees, etc. This endpoint requires either the "view" or "trade" permission.
func getAccountLedgerPath(params map[string]string) string {
	return path.Join("/accounts", params["account_id"], "ledger")
}

// AccountTransfers returns past withdrawals and deposits for an account.
func getAccountTransfersPath(params map[string]string) string {
	return path.Join("/accounts", params["account_id"], "transfers")
}

// Accounts will return a list of trading accounts from the profile of the API key.
func getAccountsPath(params map[string]string) string {
	return path.Join("/accounts")
}

// Book will return a list of open orders for a product. The amount of detail shown can be customized with the level
// parameter.
func getBookPath(params map[string]string) string {
	return path.Join("/products", params["product_id"], "book")
}

// CancelOpenOrders will try with best effort to cancel all open orders. This may require you to make the request
// multiple times until all of the open orders are deleted.
func getCancelOpenOrdersPath(params map[string]string) string {
	return path.Join("/orders")
}

// CancelOrder will cancel a single open order by order id.
func getCancelOrderPath(params map[string]string) string {
	return path.Join("/orders", params["order_id"])
}

// Candles will return historic rates for a product.
func getCandlesPath(params map[string]string) string {
	return path.Join("/products", params["product_id"], "candles")
}

// CoinbaseAccountDeposit funds from a www.coinbase.com wallet to the specified profile_id.
func getCoinbaseAccountDepositPath(params map[string]string) string {
	return path.Join("/deposits", "coinbase-account")
}

// AccountWithdraws funds from the specified profile_id to a www.coinbase.com wallet. Withdraw funds to a coinbase
// account. You can move funds between your Coinbase accounts and your Coinbase Exchange trading accounts within your
// daily limits. Moving funds between Coinbase and Coinbase Exchange is instant and free. See the Coinbase Accounts
// section for retrieving your Coinbase accounts. This endpoint requires the "transfer" permission.
func getCoinbaseAccountWithdrawalPath(params map[string]string) string {
	return path.Join("/withdrawals", "coinbase-account")
}

// ConvertCurrency converts funds from from currency to to currency. Funds are converted on the from account in the
// profile_id profile. This endpoint requires the "trade" permission. A successful conversion will be assigned a
// conversion id. The corresponding ledger entries for a conversion will reference this conversion id
func getConvertCurrencyPath(params map[string]string) string {
	return path.Join("/conversions")
}

// CreateOrder will create a new an order. You can place two types of orders: limit and market. Orders can only be
// placed if your account has sufficient funds. Once an order is placed, your account funds will be put on hold for the
// duration of the order. How much and which funds are put on hold depends on the order type and parameters specified.
func getCreateOrderPath(params map[string]string) string {
	return path.Join("/orders")
}

// CreateProfile will create a new profile. Will fail if no name is provided or if user already has max number of
// profiles.
func getCreateProfilePath(params map[string]string) string {
	return path.Join("/profiles")
}

// CreateProfileTransfer will transfer an amount of currency from one profile to another. This endpoint requires the
// "transfer" permission.
func getCreateProfileTransferPath(params map[string]string) string {
	return path.Join("/profiles", "transfer")
}

// CreateReport generates a report. Reports are either for past account history or past fills on either all accounts or
// one product's account.
func getCreateReportPath(params map[string]string) string {
	return path.Join("/reports")
}

// CryptoWithdrawal funds from the specified profile_id to an external crypto address. This endpoint requires the
// "transfer" permission. API key must belong to default profile.
func getCryptoWithdrawalPath(params map[string]string) string {
	return path.Join("/withdrawals", "crypto")
}

// Currencies returns a list of all known currencies. Note: Not all currencies may be currently in use for trading.
func getCurrenciesPath(params map[string]string) string {
	return path.Join("/currencies")
}

// Currency returns a single currency by id.
func getCurrencyPath(params map[string]string) string {
	return path.Join("/currencies", params["currency_id"])
}

// CurrencyConversion returns the currency conversion by conversion id (i.e. USD -> USDC).
func getCurrencyConversionPath(params map[string]string) string {
	return path.Join("/conversions", params["conversion_id"])
}

// DeleteProfile deletes the profile specified by profile_id and transfers all funds to the profile specified by to.
// Fails if there are any open orders on the profile to be deleted.
func getDeleteProfilePath(params map[string]string) string {
	return path.Join("/profiles", params["profile_id"], "deactivate")
}

// ExchangeLimits returns exchange limits information for a single user.
func getExchangeLimitsPath(params map[string]string) string {
	return path.Join("/users", params["user_id"], "exchange-limits")
}

// Fees returns fees rates and 30 days trailing volume.
func getFeesPath(params map[string]string) string {
	return path.Join("/fees")
}

// Fills returns a list of fills. A fill is a partial or complete match on a specific order.
func getFillsPath(params map[string]string) string {
	return path.Join("/fills")
}

// GenerateCryptoAddress will create a one-time crypto address for depositing crypto, using a wallet account id. This
// endpoint requires the "transfer" permission. API key must belong to default profile.
func getGenerateCryptoAddressPath(params map[string]string) string {
	return path.Join("/coinbase-accounts", params["account_id"], "addresses")
}

// Order returns a single order by id.
func getOrderPath(params map[string]string) string {
	return path.Join("/orders", params["order_id"])
}

// Orders will return your current open orders. Only open or un-settled orders are returned by default. As soon as an
// order is no longer open and settled, it will no longer appear in the default request. Open orders may change state
// between the request and the response depending on market conditions.
func getOrdersPath(params map[string]string) string {
	return path.Join("/orders")
}

// PaymentMethodDeposit will fund from a linked external payment method to the specified profile_id.
func getPaymentMethodDepositPath(params map[string]string) string {
	return path.Join("/deposits", "payment-method")
}

// PaymentMethodWithdrawal will fund from the specified profile_id to a linked external payment method. This endpoint
// requires the "transfer" permission. API key is restricted to the default profile.
func getPaymentMethodWithdrawalPath(params map[string]string) string {
	return path.Join("/withdrawals", "payment-method")
}

// PaymentMethods returns a list of the user's linked payment methods.
func getPaymentMethodsPath(params map[string]string) string {
	return path.Join("/payment-methods")
}

// Product will return information on a single product.
func getProductPath(params map[string]string) string {
	return path.Join("/products", params["product_id"])
}

// ProductStats will return 30day and 24hour stats for a product.
func getProductStatsPath(params map[string]string) string {
	return path.Join("/products", params["product_id"], "stats")
}

// ProductTicker will return snapshot information about the last trade (tick), best bid/ask and 24h volume.
func getProductTickerPath(params map[string]string) string {
	return path.Join("/products", params["product_id"], "ticker")
}

// Products will return a list of available currency pairs for trading.
func getProductsPath(params map[string]string) string {
	return path.Join("/products")
}

// Profile returns information for a single profile. Use this endpoint when you know the profile_id. This endpoint
// requires the "view" permission and is accessible by any profile's API key.
func getProfilePath(params map[string]string) string {
	return path.Join("/profiles", params["profile_id"])
}

// Profiles returns a list of all of the current user's profiles.
func getProfilesPath(params map[string]string) string {
	return path.Join("/profiles")
}

// RenameProfile will rename a profile. Names 'default' and 'margin' are reserved.
func getRenameProfilePath(params map[string]string) string {
	return path.Join("/profiles", params["profile_id"])
}

// Report will return a specific report by report_id.
func getReportPath(params map[string]string) string {
	return path.Join("/reports", params["report_id"])
}

// Reports returns a list of past fills/account reports.
func getReportsPath(params map[string]string) string {
	return path.Join("/reports")
}

// SignedPrices returns cryptographically signed prices ready to be posted on-chain using Compound's Open Oracle smart
// contract.
func getSignedPricesPath(params map[string]string) string {
	return path.Join("/oracle")
}

// Trades retruns a list the latest trades for a product.
func getTradesPath(params map[string]string) string {
	return path.Join("/products", params["product_id"], "trades")
}

// AccountTransfer returns information on a single transfer.
func getTransferPath(params map[string]string) string {
	return path.Join("/transfers", params["transfer_id"])
}

// Transfers is a list of in-progress and completed transfers of funds in/out of any of the user's accounts.
func getTransfersPath(params map[string]string) string {
	return path.Join("/transfers")
}

// Wallets will return all the user's available Coinbase wallets (These are the wallets/accounts that are used for
// buying and selling on www.coinbase.com)
func getWalletsPath(params map[string]string) string {
	return path.Join("/coinbase-accounts")
}

// WithdrawalFeeEstimate will return the fee estimate for the crypto withdrawal to crypto address
func getWithdrawalFeeEstimatePath(params map[string]string) string {
	return path.Join("/withdrawals", "fee-estimate")
}

// Get takes an rawPath const and rawPath arguments to parse the URL rawPath path.
func (p rawPath) Path(params map[string]string) string {
	return map[rawPath]func(map[string]string) string{
		AccountPath:                   getAccountPath,
		AccountHoldsPath:              getAccountHoldsPath,
		AccountLedgerPath:             getAccountLedgerPath,
		AccountTransfersPath:          getAccountTransfersPath,
		AccountsPath:                  getAccountsPath,
		BookPath:                      getBookPath,
		CancelOpenOrdersPath:          getCancelOpenOrdersPath,
		CancelOrderPath:               getCancelOrderPath,
		CandlesPath:                   getCandlesPath,
		CoinbaseAccountDepositPath:    getCoinbaseAccountDepositPath,
		CoinbaseAccountWithdrawalPath: getCoinbaseAccountWithdrawalPath,
		ConvertCurrencyPath:           getConvertCurrencyPath,
		CreateOrderPath:               getCreateOrderPath,
		CreateProfilePath:             getCreateProfilePath,
		CreateProfileTransferPath:     getCreateProfileTransferPath,
		CreateReportPath:              getCreateReportPath,
		CryptoWithdrawalPath:          getCryptoWithdrawalPath,
		CurrenciesPath:                getCurrenciesPath,
		CurrencyPath:                  getCurrencyPath,
		CurrencyConversionPath:        getCurrencyConversionPath,
		DeleteProfilePath:             getDeleteProfilePath,
		ExchangeLimitsPath:            getExchangeLimitsPath,
		FeesPath:                      getFeesPath,
		FillsPath:                     getFillsPath,
		GenerateCryptoAddressPath:     getGenerateCryptoAddressPath,
		OrderPath:                     getOrderPath,
		OrdersPath:                    getOrdersPath,
		PaymentMethodDepositPath:      getPaymentMethodDepositPath,
		PaymentMethodWithdrawalPath:   getPaymentMethodWithdrawalPath,
		PaymentMethodsPath:            getPaymentMethodsPath,
		ProductPath:                   getProductPath,
		ProductStatsPath:              getProductStatsPath,
		ProductTickerPath:             getProductTickerPath,
		ProductsPath:                  getProductsPath,
		ProfilePath:                   getProfilePath,
		ProfilesPath:                  getProfilesPath,
		RenameProfilePath:             getRenameProfilePath,
		ReportPath:                    getReportPath,
		ReportsPath:                   getReportsPath,
		SignedPricesPath:              getSignedPricesPath,
		TradesPath:                    getTradesPath,
		TransferPath:                  getTransferPath,
		TransfersPath:                 getTransfersPath,
		WalletsPath:                   getWalletsPath,
		WithdrawalFeeEstimatePath:     getWithdrawalFeeEstimatePath,
	}[p](params)
}

func (p rawPath) Scope() string {
	return map[rawPath]string{}[p]
}
