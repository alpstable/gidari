package websocket

type Channel struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}
