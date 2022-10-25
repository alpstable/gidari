package csv

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	gocsv "encoding/csv"

	"github.com/alpstable/gidari/internal/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	// ErrNoDir is returned when the directory does not exist.
	ErrNoDir = fmt.Errorf("directory does not exist")
)

// CSV is a wrapper for an "encoding/csv" reader/writer used to perform CRUD operations on CSV files in a given
// directory. The operations for CSV are not ACID.
type CSV struct {
	Dir string // Dir is the directory where the CSV files read/write from/to.
}

// New takes the path to a directory to store the CSV files. It returns a new CSV object for reading and writing to
// CSV files in the given directory.
func New(ctx context.Context, dir string) (*CSV, error) {
	// Check to see if the directory exists.
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return nil, ErrNoDir
	}

	return &CSV{Dir: dir}, nil
}

type row struct {
	header bool
	data   []string
}

// writeState is used to maintain state between calls to writeBody.
type writeState struct {
	// headerRowByTable is a map of table names to the header row for that table.
	headerRowByTable map[string]*row

	// mtx is a mutex used to protect the headerRowByTable map.
	mtx sync.RWMutex
}

func newWriteState() *writeState {
	return &writeState{
		headerRowByTable: make(map[string]*row),
		mtx:              sync.RWMutex{},
	}
}

// flattenStruct will flatten a "structpb.Struct" into a "flattenedStruct", which will contain the header and data
// matching 1-1 with the header index.
func flattenStruct(record *structpb.Struct) (map[string]string, error) {
	flatMap := make(map[string]string)

	for fieldName, fieldValue := range record.GetFields() {
		switch fieldValue.Kind.(type) {
		case *structpb.Value_StructValue:
			// If the value is a struct, this function should be called recursively until a non-struct
			// value is found. For instance, if {a: {b: {c: 1}}} is passed in, the flattened struct, then
			// the header should be "a.b.c" and the data should be "1".
			subStruct, err := flattenStruct(fieldValue.GetStructValue())
			if err != nil {
				return nil, err
			}

			for subFieldName, subFieldValue := range subStruct {
				flatMap[fmt.Sprintf("%s.%s", fieldName, subFieldName)] = subFieldValue
			}
		case *structpb.Value_StringValue:
			flatMap[fieldName] = fieldValue.GetStringValue()
		case *structpb.Value_NumberValue:
			flatMap[fieldName] = fmt.Sprintf("%f", fieldValue.GetNumberValue())
		case *structpb.Value_BoolValue:
			flatMap[fieldName] = fmt.Sprintf("%t", fieldValue.GetBoolValue())
		case *structpb.Value_NullValue:
			flatMap[fieldName] = ""
		default:
			return nil, fmt.Errorf("unknown value type: %s", reflect.TypeOf(fieldValue.Kind))
		}
	}

	return flatMap, nil
}

// addHeader will extend the headersRowByTable slice with any new keys from the "structpb.Struct" map.
func (wstate *writeState) addHeaders(table string, record *structpb.Struct) ([]string, error) {
	// Check the the table exists in the headerRowByTable map, if it doesn't create it.
	if _, ok := wstate.headerRowByTable[table]; !ok {
		wstate.headerRowByTable[table] = &row{header: true, data: []string{}}
	}

	flatMap, err := flattenStruct(record)
	if err != nil {
		return nil, err
	}

	// Get the existing parts of the header and put them in a set.
	headerSet := make(map[string]struct{})
	for _, header := range wstate.headerRowByTable[table].data {
		headerSet[header] = struct{}{}
	}

	// Collect the row data.
	rowData := make([]string, 0, len(flatMap))

	for fieldName := range flatMap {
		// Check to see if the header already exists, if it doesn't add it.
		if _, ok := headerSet[fieldName]; !ok {
			wstate.headerRowByTable[table].data = append(wstate.headerRowByTable[table].data, fieldName)
		}

		rowData = append(rowData, flatMap[fieldName])
	}

	return rowData, nil
}

// getHeader will return the header row for a given table.
func (wstate *writeState) getHeader(table string) *row {
	wstate.mtx.RLock()
	defer wstate.mtx.RUnlock()

	if header, ok := wstate.headerRowByTable[table]; ok {
		return header
	}

	return nil
}

// decodeUpsertRequest will convert the JSON data from a "proto.UpsertRequest" into CSV data.
//
// It is not possible to only get a header rows as the header row is determined by the keys of a map.
func decodeUpsertRequest(req *proto.UpsertRequest, state *writeState, rowCh chan<- *row) error {
	defer close(rowCh)

	// If the data is empty, then there is nothing to do.
	if reflect.DeepEqual(req.Data, []byte(``)) {
		return nil
	}

	records, err := proto.DecodeUpsertRequest(req)
	if err != nil {
		return err
	}

	// Iterate over the records and send them to the row channel.
	for _, record := range records {
		rowData, err := state.addHeaders(req.Table, record)
		if err != nil {
			return err
		}

		rowCh <- &row{data: rowData}
	}

	return err
}

// writeBody will write the body of the CSV file. While writing to the file, writeBody will persist the header row data
// which can change dynamically as the incoming data is processed. The resulting header file is returned by this
// function to be written by the caller.
//
// writeBody must have some way of maintaining a constant state between calls.
func (csv *CSV) writeBody(ctx context.Context, state *writeState, req *proto.UpsertRequest) error {
	rows := make(chan *row)
	go decodeUpsertRequest(req, state, rows)

	// Create the file if it does not exist.
	filname := filepath.Join(csv.Dir, req.Table) + ".csv"

	file, err := os.OpenFile(filname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Open a CSV writer.
	writer := gocsv.NewWriter(file)
	defer writer.Flush()

	// Write the data to the file.
	for row := range rows {
		if row.header {
			state.headerRowByTable[req.Table] = row

			continue
		}

		if err := writer.Write(row.data); err != nil {
			panic(err)
		}
	}

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

// writeHeader will write the header row to the top of the CSV file.
func (csv *CSV) writeHeader(ctx context.Context, table string, header *row) error {
	filename := filepath.Join(csv.Dir, table) + ".csv"

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	// Insert the header into the top of the file
	lines = append(lines, "")
	copy(lines[1:], lines)
	lines[0] = strings.Join(header.data, ",")

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Upsert will insert or update a record in the directory. CSV does not support "Upsert" in the traditional sense. This
// function will simply do one of two things (1) if the file defined by the "Table" column already exists, then this
// function will simply append the data to the end of the file. (2) if the file does not exist, then this function will
// create the file and write the data to it.
//
// Consistant column order is not guarunteed.
func (csv *CSV) Upsert(ctx context.Context, req <-chan *proto.UpsertRequest) (*proto.UpsertResponse, *errgroup.Group) {
	errs, ctx := errgroup.WithContext(ctx)
	errs.Go(func() error {
		// Each iteration of the request may contain records with slightly different "shape". We need to keep
		// track of the headers for each request so that we can maintain a constant "shape" for the resulting
		// CSV files.
		state := newWriteState()

		for r := range req {
			if err := csv.writeBody(ctx, state, r); err != nil {
				fmt.Println("err", err)
				return err
			}
		}

		for table, header := range state.headerRowByTable {
			if err := csv.writeHeader(ctx, table, header); err != nil {
				fmt.Println("err", err)
				return err
			}
		}

		return nil
	})

	return nil, errs
}
