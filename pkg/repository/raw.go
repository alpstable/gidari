package repository

import "encoding/json"

// Raw is a struct that holds a byte slice and a table name.
type Raw struct {
	Table string `json:"table"`
	Data  []byte `json:"data"`
}

// NewRaw will return a new Raw struct.
func NewRaw(table string, data []byte) Raw {
	return Raw{Table: table, Data: data}
}

// Encode will encode the Raw struct into a byte slice.
func (r Raw) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode will decode a byte slice into a Raw struct.
func (r Raw) Decode(b []byte) error {
	return json.Unmarshal(b, &r)
}
