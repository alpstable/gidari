package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/alpstable/gidari"
	"golang.org/x/time/rate"
)

func main() {
	ctx := context.Background()

	iter, err := gidari.NewIterator(ctx, asoiafConfig())
	if err != nil {
		log.Fatalf("failed to create iterator: %v", err)
	}

	defer iter.Close(ctx)

	urlByteSize := map[string][]int{}

	for iter.Next(ctx) {
		current := iter.Current

		url := current.GetURL()
		byteSize := len(current.GetData())

		log.Printf("URL: %s, ByteSize: %v", url, byteSize)
		urlByteSize[url] = append(urlByteSize[url], byteSize)
	}

	if err := iter.Err(); err != nil {
		log.Fatalf("iterator error: %v", err)
	}

	const (
		charURL  = "https://www.anapioficeandfire.com/api/characters"
		bookURL  = "https://www.anapioficeandfire.com/api/books"
		houseURL = "https://www.anapioficeandfire.com/api/houses"
	)

	expected := map[string][]int{
		charURL:  {339, 624, 361, 380, 319, 326, 364, 365, 324, 331},
		bookURL:  {24856, 43974, 57575, 3346, 69752, 4990, 4283, 49062, 3454, 3341},
		houseURL: {447, 745, 293, 703, 458, 358, 2187, 676, 342, 887},
	}

	if !reflect.DeepEqual(urlByteSize[charURL], expected[charURL]) {
		log.Fatalf("unexpected character byte want: %v, got: %v", expected[charURL], urlByteSize[charURL])
	}

	if !reflect.DeepEqual(urlByteSize[bookURL], expected[bookURL]) {
		log.Fatalf("unexpected book byte want: %v, got: %v", expected[bookURL], urlByteSize[bookURL])
	}

	if !reflect.DeepEqual(urlByteSize[houseURL], expected[houseURL]) {
		log.Fatalf("unexpected house byte want: %v, got: %v", expected[houseURL], urlByteSize[houseURL])
	}
}

// asoiafConfig is a sample configuration to fetch data from "https://anapioficeandfire.com".
//
// This example illustrates basic usage of the Gidari iterator.
func asoiafConfig() *gidari.Config {
	return &gidari.Config{
		URL: func() *url.URL {
			url, _ := url.Parse("https://www.anapioficeandfire.com")

			return url
		}(),
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
}
