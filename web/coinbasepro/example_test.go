package coinbasepro_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alpine-hodler/web/pkg/coinbasepro"
	"github.com/alpine-hodler/web/pkg/transport"
	"github.com/alpine-hodler/web/pkg/websocket"
	"github.com/alpine-hodler/web/tools"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// These are configurable data for testing on.  Particularly useful if coinbase decides to do a "Cancel Only"
// restriction to a currency or product.
const (
	exampleCurrency  = "ETH"
	exampleProductID = "ETH-BTC"
)

const (
	testAccountID              = "CB_PRO_ACCOUNT_ID"
	testCurrency               = "CB_PRO_ACCOUNT_CURRENCY"
	testOrderIDForCancellation = "CB_PRO_ORDER_ID_FOR_CANCELLATION"
	testProfileID              = "CB_PRO_PROFILE_ID"
	testProductID              = "CB_PRO_PRODUCT_ID"
	testUSDCWalletID           = "CB_PRO_USDC_WALLET_ID"
	testUserID                 = "CB_PRO_USER_ID"
	testWalletID               = "CB_PRO_WALLET_ID"
)

func TestExamples(t *testing.T) {
	defer tools.Quiet()()

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ! Make sure that these tests only run on the sandbox URL
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	godotenv.Load(".simple-test.env")
	os.Setenv("CB_PRO_URL", "https://api-public.sandbox.exchange.coinbase.com") // safety check

	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Cancel all orders before creating the new ones.
	cancelOpenOrders(t, client)

	// Set environment variables for integration test.
	os.Setenv(testProductID, exampleProductID)
	setAccountEnvVariables(t, client, exampleCurrency)
	setWalletEnvVariables(t, client, os.Getenv(testCurrency))
	createDepositTransfer(t, client, os.Getenv(testAccountID), os.Getenv(testCurrency), os.Getenv(testWalletID))
	createUnlikelyLimitOrder(t, client, os.Getenv(testProfileID), os.Getenv(testProductID))
	createUnlikelyLimitOrderForCancellation(t, client, os.Getenv(testProfileID), os.Getenv(testProductID))

	// Run Examples
	t.Run("Client.Account", func(t *testing.T) { ExampleClient_Account() })
	t.Run("Client.AccountHolds", func(t *testing.T) { ExampleClient_AccountHolds() })
	t.Run("Client.AccountLedger", func(t *testing.T) { ExampleClient_AccountLedger() })
	t.Run("Client.AccountTransfers", func(t *testing.T) { ExampleClient_AccountTransfers() })
	t.Run("Client.Accounts", func(t *testing.T) { ExampleClient_Accounts() })
	t.Run("Client.Book", func(t *testing.T) { ExampleClient_Book() })
	t.Run("Client.CancelOpenOrders", func(t *testing.T) { ExampleClient_CancelOpenOrders() })
	t.Run("Client.CancelOrder", func(t *testing.T) { ExampleClient_CancelOrder() })
	t.Run("Client.Candles", func(t *testing.T) { ExampleClient_Candles() })
	t.Run("Client.CoinbaseAccountDeposit", func(t *testing.T) { ExampleClient_CoinbaseAccountDeposit() })
	t.Run("Client.CoinbaseAccountWithdrawal", func(t *testing.T) { ExampleClient_CoinbaseAccountWithdrawal() })
	t.Run("Client.ConvertCurrency", func(t *testing.T) { ExampleClient_ConvertCurrency() })
	t.Run("Client.CreateOrder", func(t *testing.T) { ExampleClient_CreateOrder() })
	t.Run("Client.CreateProfile", func(t *testing.T) { ExampleClient_CreateProfile() })
	t.Run("Client.CreateProfileTransfer", func(t *testing.T) { ExampleClient_CreateProfileTransfer() })
	t.Run("Client.CreateReport", func(t *testing.T) { ExampleClient_CreateReport() })
	t.Run("Client.CryptoWithdrawal", func(t *testing.T) { ExampleClient_CryptoWithdrawal() })
	t.Run("Client.Currencies", func(t *testing.T) { ExampleClient_Currencies() })
	t.Run("Client.Currency", func(t *testing.T) { ExampleClient_Currency() })
	t.Run("Client.DeleteProfile", func(t *testing.T) { ExampleClient_DeleteProfile() })
	t.Run("Client.ExchangeLimits", func(t *testing.T) { ExampleClient_ExchangeLimits() })
	t.Run("Client.Fees", func(t *testing.T) { ExampleClient_Fees() })
	t.Run("Client.Fills", func(t *testing.T) { ExampleClient_Fills() })
	t.Run("Client.GenerateCryptoAddress", func(t *testing.T) { ExampleClient_GenerateCryptoAddress() })
	t.Run("Client.Order", func(t *testing.T) { ExampleClient_Order() })
	t.Run("Client.Orders", func(t *testing.T) { ExampleClient_Orders() })
	t.Run("Client.PaymentDepositMethod", func(t *testing.T) { ExampleClient_PaymentMethodDeposit() })
	t.Run("Client.PaymentWithdrawalMethod", func(t *testing.T) { ExampleClient_PaymentMethodWithdrawal() })
	t.Run("Client.PaymentMethods", func(t *testing.T) { ExampleClient_PaymentMethods() })
	t.Run("Client.Product", func(t *testing.T) { ExampleClient_Product() })
	t.Run("Client.ProductStats", func(t *testing.T) { ExampleClient_ProductStats() })
	t.Run("Client.ProductTicker", func(t *testing.T) { ExampleClient_ProductTicker() })
	t.Run("Client.Products", func(t *testing.T) { ExampleClient_Products() })
	t.Run("Client.Profile", func(t *testing.T) { ExampleClient_Profile() })
	t.Run("Client.Profiles", func(t *testing.T) { ExampleClient_Profiles() })
	t.Run("Client.RenameProfile", func(t *testing.T) { ExampleClient_RenameProfile() })
	t.Run("Client.Report", func(t *testing.T) { ExampleClient_Report() })
	t.Run("Client.Reports", func(t *testing.T) { ExampleClient_Reports() })
	t.Run("Client.SignedPrices", func(t *testing.T) { ExampleClient_SignedPrices() })
	t.Run("Client.Trades", func(t *testing.T) { ExampleClient_Trades() })
	t.Run("Client.Transfer", func(t *testing.T) { ExampleClient_Transfer() })
	t.Run("Client.Transfers", func(t *testing.T) { ExampleClient_Transfers() })
	t.Run("Client.Wallets", func(t *testing.T) { ExampleClient_Wallets() })
	t.Run("Client.WithdrawalFeeEstimate", func(t *testing.T) { ExampleClient_WithdrawalFeeEstimate() })
	t.Run("ProductWebsocket.Ticker", func(t *testing.T) { ExampleProductWebsocket_Ticker() })
}

