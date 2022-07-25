package query

import (
	_ "embed"
	"encoding/gob"

	"go.mongodb.org/mongo-driver/bson"
)

//go:embed sql/postgres/coinbasepro_candle_minutes.sql
var PostgresCoinbaseProCandleMinutes []byte

//go:embed sql/postgres/polygon_bar_minutes.sql
var PostgresPolygonBarMinutes []byte

//go:embed sql/postgres/columns.sql
var PostgresColumns []byte

//go:embed sql/postgres/tables.sql
var PostgresTables []byte

//go:embed sql/postgres/truncate_tables.sql
var PostgresTruncateTables []byte

type Mongo struct {
	D          bson.D
	Collection string
}

func init() {
	gob.Register(bson.D{})
	gob.Register(&Mongo{})
}
