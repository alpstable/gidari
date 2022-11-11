package main

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

func main() {
	ctx := context.Background()

	iter, err := gidari.NewIterator(ctx, asoiafConfig())
	if err != nil {
		log.Fatalf("failed to create iterator: %v", err)
	}

	iter.Next(ctx)
	fmt.Println(iter.Current)

	iter.Next(ctx)
	fmt.Println(iter.Current)

	//for itr.Next(ctx) {
	//	//fmt.Printf("current: %v\n", len(itr.Current.Data))
	//	//fmt.Printf("url: %v\n", itr.Current.Endpoint)
	//}

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