func ExampleNewClient() {
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("A Coinbase Pro client: %v", client)
}

func ExampleClient_Account() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get the account ID to look up.
	accountID := os.Getenv("CB_PRO_ACCOUNT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	account, err := client.Account(accountID)
	if err != nil {
		log.Fatalf("Error fetching account: %v", err)
	}
	fmt.Printf("account: %+v\n", account)
}

func ExampleClient_AccountHolds() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get the account ID to look up.
	accountID := os.Getenv("CB_PRO_ACCOUNT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Since the order above will forever be in an "active" status, we should
	// have a hold value waiting for us.
	holds, err := client.AccountHolds(accountID, new(coinbasepro.AccountHoldsOptions).
		SetLimit(1).
		SetBefore("2010-01-01").
		SetAfter("2080-01-01"))
	if err != nil {
		log.Fatalf("Error fetching holds: %v", err)
	}
	fmt.Printf("holds: %+v\n", holds)
}

func ExampleClient_AccountLedger() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get the account ID and profile ID for request.
	accountID := os.Getenv("CB_PRO_ACCOUNT_ID")
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ledger, err := client.AccountLedger(accountID, new(coinbasepro.AccountLedgerOptions).
		SetStartDate("2010-01-01").
		SetEndDate("2080-01-01").
		SetAfter(1526365354).
		SetLimit(1).
		SetProfileID(profileID))
	if err != nil {
		log.Fatalf("Error fetching ledger: %v", err)
	}
	fmt.Printf("ledger: %+v\n", ledger)
}

