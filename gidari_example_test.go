// Copyright 2022 The Gidari Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
package gidari_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/types/known/structpb"
)

func ExampleHTTPService_Iterator() {
	ctx := context.Background()

	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some requests and add them to the service.
	bookReq, _ := http.NewRequest(http.MethodGet, api+"/books", nil)
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses", nil)

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: bookReq}).
		Requests(&gidari.HTTPRequest{Request: charReq}).
		Requests(&gidari.HTTPRequest{Request: housReq})

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), 5))

	// byteSize will keep track of the sum of bytes for each HTTP Response's
	// body.
	var byteSize int

	for svc.HTTP.Iterator.Next(ctx) {
		current := svc.HTTP.Iterator.Current

		rsp := current.Response
		if rsp == nil {
			break
		}

		// Get the byte slice from the response body.
		body, err := io.ReadAll(current.Response.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}

		// Add the number of bytes to the sum.
		byteSize += len(body)
	}

	// Check to see if an error occurred.
	if err := svc.HTTP.Iterator.Err(); err != nil {
		log.Fatalf("failed to iterate over HTTP responses: %v", err)
	}

	fmt.Println("Total number of bytes:", byteSize)
	// Output:
	// Total number of bytes: 256179
}

type ExampleWriter struct {
	lists []*structpb.ListValue
}

func (w *ExampleWriter) Write(ctx context.Context, list *structpb.ListValue) error {
	w.lists = append(w.lists, list)

	return nil
}

func ExampleHTTPService_Upsert() {
	ctx := context.Background()

	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some requests and add them to the service.
	bookReq, _ := http.NewRequest(http.MethodGet, api+"/books", nil)
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses", nil)

	w := &ExampleWriter{}

	svc.HTTP.
		Requests(&gidari.HTTPRequest{Request: bookReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: charReq, Writer: w}).
		Requests(&gidari.HTTPRequest{Request: housReq, Writer: w})

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), 5))

	// Upsert the responses to the database.
	if err := svc.HTTP.Upsert(ctx); err != nil {
		log.Fatalf("failed to upsert HTTP responses: %v", err)
	}

	// Print the result of the mock writer.
	for _, list := range w.lists {
		fmt.Println("list size: ", len(list.Values))
	}

	// Output:
	// list size:  10
	// list size:  10
	// list size:  10
}
