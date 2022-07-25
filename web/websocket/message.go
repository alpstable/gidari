package websocket

type Message struct {
	Type     string    `json:"type"`
	Products []string  `json:"product_ids"`
	Channels []Channel `json:"channels"`
}

// NewWebsocket will return a new websocket subscription that can be used to get
// a feed of real-time market data.
func NewProductsMessage(products []string, channels []Channel) (*Message, error) {
	sub := new(Message)
	sub.Products = products
	sub.Channels = channels
	return sub, nil
}

func (msg *Message) isSub()   { msg.Type = "subscribe" }
func (msg *Message) isUnsub() { msg.Type = "unsubscribe" }

// Subscribe will use the websocket message to subscribe to a websocket
// connections, returning that connection or any errors that occurred attemptting
// to make it.
func (msg *Message) Subscribe(conn Connector) error {
	msg.isSub()
	if err := conn.WriteJSON(msg); err != nil {
		return err
	}
	return nil
}

// Unsubscribe will use the websocket message to unsubscribe to a connection.
func (msg *Message) Unsubscribe(conn Connector) error {
	msg.isUnsub()
	if err := conn.WriteJSON(msg); err != nil {
		return err
	}
	return nil
}
