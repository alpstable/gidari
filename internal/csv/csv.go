package csv

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/alpstable/gidari/internal/proto"
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

// decodeUpsertRequest will convert the JSON data from a "proto.UpsertRequest" into CSV data.
//
// It is not possible to only get a header rows as the header row is determined by the keys of a map.
func decodeUpsertRequest(req *proto.UpsertRequest, rows chan<- []interface{}) error {
	defer close(rows)

	// If the data is empty, then there is nothing to do.
	if reflect.DeepEqual(req.Data, []byte(``)) {
		return nil
	}

	records, err := proto.DecodeUpsertRequest(req)
	if err != nil {
		return err
	}

	header := make(map[string]int)

	for _, record := range records {
		rmap := record.AsMap()

		// If the header is empty, the range over the map and add the keys to the header.
		if len(header) == 0 {
			for k := range rmap {
				header[k] = len(header)
			}
		}

		// If the header is smaller than the record, then add the missing keys to the end of the header.
		if len(header) < len(rmap) {
			for k := range rmap {
				if _, ok := header[k]; !ok {
					header[k] = len(header)
				}
			}
		}

		// Iterate over the header and add the values to the values slice.
		values := make([]interface{}, len(header))
		for key, index := range header {
			values[index] = rmap[key]
		}

		rows <- values
	}

	// build the header row
	headerRow := make([]interface{}, len(header))
	for key, index := range header {
		headerRow[index] = key
	}

	rows <- headerRow

	return err
}

// Upsert will insert or update a record in the directory. CSV does not support "Upsert" in the traditional sense. This
// function will simply do one of two things (1) if the file defined by the "Table" column already exists, then this
// function will simply append the data to the end of the file. (2) if the file does not exist, then this function will
// create the file and write the data to it.
//
// Consistant column order is not guarunteed.
func (csv *CSV) Upsert(ctx context.Context, req *proto.UpsertRequest) (*proto.UpsertResponse, error) {

	return nil, nil
}
