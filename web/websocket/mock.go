package websocket

type Mock struct{}

func NewMock() (Connector, error) {
	return &Mock{}, nil
}

func (conn *Mock) ReadJSON(_ interface{}) error {
	return nil
}

func (conn *Mock) WriteJSON(_ interface{}) error {
	return nil
}
