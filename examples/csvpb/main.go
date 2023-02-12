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
)

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

	// Create a writer object with 'encoding/csv' package,
	// then use csvpb library to create list writer.
	writer := csv.NewWriter(os.Stdout)
	listWriter := csvpb.NewListWriter(writer, csvpb.WithAlphabetizeHeaders())

	// Create some requests and add them to the service.
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters/583", nil) // Jon Snow
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses/10", nil)      // House Baelish

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: charReq, Writer: listWriter}).
		Requests(&gidari.HTTPRequest{Request: housReq, Writer: listWriter})

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	requestPerSecond := 5
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), requestPerSecond))

	// Use Upsert to make requests and in our case write response data to stdout in CSV format.
	if err := svc.HTTP.Upsert(ctx); err != nil {
		log.Fatalf("failed to upsert HTTP responses: %v", err)
	}

	// Flush the writer.
	writer.Flush()
}
