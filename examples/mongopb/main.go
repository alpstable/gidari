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
	"log"
	"net/http"
	"time"

	"github.com/alpstable/gidari"
	"github.com/alpstable/mongopb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

func main() {
	ctx := context.Background()

	const api = "https://anapioficeandfire.com/api"

	// Create a mongoDB client.
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to create mongo client: %v", err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect from mongo: %v", err)
		}
	}()

	// Create a collection to store the api data.
	bookColl := client.Database("test").Collection("books")
	charColl := client.Database("test").Collection("characters")

	// Create writers for the books and characters collections.
	bookWriter := mongopb.NewListWriter(bookColl)
	charWriter := mongopb.NewListWriter(charColl)

	// Create a service that can be used to make bulk HTTP requests to the
	// API.
	svc, err := gidari.NewService(ctx)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Create some requests and add them to the service.
	charReq, _ := http.NewRequest(http.MethodGet, api+"/characters", nil)
	housReq, _ := http.NewRequest(http.MethodGet, api+"/houses", nil)

	// Wrap the HTTP Requests in the gidalri.HTTPRequest type.
	charReqWrapper := gidari.NewHTTPRequest(charReq, gidari.WithWriter(bookWriter))
	housReqWrapper := gidari.NewHTTPRequest(housReq, gidari.WithWriter(charWriter))

	svc.HTTP.Requests(charReqWrapper, housReqWrapper)

	// Add a rate limiter to the service, 5 requests per second. This can
	// help avoid "429" errors.
	requestPerSecond := 5
	svc.HTTP.RateLimiter(rate.NewLimiter(rate.Every(1*time.Second), requestPerSecond))

	// Use Upsert to make requests and in our case write response data to
	// the MongoDB collectionsMongoDB collections.
	if err := svc.HTTP.Upsert(ctx); err != nil {
		log.Fatalf("failed to upsert HTTP responses: %v", err)
	}
}
