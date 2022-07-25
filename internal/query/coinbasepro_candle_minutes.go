package query

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alpine-hodler/driver/data/proto"
	"github.com/alpine-hodler/driver/web/coinbasepro"
	"go.mongodb.org/mongo-driver/bson"
)

type coinbaseproCandleMinutes struct {
	table string
}

func newCoinbaseProCandleMinutes() (*coinbaseproCandleMinutes, error) {
	return &coinbaseproCandleMinutes{}, nil
}

func (cm *coinbaseproCandleMinutes) ReaderArgs(req *proto.ReadRequest) ([]interface{}, error) {
	cm.table = req.Table
	bytes, err := req.Options.MarshalJSON()
	if err != nil {
		return nil, err
	}
	opts := new(coinbasepro.CandlesOptions)
	if err := json.Unmarshal(bytes, opts); err != nil {
		return nil, err
	}

	required := req.Required.AsMap()
	productID := required["product_id"]
	if productID == nil {
		return nil, fmt.Errorf("product_id is a required field")
	}
	productID = productID.(string)

	args := []interface{}{productID}
	if opts.Start != nil {
		start, err := time.Parse(time.RFC3339, *opts.Start)
		if err != nil {
			return nil, err
		}
		args = append(args, start.Unix())
	} else {
		args = append(args, 0)
	}

	if opts.End != nil {
		end, err := time.Parse(time.RFC3339, *opts.End)
		if err != nil {
			return nil, err
		}
		args = append(args, end.Unix())
	} else {
		// 	Fri Dec 31 9999 07:00:00 GMT+0000
		args = append(args, 253402239600)
	}

	return args, nil
}

func (cm *coinbaseproCandleMinutes) ReaderQuery(stg StorageType, args ...interface{}) ([]byte, error) {
	switch stg {
	case MongoStorage:
		doc := bson.D{{"product_id", args[0]}}
		if len(args) > 0 && args[1] != nil {
			doc = append(doc, bson.E{"unix", bson.D{{"$gte", args[1]}}})
		}
		if len(args) > 1 && args[2] != nil {
			doc = append(doc, bson.E{"unix", bson.D{{"$lte", args[2]}}})
		}

		m := new(Mongo)
		m.D = doc
		m.Collection = cm.table
		var inputBuffer bytes.Buffer
		if err := gob.NewEncoder(&inputBuffer).Encode(m); err != nil {
			return nil, err
		}
		return inputBuffer.Bytes(), nil
	case PostgresStorage:
		return PostgresCoinbaseProCandleMinutes, nil
	default:
		return nil, fmt.Errorf("storage type %q not support for coinbaseproCandleMinutes", stg)
	}
}
