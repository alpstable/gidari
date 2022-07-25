package websocket

type Connector interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
}

// DefaultWebsocketConnector returns a new websocket as the default connection.
func DefaultConnector(url string) (Connector, error) {
	return New(url)
}
