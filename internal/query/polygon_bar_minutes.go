package query

import (
	"fmt"

	"github.com/alpine-hodler/driver/data/proto"
)

type polygonBarMinutes struct {
	table string
}

func newPolygonBarMinutes() (*polygonBarMinutes, error) {
	return &polygonBarMinutes{}, nil
}

func (cm *polygonBarMinutes) ReaderArgs(req *proto.ReadRequest) ([]interface{}, error) {
	cm.table = req.Table

	required := req.Required.AsMap()

	ticker := required["ticker"]
	if ticker == nil {
		return nil, fmt.Errorf("ticker is a required field")
	}
	ticker = ticker.(string)

	adjusted := required["adjusted"]
	if adjusted == nil {
		return nil, fmt.Errorf("adjusted is a required field")
	}
	adjusted = adjusted.(bool)

	start := required["start"]
	if start == nil {
		return nil, fmt.Errorf("start is a required field")
	}
	start = start.(float64)

	end := required["end"]
	if end == nil {
		return nil, fmt.Errorf("end is a required field")
	}
	end = end.(float64)

	args := []interface{}{ticker, adjusted, start, end}

	return args, nil
}

func (cm *polygonBarMinutes) ReaderQuery(stg StorageType, args ...interface{}) ([]byte, error) {
	switch stg {
	case MongoStorage:
		return nil, fmt.Errorf("mongodb not supported for polygon bar minutes")
	case PostgresStorage:
		return PostgresPolygonBarMinutes, nil
	default:
		return nil, fmt.Errorf("storage type %q not support for polygon_bar_minutes", stg)
	}
}
