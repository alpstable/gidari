// Copyright 2023 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0

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

	// Wrap the HTTP Requests in the gidalri.HTTPRequest type.
	charReqWrapper := gidari.NewHTTPRequest(charReq, gidari.WithWriters(listWriter))
	housReqWrapper := gidari.NewHTTPRequest(housReq, gidari.WithWriters(listWriter))

	svc.HTTP.Requests(charReqWrapper, housReqWrapper)

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	requestPerSecond := 5
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), requestPerSecond))

	// Write the response body to the CSV file.
	if err := svc.HTTP.Store(ctx); err != nil {
		log.Fatalf("failed to upsert HTTP responses: %v", err)
	}

	// Flush the writer.
	writer.Flush()
}