func ExampleClient_AccountTransfers() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get the account ID for request.
	accountID := os.Getenv("CB_PRO_ACCOUNT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Get the resulting Account Transfers.
	transfers, err := client.AccountTransfers(accountID, new(coinbasepro.AccountTransfersOptions).
		SetBefore("2010-01-01").
		SetAfter("2080-01-01").
		SetType(coinbasepro.TransferMethodDeposit).
		SetLimit(1))
	if err != nil {
		log.Fatalf("Error fetching account transfers: %v", err)
	}
	fmt.Printf("account transfers: %+v\n", transfers)
}

func ExampleClient_Accounts() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	accounts, err := client.Accounts()
	if err != nil {
		log.Fatalf("Error fetching accounts: %v", err)
	}
	fmt.Printf("accounts: %+v\n", accounts)
}

func ExampleClient_Book() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	book, err := client.Book("BTC-USD", new(coinbasepro.BookOptions).SetLevel(1))
	if err != nil {
		log.Fatalf("Error fetching book: %v", err)
	}
	fmt.Printf("book: %+v\n", book)
}

func ExampleClient_CancelOpenOrders() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get the profile ID for request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Cancel the order by simply cancelling all BTC-USD orders for some profile ID.
	_, err = client.CancelOpenOrders(new(coinbasepro.CancelOpenOrdersOptions).
		SetProductID("BTC-USD").
		SetProfileID(profileID))
	if err != nil {
		log.Fatalf("Error canceling open orders: %v", err)
	}
}

func ExampleClient_CancelOrder() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get an orderID for cancellation
	orderID := os.Getenv("CB_PRO_ORDER_ID_FOR_CANCELLATION")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Cancel the order by using the order's ID.
	_, err = client.CancelOrder(orderID)
	if err != nil {
		log.Fatalf("Error canceling order: %v", err)
	}
}

func ExampleClient_Candles() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	startTimestamp := time.Date(2018, 11, 9, 17, 11, 8, 0, time.UTC).Format(time.RFC3339)
	endTimestamp := time.Date(2018, 11, 9, 18, 11, 8, 0, time.UTC).Format(time.RFC3339)

	candles, err := client.Candles("BTC-USD", new(coinbasepro.CandlesOptions).
		SetStart(startTimestamp).
		SetEnd(endTimestamp).
		SetGranularity(coinbasepro.Granularity60))
	if err != nil {
		log.Fatalf("Error canceling order: %v", err)
	}
	fmt.Printf("candle example: %+v\n", (*candles)[0])
}

func ExampleClient_CoinbaseAccountDeposit() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID to deposit into and a wallet ID to withdraw from.
	accountCurrency := os.Getenv("CB_PRO_ACCOUNT_CURRENCY")
	profileID := os.Getenv("CB_PRO_PROFILE_ID")
	walletID := os.Getenv("CB_PRO_WALLET_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Make the deposit.
	_, err = client.CoinbaseAccountDeposit(new(coinbasepro.CoinbaseAccountDepositOptions).
		SetProfileID(profileID).
		SetCoinbaseAccountID(walletID).
		SetAmount(1).
		SetCurrency(accountCurrency))
	if err != nil {
		log.Fatalf("Error making deposit: %v", err)
	}
}

func ExampleClient_CoinbaseAccountWithdrawal() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID to withdraw from and a wallet ID to deposit int.
	accountCurrency := os.Getenv("CB_PRO_ACCOUNT_CURRENCY")
	profileID := os.Getenv("CB_PRO_PROFILE_ID")
	walletID := os.Getenv("CB_PRO_WALLET_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Withdraw the deposit.
	_, err = client.CoinbaseAccountWithdrawal(new(coinbasepro.CoinbaseAccountWithdrawalOptions).
		SetProfileID(profileID).
		SetCoinbaseAccountID(walletID).
		SetAmount(1).
		SetCurrency(accountCurrency))
	if err != nil {
		log.Fatalf("Error making withdrawal: %v", err)
	}
}

