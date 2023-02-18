// Copyright 2023 The Gidari Authors.
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
	"os"
	"time"

	"github.com/alpstable/gidari"
	"github.com/alpstable/gidari/auth"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/types/known/structpb"
)

func ExampleHTTPIteratorService_Next() {
	ctx := context.TODO()

	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some requests and add them to the service.
	charReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/houses", nil)

	// Wrap the HTTP Requests in the gidalri.HTTPRequest type.
	charReqWrapper := gidari.NewHTTPRequest(charReq)
	housReqWrapper := gidari.NewHTTPRequest(housReq)

	// Add the wrapped HTTP requests to the HTTP Service.
	svc.HTTP.Requests(charReqWrapper, housReqWrapper)

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
	// Total number of bytes: 10455
}

type ExampleWriter struct {
	lists []*structpb.ListValue
}

func (w *ExampleWriter) Write(ctx context.Context, list *structpb.ListValue) error {
	w.lists = append(w.lists, list)

	return nil
}

func ExampleHTTPService_Upsert() {
	ctx := context.TODO()

	const api = "https://anapioficeandfire.com/api"

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some HTTP Requests.
	charReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/houses", nil)

	// Create a writer to write the data.
	w := &ExampleWriter{}

	// Wrap the HTTP Requests in the gidalri.HTTPRequest type.
	charReqWrapper := gidari.NewHTTPRequest(charReq, gidari.WithWriter(w))
	housReqWrapper := gidari.NewHTTPRequest(housReq, gidari.WithWriter(w))

	// Add the wrapped HTTP requests to the HTTP Service.
	svc.HTTP.Requests(charReqWrapper, housReqWrapper)

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
}

func ExampleWithAuth() {
	ctx := context.TODO()

	const api = "https://api-public.sandbox.exchange.coinbase.com"

	key := os.Getenv("COINBASE_API_KEY")
	secret := os.Getenv("COINBASE_API_SECRET")
	passphrase := os.Getenv("COINBASE_API_PASSPHRASE")

	// If these environment variables are not set, then skip the example.
	if key == "" || secret == "" || passphrase == "" {
		return
	}

	// First we create a service that can be used to make bulk HTTP
	// requests to the API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Set a round tripper that will sign the requests.
	roundTripper := auth.NewCoinbaseRoundTrip(key, secret, passphrase)
	withAuth := gidari.WithAuth(roundTripper)

	// Create some requests and add them to the service.
	accounts, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/accounts", nil)
	currencies, _ := http.NewRequestWithContext(ctx, http.MethodGet, api+"/currencies", nil)

	// Wrap the HTTP Requests in the gidalri.HTTPRequest type.
	accountsW := gidari.NewHTTPRequest(accounts, withAuth)
	currenciesW := gidari.NewHTTPRequest(currencies, withAuth)

	// Add the wrapped HTTP requests to the HTTP Service.
	svc.HTTP.Requests(accountsW, currenciesW)

	// Get the status code for the responses.
	for svc.HTTP.Iterator.Next(ctx) {
		current := svc.HTTP.Iterator.Current

		rsp := current.Response
		if rsp == nil {
			break
		}

		fmt.Println("status code:", rsp.Request.URL.Path, rsp.StatusCode)
	}

	// Unordered Output:
	// status code: /accounts 200
	// status code: /currencies 200
}
