package repository

// Raw is a struct that holds a table name and a byte slice of data to upsert into that table.
type Raw struct {
	Table string `json:"table"`
	Data  []byte `json:"data"`
}

// NewRaw will return a new Raw struct.
func NewRaw(table string, data []byte) Raw {
	return Raw{Table: table, Data: data}
}