func ExampleClient_ConvertCurrency() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Convert the USD into USDC.
	_, err = client.ConvertCurrency(new(coinbasepro.ConvertCurrencyOptions).
		SetAmount(1).
		SetFrom("USD").
		SetTo("USDC").
		SetProfileID(profileID))
	if err != nil {
		log.Fatalf("Error converting currency: %v", err)
	}
}

func ExampleClient_CreateOrder() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")
	productID := os.Getenv("CB_PRO_PRODUCT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	order, err := client.CreateOrder(new(coinbasepro.CreateOrderOptions).
		SetProfileID(profileID).
		SetType(coinbasepro.OrderTypeLimit).
		SetSide(coinbasepro.SideSell).
		SetSTP(coinbasepro.STPDc).
		SetStop(coinbasepro.StopLoss).
		SetTimeInForce(coinbasepro.TimeInForceGTC).
		SetCancelAfter(coinbasepro.CancelAfterMin).
		SetProductID(productID).
		SetStopPrice(1.0).
		SetSize(1.0).
		SetPrice(1.0))
	if err != nil {
		log.Fatal(err)
	}

	// Cancel the order since it will almost definitely never get filled.
	if _, err := client.CancelOrder(order.ID); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_CreateProfile() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_CreateProfileTransfer() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_CreateReport() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	startTimestamp := time.Date(2018, 11, 9, 17, 11, 8, 0, time.UTC).Format(time.RFC3339)
	endTimestamp := time.Date(2018, 11, 9, 18, 11, 8, 0, time.UTC).Format(time.RFC3339)

	_, err = client.CreateReport(new(coinbasepro.CreateReportOptions).
		SetType(coinbasepro.ReportTypeFills).
		SetFormat(coinbasepro.FileFormatPdf).
		SetProfileID(profileID).
		SetProductID("BTC-USD").
		SetStartDate(startTimestamp).
		SetEndDate(endTimestamp))
	if err != nil {
		log.Fatalf("Error creating report: %v", err)
	}
}

func ExampleClient_CryptoWithdrawal() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a wallet ID for the request.
	walletID := os.Getenv("CB_PRO_WALLET_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Get an address for the target wallet.
	address, _ := client.GenerateCryptoAddress(walletID)

	// Withdraw the funds using the generated USDC wallet address.
	_, err = client.CryptoWithdrawal(new(coinbasepro.CryptoWithdrawalOptions).
		SetCryptoAddress(address.Address).
		SetAmount(1).
		SetCurrency("USDC"))
	if err != nil {
		log.Fatalf("Error withdrawing crypto: %v", err)
	}
}

func ExampleClient_Currencies() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	currencies, err := client.Currencies()
	if err != nil {
		log.Fatalf("Error listing currencies: %v", err)
	}
	fmt.Printf("currencies: %+v\n", currencies)
}

func ExampleClient_Currency() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a valid currency.
	accountCurrency := os.Getenv("CB_PRO_ACCOUNT_CURRENCY")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	currency, err := client.Currency(accountCurrency)
	if err != nil {
		log.Fatalf("Error fetching currency by ID: %v", err)
	}
	fmt.Printf("currency: %+v\n", currency)
}

func ExampleClient_DeleteProfile() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_ExchangeLimits() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a valid User ID from the account's profile.
	userID := os.Getenv("CB_PRO_USER_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	exchangeLimits, err := client.ExchangeLimits(userID)
	if err != nil {
		log.Fatalf("Error fetching exchange limits: %v", err)
	}
	fmt.Printf("exchange limits: %+v\n", exchangeLimits)
}

func ExampleClient_Fees() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	fees, err := client.Fees()
	if err != nil {
		log.Fatalf("Error fetching fees: %v", err)
	}
	fmt.Printf("fees: %+v\n", fees)
}

