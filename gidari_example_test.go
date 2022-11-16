package gidari_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
)

func ExampleNewIterator() {
	ctx := context.Background()

	// For this example, we will query the "A Song of Ice and Fire" API.
	url, err := url.Parse("https://anapioficeandfire.com")
	if err != nil {
		log.Fatal(err)
	}

	// Create a configuration for the iterator.
	config := &gidari.Config{
		URL: url,
		Requests: []*gidari.Request{
			{
				Endpoint:    "/api/books",
				Method:      http.MethodGet,
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
			},
			{
				Endpoint:    "/api/characters",
				Method:      http.MethodGet,
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
			},
			{
				Endpoint:    "/api/houses",
				Method:      http.MethodGet,
				RateLimiter: rate.NewLimiter(rate.Every(time.Second), 5),
			},
		},
	}

	// Create an iterator for the configuration.
	iter, err := gidari.NewIterator(ctx, config)
	if err != nil {
		log.Fatalf("failed to create iterator: %v", err)
	}

	defer iter.Close(ctx)

	// byteSize will keep track of the number of bytes for the JSON response in each request.
	var byteSize int

	for iter.Next(ctx) {
		current := iter.Current
		byteSize += len(current.GetData())
	}

	if err := iter.Err(); err != nil {
		log.Fatalf("iterator error: %v", err)
	}

	fmt.Println("Total number of byte:", byteSize)
	// Output:
	// Total number of byte: 256146
}
