package tools

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/alpine-hodler/gidari/proto"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Encoder interface {
	EncodeBody() (io.Reader, error)
	EncodeQuery(*http.Request)
}

func AssingRecordBSONDocument(req *structpb.Struct, doc *bson.D) error {
	data, err := bson.Marshal(req.AsMap())
	if err != nil {
		return err
	}
	err = bson.Unmarshal(data, doc)
	return err
}

// AssignReadOptions will assign an options struct to the read request.
func AssignReadOptions(req *proto.ReadRequest, opts Encoder) error {
	bytes, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	optsMap := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &optsMap); err != nil {
		return err
	}

	req.Options, err = structpb.NewStruct(optsMap)
	if err != nil {
		return err
	}
	return nil
}

func AssignReadResponseRecords(rsp *proto.ReadResponse, dest interface{}) error {
	v := reflect.ValueOf(dest).Elem()
	switch v.Kind() {
	case reflect.Slice:
		structType := reflect.TypeOf(dest).Elem().Elem()

		// If the structType is a poitner, then make it not a pointer.
		nonptrStructType := structType
		if structType.Kind() == reflect.Pointer {
			nonptrStructType = structType.Elem()
		}

		result := reflect.MakeSlice(reflect.SliceOf(structType), len(rsp.Records), len(rsp.Records))
		for idx, record := range rsp.Records {
			bytes, err := record.MarshalJSON()
			if err != nil {
				return err
			}

			model := reflect.New(nonptrStructType).Interface()
			err = json.Unmarshal(bytes, &model)
			if err != nil {
				return err
			}
			result.Index(idx).Set(reflect.ValueOf(model))
		}
		reflect.ValueOf(dest).Elem().Set(result)
	default:
		return fmt.Errorf("assign read response does not support type %T", v)
	}
	return nil
}

// AssignReadRequired will assign a key value to the required struct on the read request.
func AssignReadRequired(req *proto.ReadRequest, key string, val interface{}) error {
	var err error
	if req.Required == nil {
		req.Required, err = structpb.NewStruct(map[string]interface{}{})
		if err != nil {
			return err
		}
	}
	m := req.Required.AsMap()
	m[key] = val
	req.Required, err = structpb.NewStruct(m)
	return nil
}

// AssignStructs will convert SQL rows into structpb.Struct value and append the slice passed into the function,
// genreralizing the process of createing JSON objects from SQL rows.
func AssignStructs(rows *sql.Rows, val *[]*structpb.Struct) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				// The postgres driver treats numbers & decimal columns as []uint8. We chose to parse
				// these values into strings.
				f, err := strconv.ParseFloat(string((*val).([]byte)), 64)
				if err != nil {
					return err
				}
				*val = f
			}
			m[colName] = *val
		}

		// Encoded the data and append it to the response tables.
		encodedData, _ := json.Marshal(m)
		pbstruct := &structpb.Struct{}
		if err = pbstruct.UnmarshalJSON(encodedData); err != nil {
			return err
		}
		*val = append(*val, pbstruct)
	}

	return nil
}

// MakeRecordsRequest will parse a slice of data into a records slice.
func MakeRecordsRequest(data interface{}, records *[]*structpb.Struct) error {
	var out []interface{}
	rv := reflect.ValueOf(data)
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
	case reflect.Map:
		out = append(out, rv.Interface())
	default:
		return fmt.Errorf("record type not supported: %v", rv.Kind())
	}

	for _, r := range out {
		record, _ := json.Marshal(r)
		rec := new(structpb.Struct)
		err := rec.UnmarshalJSON(record)
		if err != nil {
			return fmt.Errorf("error unmarshaling record: %v", err)
		}
		*records = append(*records, rec)
	}

	return nil
}

// UpsertDataType are the supported types for decoding upsert records.
type UpsertDataType uint8

const (
	// UpsertDataJSON is the default upsert data type.
	UpsertDataJSON UpsertDataType = iota
)

// DecodeUpsertRecords will decode the records from the upsert request into a slice of structs.
func DecodeUpsertRecords(req *proto.UpsertRequest) ([]*structpb.Struct, error) {
	switch UpsertDataType(req.DataType) {
	case UpsertDataJSON:
		var records []*structpb.Struct
		var data interface{}
		if err := json.Unmarshal(req.Data, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raw data: %w", err)
		}

		if err := MakeRecordsRequest(data, &records); err != nil {
			return nil, fmt.Errorf("error making records request: %v", err)
		}
		return records, nil
	}
	return nil, fmt.Errorf("unsupported data type: %v", req.DataType)
}

// PartitionStructs ensures that the request structures are partitioned into size n or less-sized chunks of data, to
// comply with insert requirements.
func PartitionStructs(n int, slice []*structpb.Struct) [][]*structpb.Struct {
	var chunks [][]*structpb.Struct
	for len(slice) > 0 {
		if len(slice) < n {
			n = len(slice)
		}
		chunks = append(chunks, slice[0:n])
		slice = slice[n:]
	}
	return chunks
}
