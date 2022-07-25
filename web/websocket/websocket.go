package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Websocket holds the connection and response data for a coinbase websocket.
type Websocket struct {
	*websocket.Conn
	Response *http.Response
}

type Creator func(string) (Connector, error)

// NewWebsocket will return a new coinbase websocket connection
func New(url string) (Connector, error) {
	conn := new(Websocket)
	var dialer websocket.Dialer
	var err error
	conn.Conn, conn.Response, err = dialer.Dial(url, nil)
	return conn, err
}
