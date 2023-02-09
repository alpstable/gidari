package main

import (
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alpstable/csvpb"
	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/types/known/structpb"
)

type ExampleWriter struct {
	lists []*structpb.ListValue
}

func (w *ExampleWriter) Write(cxt context.Context, list *structpb.ListValue) error {
	w.lists = append(w.lists, list)
	return nil
}

/**
	Using the Gidari and CSVPB library to write HTTP response data to stdout as CSV
**/

func main() {
	ctx := context.Background()
	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	w := &ExampleWriter{}

	// Create some requests and add them to the service.
	bookReq, _ := http.NewRequest(http.MethodGet, api+"/books/1", nil)        // A Game of Thrones
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters/583", nil) // Jon Snow
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses/10", nil)      // House Baelish

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: bookReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: charReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: housReq, Writer: w})

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), 5))

	// Use Upsert to make requests and in our case gain access to the response data.
	if err := svc.HTTP.Upsert(ctx); err != nil {
		log.Fatalf("failed to upsert HTTP responses: %v", err)
	}

	// Create a writer object with 'encoding/csv' package,
	// then use csvpb library to create list writer.
	writer := csv.NewWriter(os.Stdout)
	listWriter := csvpb.NewListWriter(writer, csvpb.WithAlphabetizeHeaders())

	// Range slice of structpb.ListValues and write them to the list writer,
	// writing the response data to stdout as CSV
	for _, list := range w.lists {
		if err := listWriter.Write(context.TODO(), list); err != nil {
			log.Fatalf("failed to write list: %v", err)
		}
	}

	// Flush the writer.
	writer.Flush()
}