func ExampleClient_Fills() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")
	productID := os.Getenv("CB_PRO_PRODUCT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	fills, err := client.Fills(new(coinbasepro.FillsOptions).
		SetAfter(1526365354).
		SetBefore(1652574195).
		SetLimit(1).
		SetProductID(productID).
		SetProfileID(profileID))
	if err != nil {
		log.Fatalf("Error fetching fills: %v", err)
	}
	fmt.Printf("fills: %+v\n", fills)
}

func ExampleClient_GenerateCryptoAddress() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a wallet ID for the request.
	walletID := os.Getenv("CB_PRO_WALLET_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	address, err := client.GenerateCryptoAddress(walletID)
	if err != nil {
		log.Fatalf("Error fetching address: %v", err)
	}
	fmt.Printf("USD wallet address: %+v\n", address)
}

func ExampleClient_Order() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Lookup the order.
	order, _ := client.Order("your-order-id")
	fmt.Printf("order: %+v\n", order)
}

func ExampleClient_Orders() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")
	productID := os.Getenv("CB_PRO_PRODUCT_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	startTimestamp := time.Date(2018, 11, 9, 17, 11, 8, 0, time.UTC).Format(time.RFC3339)
	endTimestamp := time.Date(2025, 11, 9, 18, 11, 8, 0, time.UTC).Format(time.RFC3339)

	// Lookup the order.
	orders, err := client.Orders(new(coinbasepro.OrdersOptions).
		SetAfter("2023-01-01").
		SetBefore("2010-01-01").
		SetStartDate(startTimestamp).
		SetEndDate(endTimestamp).
		SetLimit(1).
		SetProductID(productID).
		SetProfileID(profileID))
	if err != nil {
		log.Fatalf("Error fetching order: %v", err)
	}
	fmt.Printf("orders: %+v\n", orders)
}

func ExampleClient_PaymentMethodDeposit() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_PaymentMethodWithdrawal() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_PaymentMethods() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	paymentMethods, err := client.PaymentMethods()
	if err != nil {
		log.Fatalf("Error fetching payment methods: %v", err)
	}
	fmt.Printf("payment methods: %+v\n", paymentMethods)
}

func ExampleClient_Product() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	product, err := client.Product("BTC-USD")
	if err != nil {
		log.Fatalf("Error fetching product: %v", err)
	}
	fmt.Printf("product: %+v\n", product)
}

func ExampleClient_ProductStats() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	stats, err := client.ProductStats("BTC-USD")
	if err != nil {
		log.Fatalf("Error fetching stats: %v", err)
	}
	fmt.Printf("stats: %+v\n", stats)
}

func ExampleClient_ProductTicker() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ticker, err := client.ProductTicker("BTC-USD")
	if err != nil {
		log.Fatalf("Error fetching ticker: %v", err)
	}
	fmt.Printf("ticker: %+v\n", ticker)
}

func ExampleClient_Products() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	products, err := client.Products(new(coinbasepro.ProductsOptions).SetType("USD-BTC"))
	if err != nil {
		log.Fatalf("Error fetching products: %v", err)
	}
	fmt.Printf("products: %+v\n", products)
}

func ExampleClient_Profile() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	profile, err := client.Profile(profileID, new(coinbasepro.ProfileOptions).SetActive(true))
	if err != nil {
		log.Fatalf("Error fetching profile: %v", err)
	}
	fmt.Printf("profile: %+v\n", profile)
}

func ExampleClient_Profiles() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	profiles, err := client.Profiles(new(coinbasepro.ProfilesOptions).SetActive(true))
	if err != nil {
		log.Fatalf("Error fetching profiles: %v", err)
	}
	fmt.Printf("profiles: %+v\n", profiles)
}

func ExampleClient_RenameProfile() {
	// TODO: Figure out why we get a 403 HTTP Status
}

func ExampleClient_Report() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	report, _ := client.Report("your-report-id")
	fmt.Printf("report: %+v\n", report)
}

