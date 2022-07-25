package coinbasepro

import "github.com/alpine-hodler/web/internal"

// * This is a generated file, do not edit

// Account will return data for a single account. Use this endpoint when you know the account_id. API key must belong to
// the same profile as the account.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccount
func (c *Client) Account(accountId string) (m *Account, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(AccountPath),
		internal.HTTPWithParams(map[string]string{
			"account_id": accountId}),
		internal.HTTPWithRatelimiter(getRateLimiter(AccountRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// AccountHolds will return the holds of an account that belong to the same profile as the API key. Holds are placed on
// an account for any active orders or pending withdraw requests. As an order is filled, the hold amount is updated. If
// an order is canceled, any remaining hold is removed. For withdrawals, the hold is removed after it is completed.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccountholds
func (c *Client) AccountHolds(accountId string, opts *AccountHoldsOptions) (m []*AccountHold, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(AccountHoldsPath),
		internal.HTTPWithParams(map[string]string{
			"account_id": accountId}),
		internal.HTTPWithRatelimiter(getRateLimiter(AccountHoldsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// AccountLedger returns ledger activity for an account. This includes anything that would affect the accounts balance -
// transfers, trades, fees, etc. This endpoint requires either the "view" or "trade" permission.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccountledger
func (c *Client) AccountLedger(accountId string, opts *AccountLedgerOptions) (m []*AccountLedger, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(AccountLedgerPath),
		internal.HTTPWithParams(map[string]string{
			"account_id": accountId}),
		internal.HTTPWithRatelimiter(getRateLimiter(AccountLedgerRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// AccountTransfers returns past withdrawals and deposits for an account.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccounttransfers
func (c *Client) AccountTransfers(accountId string, opts *AccountTransfersOptions) (m []*Transfer, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(AccountTransfersPath),
		internal.HTTPWithParams(map[string]string{
			"account_id": accountId}),
		internal.HTTPWithRatelimiter(getRateLimiter(AccountTransfersRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Accounts will return a list of trading accounts from the profile of the API key.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getaccounts
func (c *Client) Accounts() (m []*Account, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(AccountsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(AccountsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Book will return a list of open orders for a product. The amount of detail shown can be customized with the level
// parameter.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproductbook
func (c *Client) Book(productId string, opts *BookOptions) (m *Book, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(BookPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(BookRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CancelOpenOrders will try with best effort to cancel all open orders. This may require you to make the request
// multiple times until all of the open orders are deleted.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_deleteorders
func (c *Client) CancelOpenOrders(opts *CancelOpenOrdersOptions) (m []*string, _ error) {
	req, _ := internal.HTTPNewRequest("DELETE", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(CancelOpenOrdersPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CancelOpenOrdersRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CancelOrder will cancel a single open order by order id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_deleteorder
func (c *Client) CancelOrder(orderId string) (m string, _ error) {
	req, _ := internal.HTTPNewRequest("DELETE", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CancelOrderPath),
		internal.HTTPWithParams(map[string]string{
			"order_id": orderId}),
		internal.HTTPWithRatelimiter(getRateLimiter(CancelOrderRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Candles will return historic rates for a product.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproductcandles
func (c *Client) Candles(productId string, opts *CandlesOptions) (m *Candles, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(CandlesPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(CandlesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CoinbaseAccountDeposit funds from a www.coinbase.com wallet to the specified profile_id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postdepositcoinbaseaccount
func (c *Client) CoinbaseAccountDeposit(opts *CoinbaseAccountDepositOptions) (m *Deposit, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CoinbaseAccountDepositPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CoinbaseAccountDepositRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// AccountWithdraws funds from the specified profile_id to a www.coinbase.com wallet. Withdraw funds to a coinbase
// account. You can move funds between your Coinbase accounts and your Coinbase Exchange trading accounts within your
// daily limits. Moving funds between Coinbase and Coinbase Exchange is instant and free. See the Coinbase Accounts
// section for retrieving your Coinbase accounts. This endpoint requires the "transfer" permission.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postwithdrawcoinbaseaccount
func (c *Client) CoinbaseAccountWithdrawal(opts *CoinbaseAccountWithdrawalOptions) (m *Withdrawal, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CoinbaseAccountWithdrawalPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CoinbaseAccountWithdrawalRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// ConvertCurrency converts funds from from currency to to currency. Funds are converted on the from account in the
// profile_id profile. This endpoint requires the "trade" permission. A successful conversion will be assigned a
// conversion id. The corresponding ledger entries for a conversion will reference this conversion id
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postconversion
func (c *Client) ConvertCurrency(opts *ConvertCurrencyOptions) (m *CurrencyConversion, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ConvertCurrencyPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(ConvertCurrencyRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CreateOrder will create a new an order. You can place two types of orders: limit and market. Orders can only be
// placed if your account has sufficient funds. Once an order is placed, your account funds will be put on hold for the
// duration of the order. How much and which funds are put on hold depends on the order type and parameters specified.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postorders
func (c *Client) CreateOrder(opts *CreateOrderOptions) (m *CreateOrder, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CreateOrderPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CreateOrderRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CreateProfile will create a new profile. Will fail if no name is provided or if user already has max number of
// profiles.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postprofile
func (c *Client) CreateProfile(opts *CreateProfileOptions) (m *Profile, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CreateProfilePath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CreateProfileRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CreateProfileTransfer will transfer an amount of currency from one profile to another. This endpoint requires the
// "transfer" permission.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postprofiletransfer
func (c *Client) CreateProfileTransfer(opts *CreateProfileTransferOptions) error {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return internal.HTTPFetch(nil, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CreateProfileTransferPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CreateProfileTransferRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CreateReport generates a report. Reports are either for past account history or past fills on either all accounts or
// one product's account.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postreports
func (c *Client) CreateReport(opts *CreateReportOptions) (m *CreateReport, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CreateReportPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CreateReportRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CryptoWithdrawal funds from the specified profile_id to an external crypto address. This endpoint requires the
// "transfer" permission. API key must belong to default profile.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postwithdrawcrypto
func (c *Client) CryptoWithdrawal(opts *CryptoWithdrawalOptions) (m *Withdrawal, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CryptoWithdrawalPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CryptoWithdrawalRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Currencies returns a list of all known currencies. Note: Not all currencies may be currently in use for trading.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getcurrencies
func (c *Client) Currencies() (m []*Currency, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CurrenciesPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(CurrenciesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Currency returns a single currency by id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getcurrency
func (c *Client) Currency(currencyId string) (m *Currency, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(CurrencyPath),
		internal.HTTPWithParams(map[string]string{
			"currency_id": currencyId}),
		internal.HTTPWithRatelimiter(getRateLimiter(CurrencyRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// CurrencyConversion returns the currency conversion by conversion id (i.e. USD -> USDC).
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getconversion
func (c *Client) CurrencyConversion(conversionId string, opts *CurrencyConversionOptions) (m *CurrencyConversion, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(CurrencyConversionPath),
		internal.HTTPWithParams(map[string]string{
			"conversion_id": conversionId}),
		internal.HTTPWithRatelimiter(getRateLimiter(CurrencyConversionRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// DeleteProfile deletes the profile specified by profile_id and transfers all funds to the profile specified by to.
// Fails if there are any open orders on the profile to be deleted.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_putprofiledeactivate
func (c *Client) DeleteProfile(profileId string, opts *DeleteProfileOptions) error {
	req, _ := internal.HTTPNewRequest("PUT", "", opts)
	return internal.HTTPFetch(nil, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(DeleteProfilePath),
		internal.HTTPWithParams(map[string]string{
			"profile_id": profileId}),
		internal.HTTPWithRatelimiter(getRateLimiter(DeleteProfileRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// ExchangeLimits returns exchange limits information for a single user.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getuserexchangelimits
func (c *Client) ExchangeLimits(userId string) (m *ExchangeLimits, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ExchangeLimitsPath),
		internal.HTTPWithParams(map[string]string{
			"user_id": userId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ExchangeLimitsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Fees returns fees rates and 30 days trailing volume.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getfees
func (c *Client) Fees() (m *Fees, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(FeesPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(FeesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Fills returns a list of fills. A fill is a partial or complete match on a specific order.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getfills
func (c *Client) Fills(opts *FillsOptions) (m []*Fill, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(FillsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(FillsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// GenerateCryptoAddress will create a one-time crypto address for depositing crypto, using a wallet account id. This
// endpoint requires the "transfer" permission. API key must belong to default profile.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postcoinbaseaccountaddresses
func (c *Client) GenerateCryptoAddress(accountId string) (m *CryptoAddress, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(GenerateCryptoAddressPath),
		internal.HTTPWithParams(map[string]string{
			"account_id": accountId}),
		internal.HTTPWithRatelimiter(getRateLimiter(GenerateCryptoAddressRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Order returns a single order by id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getorder
func (c *Client) Order(orderId string) (m *Order, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(OrderPath),
		internal.HTTPWithParams(map[string]string{
			"order_id": orderId}),
		internal.HTTPWithRatelimiter(getRateLimiter(OrderRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Orders will return your current open orders. Only open or un-settled orders are returned by default. As soon as an
// order is no longer open and settled, it will no longer appear in the default request. Open orders may change state
// between the request and the response depending on market conditions.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getorders
func (c *Client) Orders(opts *OrdersOptions) (m []*Order, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(OrdersPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(OrdersRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// PaymentMethodDeposit will fund from a linked external payment method to the specified profile_id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postdepositpaymentmethod
func (c *Client) PaymentMethodDeposit(opts *PaymentMethodDepositOptions) (m *Deposit, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(PaymentMethodDepositPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(PaymentMethodDepositRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// PaymentMethodWithdrawal will fund from the specified profile_id to a linked external payment method. This endpoint
// requires the "transfer" permission. API key is restricted to the default profile.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_postwithdrawpaymentmethod
func (c *Client) PaymentMethodWithdrawal(opts *PaymentMethodWithdrawalOptions) (m *Withdrawal, _ error) {
	req, _ := internal.HTTPNewRequest("POST", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(PaymentMethodWithdrawalPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(PaymentMethodWithdrawalRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// PaymentMethods returns a list of the user's linked payment methods.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getpaymentmethods
func (c *Client) PaymentMethods() (m []*PaymentMethod, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(PaymentMethodsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(PaymentMethodsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Product will return information on a single product.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproduct
func (c *Client) Product(productId string) (m *Product, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ProductPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ProductRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// ProductStats will return 30day and 24hour stats for a product.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproductstats
func (c *Client) ProductStats(productId string) (m *ProductStats, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ProductStatsPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ProductStatsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// ProductTicker will return snapshot information about the last trade (tick), best bid/ask and 24h volume.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproductticker
func (c *Client) ProductTicker(productId string) (m *ProductTicker, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ProductTickerPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ProductTickerRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Products will return a list of available currency pairs for trading.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproducts
func (c *Client) Products(opts *ProductsOptions) (m []*Product, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(ProductsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(ProductsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Profile returns information for a single profile. Use this endpoint when you know the profile_id. This endpoint
// requires the "view" permission and is accessible by any profile's API key.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getprofile
func (c *Client) Profile(profileId string, opts *ProfileOptions) (m *Profile, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(ProfilePath),
		internal.HTTPWithParams(map[string]string{
			"profile_id": profileId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ProfileRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Profiles returns a list of all of the current user's profiles.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getprofiles
func (c *Client) Profiles(opts *ProfilesOptions) (m []*Profile, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(ProfilesPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(ProfilesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// RenameProfile will rename a profile. Names 'default' and 'margin' are reserved.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_putprofile
func (c *Client) RenameProfile(profileId string, opts *RenameProfileOptions) (m *Profile, _ error) {
	req, _ := internal.HTTPNewRequest("PUT", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(RenameProfilePath),
		internal.HTTPWithParams(map[string]string{
			"profile_id": profileId}),
		internal.HTTPWithRatelimiter(getRateLimiter(RenameProfileRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Report will return a specific report by report_id.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getreport
func (c *Client) Report(reportId string) (m *Report, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(ReportPath),
		internal.HTTPWithParams(map[string]string{
			"report_id": reportId}),
		internal.HTTPWithRatelimiter(getRateLimiter(ReportRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Reports returns a list of past fills/account reports.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getreports
func (c *Client) Reports(opts *ReportsOptions) (m []*Report, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(ReportsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(ReportsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// SignedPrices returns cryptographically signed prices ready to be posted on-chain using Compound's Open Oracle smart
// contract.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getcoinbasepriceoracle
func (c *Client) SignedPrices() (m *Oracle, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(SignedPricesPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(SignedPricesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Trades retruns a list the latest trades for a product.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getproducttrades
func (c *Client) Trades(productId string, opts *TradesOptions) (m []*Trade, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(TradesPath),
		internal.HTTPWithParams(map[string]string{
			"product_id": productId}),
		internal.HTTPWithRatelimiter(getRateLimiter(TradesRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// AccountTransfer returns information on a single transfer.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_gettransfer
func (c *Client) Transfer(transferId string) (m *Transfer, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(TransferPath),
		internal.HTTPWithParams(map[string]string{
			"transfer_id": transferId}),
		internal.HTTPWithRatelimiter(getRateLimiter(TransferRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Transfers is a list of in-progress and completed transfers of funds in/out of any of the user's accounts.
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_gettransfers
func (c *Client) Transfers() (m []*Transfer, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(TransfersPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(TransfersRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// Wallets will return all the user's available Coinbase wallets (These are the wallets/accounts that are used for
// buying and selling on www.coinbase.com)
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getcoinbaseaccounts
func (c *Client) Wallets() (m []*Wallet, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", nil)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(nil),
		internal.HTTPWithEndpoint(WalletsPath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(WalletsRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}

// WithdrawalFeeEstimate will return the fee estimate for the crypto withdrawal to crypto address
//
// source: https://docs.cloud.coinbase.com/exchange/reference/exchangerestapi_getwithdrawfeeestimate
func (c *Client) WithdrawalFeeEstimate(opts *WithdrawalFeeEstimateOptions) (m *WithdrawalFeeEstimate, _ error) {
	req, _ := internal.HTTPNewRequest("GET", "", opts)
	return m, internal.HTTPFetch(&m, internal.HTTPWithClient(c.Client),
		internal.HTTPWithEncoder(opts),
		internal.HTTPWithEndpoint(WithdrawalFeeEstimatePath),
		internal.HTTPWithParams(nil),
		internal.HTTPWithRatelimiter(getRateLimiter(WithdrawalFeeEstimateRatelimiter, 5)),
		internal.HTTPWithRequest(req))
}
