package coinbasepro

import "github.com/alpine-hodler/web/pkg/websocket"

type ProductWebsocket struct {
	conn websocket.Connector
}

// NewWebsocket will create a connection to the coinbase websocket and
// return a singleton that can be used to open channels that stream product
// data via a websocket.
func NewWebsocket(ws websocket.Creator) *ProductWebsocket {
	productWebsocket := new(ProductWebsocket)
	productWebsocket.conn, _ = ws(websocketURL)
	return productWebsocket
}

// Ticker ticker uses the ProductWebsocket connection to query coinbase for ticket data, then it puts that data onto a
// channel for model.CoinbaseTicker
func (productWebsocket *ProductWebsocket) Ticker(products ...string) *AsyncTicker {
	return newAsyncTicker(productWebsocket.conn, products...)
}
