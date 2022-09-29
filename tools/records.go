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

var (
	ErrUnsupportedDataType     = fmt.Errorf("unsupported data type")
	ErrFailedToAssertInterface = fmt.Errorf("failed to assert interface")
	ErrFailedToMarshalBSON     = fmt.Errorf("failed to marshal bson")
	ErrFailedToUnmarshalBSON   = fmt.Errorf("failed to unmarshal bson")
	ErrFailedToMarshalJSON     = fmt.Errorf("failed to marshal json")
	ErrFailedToUnmarshalJSON   = fmt.Errorf("failed to unmarshal json")
	ErrFailedToCreateStruct    = fmt.Errorf("failed to create struct")
	ErrFailedToScanRow         = fmt.Errorf("failed to scan row")
	ErrFailedToParseFloat      = fmt.Errorf("failed to parse float")
	ErrFailedToDecodeRecords   = fmt.Errorf("failed to decode records")
	ErrFailedToGetColumns      = fmt.Errorf("failed to get columns")
)

type Encoder interface {
	EncodeBody() (io.Reader, error)
	EncodeQuery(*http.Request)
}

func AssingRecordBSONDocument(req *structpb.Struct, doc *bson.D) error {
	data, err := bson.Marshal(req.AsMap())
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToMarshalBSON, err)
	}

	err = bson.Unmarshal(data, doc)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToUnmarshalBSON, err)
	}

	return nil
}

// AssignReadOptions will assign an options struct to the read request.
func AssignReadOptions(req *proto.ReadRequest, opts Encoder) error {
	bytes, err := json.Marshal(opts)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToMarshalJSON, err)
	}

	optsMap := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &optsMap); err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToUnmarshalJSON, err)
	}

	req.Options, err = structpb.NewStruct(optsMap)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToCreateStruct, err)
	}

	return nil
}

func AssignReadResponseRecords(rsp *proto.ReadResponse, dest interface{}) error {
	elem := reflect.ValueOf(dest).Elem()
	switch elem.Kind() {
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
				return fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
			}

			model := reflect.New(nonptrStructType).Interface()

			err = json.Unmarshal(bytes, &model)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
			}

			result.Index(idx).Set(reflect.ValueOf(model))
		}

		reflect.ValueOf(dest).Elem().Set(result)
	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Map, reflect.Pointer, reflect.String, reflect.Struct,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
		return fmt.Errorf("%w: %v", ErrUnsupportedDataType, elem.Kind())
	}

	return nil
}

// AssignReadRequired will assign a key value to the required struct on the read request.
func AssignReadRequired(req *proto.ReadRequest, key string, val interface{}) error {
	var err error
	if req.Required == nil {
		req.Required, err = structpb.NewStruct(map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("%v: %w", ErrFailedToCreateStruct, err)
		}
	}

	m := req.Required.AsMap()
	m[key] = val

	req.Required, err = structpb.NewStruct(m)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToCreateStruct, err)
	}

	return nil
}

// scanColumnPointers will scan a row into a slice of pointers, returning the slice.
func scanColumnPointers(rows *sql.Rows, columns []string) ([]interface{}, error) {
	// Create a slice of interface{}'s to represent each column,
	// and a second slice to contain pointers to each item in the columns slice.
	columnReps := make([]interface{}, len(columns))
	columnPointers := make([]interface{}, len(columns))

	for i := range columns {
		columnPointers[i] = &columnReps[i]
	}

	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		return nil, fmt.Errorf("%v: %w", ErrFailedToScanRow, err)
	}

	return columnPointers, nil
}

// sortColumnValues will sort the values of a row into a map of column name to value.
func sortColumnValues(rows *sql.Rows, columns []string) (map[string]interface{}, error) {
	columnPointers, err := scanColumnPointers(rows, columns)
	if err != nil {
		return nil, err
	}

	// Create our map, and retrieve the value for each column from the pointers slice,
	// storing it in the map with the name of the column as the key.
	colVal := make(map[string]interface{})

	for i, colName := range columns {
		val, ok := columnPointers[i].(*interface{})
		if !ok {
			return nil, ErrFailedToAssertInterface
		}

		if reflect.TypeOf(*val) == reflect.TypeOf([]byte{}) {
			// The postgres driver treats numbers & decimal columns as []uint8. We chose to parse
			// these values into strings.
			f, err := strconv.ParseFloat(string((*val).([]byte)), strconv.IntSize)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", ErrFailedToParseFloat, err)
			}

			*val = f
		}

		colVal[colName] = *val
	}

	return colVal, nil
}

// AssignStructs will convert SQL rows into structpb.Struct value and append the slice passed into the function,
// genreralizing the process of createing JSON objects from SQL rows.
func AssignStructs(rows *sql.Rows, val *[]*structpb.Struct) error {
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("%v: %w", ErrFailedToGetColumns, err)
	}
	defer rows.Close()

	for rows.Next() {
		// Scan the rows into a slice of pointers, then sort them by column name.
		colVal, err := sortColumnValues(rows, cols)
		if err != nil {
			return err
		}

		// Encoded the data and append it to the response tables.
		encodedData, err := json.Marshal(colVal)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
		}

		pbstruct := &structpb.Struct{}
		if err = pbstruct.UnmarshalJSON(encodedData); err != nil {
			return fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
		}

		*val = append(*val, pbstruct)
	}

	return nil
}

// decodeRecords will parse a slice of data into a records slice.
func decodeRecords(data interface{}) ([]*structpb.Struct, error) {
	var out []interface{}

	dataValue := reflect.ValueOf(data)
	switch dataValue.Kind() {
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			out = append(out, dataValue.Index(i).Interface())
		}
	case reflect.Map:
		out = append(out, dataValue.Interface())

	case reflect.Array, reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32,
		reflect.Float64, reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Interface, reflect.Invalid, reflect.Pointer, reflect.String, reflect.Struct,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, dataValue.Kind())
	}

	records := make([]*structpb.Struct, 0)

	for _, r := range out {
		record, err := json.Marshal(r)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToMarshalJSON, err)
		}

		rec := new(structpb.Struct)

		err = rec.UnmarshalJSON(record)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
		}

		records = append(records, rec)
	}

	return records, nil
}

// UpsertDataType are the supported types for decoding upsert records.
type UpsertDataType uint8

const (
	// UpsertDataJSON is the default upsert data type.
	UpsertDataJSON UpsertDataType = iota
)

// DecodeUpsertRecords will decode the records from the upsert request into a slice of structs.
func DecodeUpsertRecords(req *proto.UpsertRequest) ([]*structpb.Struct, error) {
	if UpsertDataType(req.DataType) == UpsertDataJSON {
		var data interface{}
		if err := json.Unmarshal(req.Data, &data); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToUnmarshalJSON, err)
		}

		records, err := decodeRecords(data)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedToDecodeRecords, err)
		}

		return records, nil
	}

	return nil, fmt.Errorf("%w: %v", ErrUnsupportedDataType, req.DataType)
}

// PartitionStructs ensures that the request structures are partitioned into size n or less-sized chunks of data, to
// comply with insert requirements.
func PartitionStructs(size int, slice []*structpb.Struct) [][]*structpb.Struct {
	var chunks [][]*structpb.Struct

	for len(slice) > 0 {
		if len(slice) < size {
			size = len(slice)
		}

		chunks = append(chunks, slice[0:size])

		slice = slice[size:]
	}

	return chunks
}