func ExampleClient_Reports() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a profile ID for the request.
	profileID := os.Getenv("CB_PRO_PROFILE_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	afterTimestamp := time.Date(2018, 11, 9, 17, 11, 8, 0, time.UTC).Format(time.RFC3339)

	reports, err := client.Reports(new(coinbasepro.ReportsOptions).
		SetAfter(afterTimestamp).
		SetIgnoredExpired(true).
		SetLimit(1).
		SetPortfolioID(profileID).
		SetType(coinbasepro.ReportTypeAccounts))
	if err != nil {
		log.Fatalf("Error fetching reports: %v", err)
	}
	fmt.Printf("reports: %+v\n", reports)
}

func ExampleClient_SignedPrices() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	prices, err := client.SignedPrices()
	if err != nil {
		log.Fatalf("Error fetching prices: %v", err)
	}
	fmt.Printf("prices: %+v\n", prices)
}

func ExampleClient_Trades() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	trades, err := client.Trades("BTC-USD", new(coinbasepro.TradesOptions).
		SetAfter(1526365354).
		SetBefore(1652574165).
		SetLimit(1))
	if err != nil {
		log.Fatalf("Error fetching prices: %v", err)
	}
	fmt.Printf("prices: %+v\n", trades)
}

func ExampleClient_Transfer() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Lookup the deposit.
	transfer, _ := client.Transfer("some-transfer-id")
	fmt.Printf("transfer: %+v\n", transfer)
}

func ExampleClient_Transfers() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	transfers, err := client.Transfers()
	if err != nil {
		log.Fatalf("Error fetching transfers: %v", err)
	}
	fmt.Printf("transfers: %+v\n", transfers)
}

func ExampleClient_Wallets() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	wallets, err := client.Wallets()
	if err != nil {
		log.Fatalf("Error fetching wallets: %v", err)
	}
	fmt.Printf("wallets: %+v\n", wallets)
}

func ExampleClient_WithdrawalFeeEstimate() {
	// Read credentials from environment variables.
	url := os.Getenv("CB_PRO_URL")
	passphrase := os.Getenv("CB_PRO_ACCESS_PASSPHRASE")
	key := os.Getenv("CB_PRO_ACCESS_KEY")
	secret := os.Getenv("CB_PRO_SECRET")

	// Get a USD Wallet ID for the request.
	walletID := os.Getenv("CB_PRO_USDC_WALLET_ID")

	// Get a new client using an API Key for authentication.
	client, err := coinbasepro.NewClient(context.TODO(), transport.NewAPIKey().
		SetKey(key).
		SetPassphrase(passphrase).
		SetSecret(secret).
		SetURL(url))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	address, err := client.GenerateCryptoAddress(walletID)
	if err != nil {
		log.Fatalf("Error generating crypto address: %v", err)
	}

	estimates, err := client.WithdrawalFeeEstimate(new(coinbasepro.WithdrawalFeeEstimateOptions).
		SetCryptoAddress(address.Address).
		SetCurrency("USDC"))
	if err != nil {
		log.Fatalf("Error fetching estimates: %v", err)
	}
	fmt.Printf("estimates: %+v\n", estimates)
}

func ExampleProductWebsocket_Ticker() {
	ws := coinbasepro.NewWebsocket(websocket.DefaultConnector)

	// initialize the ticker object to channel product messages
	ticker := ws.Ticker("ETH-USD")

	// start a go routine that passes product messages concerning ETH-USD currency pair to a channel on the ticker struct.
	ticker.Open()
	go func() {
		// Next we range over the product message channel and print the product messages.
		for productMsg := range ticker.Channel() {
			fmt.Printf("ETH-USD Price @ %v: %v\n", productMsg.Time, productMsg.Price)
		}
	}()

	// Let the product messages print for 5 seconds.
	time.Sleep(5 * time.Second)

	// Then close the ticker channel, this will unsubscribe from the websocket and close the underlying channel that the
	// messages read to.
	ticker.Close()
}

// cancelOpenOrders will cancel all of the open orders enqueued to free up cash to make new orders for the test.
func cancelOpenOrders(t *testing.T, client *coinbasepro.Client) {
	// Check to see if there are already transfers for this account.
	_, err := client.CancelOpenOrders(nil)
	require.NoError(t, err)
}

