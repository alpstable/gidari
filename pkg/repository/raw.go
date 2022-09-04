package repository

// Raw is a struct that holds a byte slice and a table name.
type Raw struct {
	Table string `json:"table"`
	Data  []byte `json:"data"`
}

// NewRaw will return a new Raw struct.
func NewRaw(table string, data []byte) Raw {
	return Raw{Table: table, Data: data}
}