// createDepositTransfer will create a transfer for the given account.  If there are already transfer for the account,
// then this this function does nothing.
func createDepositTransfer(t *testing.T, client *coinbasepro.Client, accountID string, currency string, walletID string) {
	// Check to see if there are already transfers for this account.
	transfers, err := client.AccountTransfers(accountID, new(coinbasepro.AccountTransfersOptions).
		SetBefore("2010-01-01").
		SetAfter("2080-01-01").
		SetType(coinbasepro.TransferMethodDeposit))
	require.NoError(t, err)
	if len(transfers) > 0 {
		return
	}

	_, err = client.CoinbaseAccountDeposit(new(coinbasepro.CoinbaseAccountDepositOptions).
		SetCoinbaseAccountID(walletID).
		SetAmount(1.0).
		SetCurrency(currency))
	require.NoError(t, err)
}

// createUnlikelyLimitOrderForCancellation will set an unlikely limit order to be canceled in testing.
func createUnlikelyLimitOrderForCancellation(t *testing.T, client *coinbasepro.Client, profileID string, productID string) {
	order, err := client.CreateOrder(new(coinbasepro.CreateOrderOptions).
		SetProfileID(profileID).
		SetType(coinbasepro.OrderTypeLimit).
		SetSide(coinbasepro.SideSell).
		SetSTP(coinbasepro.STPDc).
		SetStop(coinbasepro.StopLoss).
		SetTimeInForce(coinbasepro.TimeInForceGTC).
		SetCancelAfter(coinbasepro.CancelAfterMin).
		SetProductID(productID).
		SetStopPrice(1.0).
		SetSize(1.0).
		SetPrice(1.0))
	fmt.Println(productID)
	require.NoError(t, err)
	os.Setenv(testOrderIDForCancellation, order.ID)
}

// setAccountIDEnvVariable will set an accountID to the environment variables for the given currency.
func setAccountEnvVariables(t *testing.T, client *coinbasepro.Client, currency string) {
	accounts, err := client.Accounts()
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	for _, account := range accounts {
		if account.Currency == currency {
			os.Setenv(testAccountID, account.ID)
			os.Setenv(testCurrency, account.Currency)
			os.Setenv(testProfileID, account.ProfileID)

			profile, err := client.Profile(account.ProfileID, nil)
			require.NoError(t, err)
			os.Setenv(testUserID, profile.UserID)
			return
		}
	}
}

// setWalletID will set a walletID to the environment variabels for the given currency.
func setWalletEnvVariables(t *testing.T, client *coinbasepro.Client, currency string) {
	wallets, err := client.Wallets()
	require.NoError(t, err)
	require.NotEmpty(t, wallets)
	for _, wallet := range wallets {
		if wallet.Currency == currency {
			os.Setenv(testWalletID, wallet.ID)
		}
		if wallet.Currency == "USDC" {
			os.Setenv(testUSDCWalletID, wallet.ID)
		}
	}
	return
}

// createUnlikelyLimitOrder will place a limit order that should never execute for the given productID.  If an order
// has already been placed, then this function does nothing.
func createUnlikelyLimitOrder(t *testing.T, client *coinbasepro.Client, profileID string, productID string) {
	// Check to see if there are already orders for this product.
	orders, err := client.Orders(new(coinbasepro.OrdersOptions).SetProductID(productID).SetLimit(1))
	require.NoError(t, err)
	if len(orders) > 0 {
		return
	}

	_, err = client.CreateOrder(new(coinbasepro.CreateOrderOptions).
		SetProfileID(profileID).
		SetType(coinbasepro.OrderTypeLimit).
		SetSide(coinbasepro.SideSell).
		SetSTP(coinbasepro.STPDc).
		SetStop(coinbasepro.StopLoss).
		SetTimeInForce(coinbasepro.TimeInForceGTC).
		SetCancelAfter(coinbasepro.CancelAfterMin).
		SetProductID(productID).
		SetStopPrice(1.0).
		SetSize(1.0).
		SetPrice(1.0))
	require.NoError(t, err)
}
